package scanner

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"webprobe/config"
)

// extractLinks提取给定HTML文档中的所有链接
func extractLinks(baseURL string, body *html.Node) []string {
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					//忽略javascript:;链接
					if strings.HasPrefix(link.String(), "javascript:") {
						continue
					}

					//在extractLinks函数中
					if link.IsAbs() {
						links = append(links, link.String())
					} else {
						//使用url.Parse解析基础URL，并使用ResolveReference来正确拼接相对URL
						base, err := url.Parse(baseURL)
						if err != nil {
							continue
						}
						resolvedURL := base.ResolveReference(link)
						links = append(links, resolvedURL.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)
	return links
}

// fetchAndParse获取给定URL的HTML内容并解析链接
func fetchAndParse(u string, timeout time.Duration) ([]string, error) {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	resp, err := client.Get(u)
	if err != nil {
		log.Printf("爬取 URL 失败: %v 连接失败", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("爬取 URL 失败: %s 返回值为 %d", u, resp.StatusCode)
		return nil, fmt.Errorf("爬取 URL 失败: %s 返回值为 %d", u, resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("爬取 URL 失败: %v 解析失败", err)
		return nil, fmt.Errorf("parsing %s as HTML: %v", u, err)
	}

	baseURL := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host
	return extractLinks(baseURL, doc), nil
}

// Crawler函数爬取给定的URL并按照深度爬取链接
func Crawler(startURL string, config config.CrawlConfig, urls []URLDepth) []URLDepth {
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	var crawl func(string, int, string)                          //添加一个父链接参数
	crawl = func(u string, currentDepth int, fatherURL string) { //接收父链接参数
		defer wg.Done()
		if currentDepth >= config.Depth {
			return
		}

		links, err := fetchAndParse(u, config.Timeout)
		if err != nil {
			return
		}

		mu.Lock()
		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				urls = append(urls, URLDepth{URL: link, Depth: currentDepth + 1, FatherURL: u}) //记录父链接
				wg.Add(1)
				go crawl(link, currentDepth+1, u) //传递当前URL作为子URL的父链接
			}
		}
		mu.Unlock()
	}

	// 在开始时，父链接为空或可以是自身
	urls = append(urls, URLDepth{URL: startURL, Depth: 0, FatherURL: ""})
	visited[startURL] = true
	// 最外层的goroutine启动前计数器加1
	wg.Add(1)
	go crawl(startURL, 0, "") //初始调用时父链接为空

	wg.Wait()

	return urls
}
