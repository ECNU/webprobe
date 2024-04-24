package handler

import (
	"html/template"
	"log"
	"os"
	"time"
	"webprobe/scanner"
)

// GenerateHTMLFile 生成包含URL数据的HTML文件
func GenerateHTMLFile(urlDataWithReachability []scanner.URLDataWithReachability, outputFilePath string, depth int) {

	funcMap := template.FuncMap{
		"percent": func(count, total int) float64 {
			if total == 0 {
				return 0
			}
			return (float64(count) / float64(total)) * 100
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"toMillis": func(d time.Duration) float64 {
			return float64(d) / float64(time.Millisecond)
		},
	}
	tmpl, err := template.New("template.html").Funcs(funcMap).ParseFiles("templates/template.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	tmplData := struct {
		MaxDepth                int
		URLDataWithReachability []scanner.URLDataWithReachability
	}{
		MaxDepth:                depth,
		URLDataWithReachability: urlDataWithReachability,
	}

	// 创建文件
	file, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Error creating HTML file: %v", err)
	}
	defer file.Close()

	// 将模板输出到文件
	if err := tmpl.Execute(file, tmplData); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
