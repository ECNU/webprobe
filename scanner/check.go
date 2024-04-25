package scanner

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"webprobe/config"
)

func Check(urls []URLDepth, config config.CheckConfig, urlData []URLStatus) []URLStatus {
	maxConcurrency := config.Concurrency
	ctx, cancel := context.WithTimeout(context.Background(), config.ContextTimeout*time.Second)
	defer cancel()

	sem := semaphore.NewWeighted(int64(maxConcurrency))
	var wg sync.WaitGroup

	// 消费者仅append 快于生产者 故创建非缓冲通道
	statusCh := make(chan URLStatus, 0)
	errCh := make(chan error, 1)

	go func() {
		for status := range statusCh {
			urlData = append(urlData, status)
		}
	}()

	for _, url := range urls {
		if config.UseIPV6 {
			wg.Add(2)
			go checkURLWithSemaphore(config, ctx, sem, url, "ipv4", &wg, statusCh, errCh)
			go checkURLWithSemaphore(config, ctx, sem, url, "ipv6", &wg, statusCh, errCh)
		} else {
			wg.Add(1)
			go checkURLWithSemaphore(config, ctx, sem, url, "ipv4", &wg, statusCh, errCh)
		}
	}
	wg.Wait()
	close(statusCh)
	close(errCh)
	return urlData
}

// checkURLWithSemaphore封装checkURL调用，包括信号量逻辑
func checkURLWithSemaphore(config config.CheckConfig, ctx context.Context, sem *semaphore.Weighted, url URLDepth, ipVersion string, wg *sync.WaitGroup, statusCh chan<- URLStatus, errCh chan<- error) {
	defer wg.Done()
	if err := sem.Acquire(ctx, 1); err != nil {
		errCh <- fmt.Errorf("failed to acquire semaphore for %s: %w", url.URL, err)
		return
	}
	defer sem.Release(1)

	var status URLStatus
	var err error

	//尝试检测直到达到最大重试次数或者成功
	for i := 0; i <= config.Retry.Time; i++ {
		status, err = checkURL(config, ctx, url, ipVersion)
		if err == nil {
			statusCh <- status
			return // 检测成功，发送状态并返回
		}
		if i <= config.Retry.Time-1 {
			time.Sleep(config.Retry.Interval * time.Second) // 在重试之间添加间隔
			log.Printf("该url前次检测失败，第 %d 次重试检测: %s", i+1, url.URL)
		}
	}

	// 如果达到重试次数后仍失败，记录失败状态
	if err != nil {
		failedStatus := URLStatus{
			FatherURL: url.FatherURL,
			Depth:     url.Depth,
			URL:       url.URL,
			IPVersion: ipVersion,
		}
		statusCh <- failedStatus
	}
}

// 检查URL并返回其状态
func checkURL(config config.CheckConfig, ctx context.Context, url URLDepth, ipVersion string) (URLStatus, error) {
	//log.Printf("正在检测: %s ", url.URL)
	start := time.Now()

	//创建HTTP客户端
	client := makeHTTPClient(config, ipVersion)

	// 使用带有超时的子上下文创建请求，用于对不同的url进行不同的超时控制
	//reqCtx, cancel := context.WithTimeout(ctx, config.HttpClientTimeout*time.Second)
	//defer cancel()
	//req, err := http.NewRequestWithContext(reqCtx, "GET", url.URL, nil)

	//使用统一ctx创建GET请求
	req, err := http.NewRequestWithContext(ctx, "GET", url.URL, nil)
	if err != nil {
		return URLStatus{FatherURL: url.FatherURL, Depth: url.Depth, URL: url.URL, IPVersion: ipVersion}, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return URLStatus{FatherURL: url.FatherURL, Depth: url.Depth, URL: url.URL, IPVersion: ipVersion}, err
	}
	defer closeBody(resp.Body)
	status := URLStatus{
		FatherURL:  url.FatherURL,
		Depth:      url.Depth,
		URL:        url.URL,
		IPVersion:  ipVersion,
		Up:         true,
		StatusCode: resp.StatusCode,
		Latency:    time.Since(start),
	}
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		status.CertExpire = resp.TLS.PeerCertificates[0].NotAfter
	}

	return status, nil
}

// makeHTTPClient创建一个根据IP版本配置的HTTP客户端
func makeHTTPClient(checkConfig config.CheckConfig, ipVersion string) *http.Client {
	dialer := &net.Dialer{
		Timeout: checkConfig.DialerTimeout * time.Second,
		//KeepAlive: checkConfig.DialerAliveTimeout * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if ipVersion == "ipv6" {
				return dialer.DialContext(ctx, "tcp6", addr)
			}
			return dialer.DialContext(ctx, "tcp4", addr)
		},
	}

	return &http.Client{
		Timeout:   checkConfig.HttpClientTimeout * time.Second,
		Transport: transport,
	}
}

// closeBody安全地关闭响应体
func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		fmt.Printf("Failed to close response body: %v\n", err)
	}
}
