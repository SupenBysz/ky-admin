package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 定义全局命令行参数
var (
	// 通用参数
	mode = flag.String("mode", "summary", "操作模式: summary或server")
	dir  = flag.String("dir", "", "报告目录")

	// summary模式参数
	outputFile = flag.String("output", "", "输出HTML文件路径")
	inputFile  = flag.String("input", "", "测试输出日志文件路径")

	// server模式参数
	port = flag.Int("port", 8089, "服务器端口")
)

// Test结果结构体
type TestResult struct {
	Package   string
	TestName  string
	Status    string // PASS, FAIL, SKIP
	Duration  string
	Output    []string
	Error     string
	Timestamp time.Time
}

// -- 测试摘要生成器功能 --

// 解析测试输出
func parseTestOutput(content string) []TestResult {
	var results []TestResult
	var currentTest *TestResult

	// 正则表达式匹配测试行
	testStartRegex := regexp.MustCompile(`^=== RUN\s+(.+)$`)
	testEndRegex := regexp.MustCompile(`^\s*--- (PASS|FAIL|SKIP):\s+(.+)\s+\((.+)\)$`)
	packageRegex := regexp.MustCompile(`^(?:ok|FAIL)\s+(.+)\s+([0-9.]+s|\[.*?\])(?:\s+coverage:\s+(.+)%)?$`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// 匹配测试开始
		if matches := testStartRegex.FindStringSubmatch(line); len(matches) > 0 {
			testName := matches[1]
			if currentTest != nil && currentTest.TestName == testName {
				continue
			}

			currentTest = &TestResult{
				TestName:  testName,
				Output:    []string{},
				Timestamp: time.Now(),
			}
			continue
		}

		// 匹配测试结束
		if matches := testEndRegex.FindStringSubmatch(line); len(matches) > 0 {
			if currentTest == nil {
				continue
			}

			status := matches[1]
			testName := matches[2]
			duration := matches[3]

			// 确保这是对应的测试结束
			if !strings.HasSuffix(currentTest.TestName, testName) {
				continue
			}

			currentTest.Status = status
			currentTest.Duration = duration

			results = append(results, *currentTest)
			currentTest = nil
			continue
		}

		// 收集当前测试的输出
		if currentTest != nil {
			currentTest.Output = append(currentTest.Output, line)
			if strings.Contains(line, "Error:") || strings.Contains(line, "Failure:") {
				currentTest.Error += line + "\n"
			}
		}

		// 匹配包结果
		if matches := packageRegex.FindStringSubmatch(line); len(matches) > 0 {
			packageName := matches[1]
			duration := matches[2]

			// 将包结果也添加到测试结果中
			var status string
			if strings.HasPrefix(matches[0], "ok") {
				status = "PASS"
			} else {
				status = "FAIL"
			}
			results = append(results, TestResult{
				Package:   packageName,
				TestName:  "（包总结）",
				Status:    status,
				Duration:  duration,
				Timestamp: time.Now(),
			})
		}
	}

	return results
}

// 生成HTML测试摘要
func generateHTMLSummary(results []TestResult) string {
	// 计算统计数据
	var passed, failed, skipped int
	for _, result := range results {
		switch result.Status {
		case "PASS":
			passed++
		case "FAIL":
			failed++
		case "SKIP":
			skipped++
		}
	}

	// 生成HTML
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>测试摘要报告</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1, h2 {
            color: #333;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }
        .summary {
            display: flex;
            margin: 20px 0;
            font-size: 18px;
        }
        .summary-item {
            flex: 1;
            text-align: center;
            padding: 15px;
            border-radius: 5px;
            margin: 0 10px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .pass {
            background-color: #e8f5e9;
            color: #2e7d32;
        }
        .fail {
            background-color: #ffebee;
            color: #c62828;
        }
        .skip {
            background-color: #fff8e1;
            color: #f9a825;
        }
        .test-list {
            list-style: none;
            padding: 0;
        }
        .test-item {
            border-bottom: 1px solid #eee;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 5px;
        }
        .test-item.pass {
            border-left: 5px solid #4caf50;
        }
        .test-item.fail {
            border-left: 5px solid #f44336;
        }
        .test-item.skip {
            border-left: 5px solid #ffeb3b;
        }
        .test-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .test-name {
            font-weight: bold;
            font-size: 16px;
            flex-grow: 1;
        }
        .test-badge {
            padding: 5px 10px;
            border-radius: 3px;
            color: white;
            font-size: 14px;
        }
        .test-badge.pass {
            background-color: #4caf50;
        }
        .test-badge.fail {
            background-color: #f44336;
        }
        .test-badge.skip {
            background-color: #ffeb3b;
            color: #333;
        }
        .test-duration {
            color: #757575;
            font-size: 14px;
            margin-left: 10px;
        }
        .test-package {
            color: #757575;
            font-size: 14px;
            margin-bottom: 5px;
        }
        .error-output {
            background-color: #ffebee;
            padding: 10px;
            border-radius: 3px;
            color: #b71c1c;
            white-space: pre-wrap;
            margin-top: 10px;
            font-family: monospace;
        }
        .test-output {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 3px;
            white-space: pre-wrap;
            margin-top: 10px;
            font-family: monospace;
            max-height: 300px;
            overflow-y: auto;
        }
        .timestamp {
            color: #757575;
            font-size: 12px;
            float: right;
        }
        .toggle-button {
            background: none;
            border: none;
            cursor: pointer;
            color: #2196f3;
            font-size: 14px;
            padding: 0;
            margin-left: 10px;
        }
        .hidden {
            display: none;
        }
    </style>
    <script>
        function toggleOutput(id) {
            var element = document.getElementById(id);
            var button = document.getElementById('btn-'+id);
            if (element.classList.contains('hidden')) {
                element.classList.remove('hidden');
                button.textContent = '隐藏输出';
            } else {
                element.classList.add('hidden');
                button.textContent = '显示输出';
            }
        }
    </script>
</head>
<body>
    <div class="container">
        <h1>测试摘要报告</h1>
        
        <div class="summary">
            <div class="summary-item pass">
                <div>通过</div>
                <div>%d</div>
            </div>
            <div class="summary-item fail">
                <div>失败</div>
                <div>%d</div>
            </div>
            <div class="summary-item skip">
                <div>跳过</div>
                <div>%d</div>
            </div>
        </div>
        
        <h2>测试详情</h2>
        <ul class="test-list">
`, passed, failed, skipped)

	// 添加测试结果
	for i, result := range results {
		statusClass := strings.ToLower(result.Status)
		outputID := fmt.Sprintf("output-%d", i)

		html += fmt.Sprintf(`
        <li class="test-item %s">
            <div class="test-header">
                <span class="test-badge %s">%s</span>
                <span class="test-name">%s</span>
                <span class="test-duration">%s</span>`,
			statusClass, statusClass, result.Status, result.TestName, result.Duration)

		if len(result.Output) > 0 {
			html += fmt.Sprintf(`
                <button id="btn-%s" class="toggle-button" onclick="toggleOutput('%s')">显示输出</button>`,
				outputID, outputID)
		}

		html += `
            </div>`

		if result.Package != "" {
			html += fmt.Sprintf(`
            <div class="test-package">包: %s</div>`, result.Package)
		}

		if result.Error != "" {
			html += fmt.Sprintf(`
            <div class="error-output">%s</div>`, result.Error)
		}

		if len(result.Output) > 0 {
			html += fmt.Sprintf(`
            <div id="%s" class="test-output hidden">%s</div>`,
				outputID, strings.Join(result.Output, "\n"))
		}

		html += `
            <div class="timestamp">` + result.Timestamp.Format("2006-01-02 15:04:05") + `</div>
        </li>`
	}

	html += `
        </ul>
    </div>
</body>
</html>
`
	return html
}

// 总结测试生成执行函数
func generateSummary() {
	var (
		inputContent string
		outputPath   string
	)

	// 检查输入文件
	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法读取输入文件 %s: %v\n", *inputFile, err)
			os.Exit(1)
		}
		inputContent = string(data)
	} else {
		fmt.Fprintln(os.Stderr, "请使用 -input 指定输入文件")
		os.Exit(1)
	}

	// 确定输出路径
	if *outputFile != "" {
		outputPath = *outputFile
	} else if *dir != "" {
		// 在报告目录中创建
		outputPath = filepath.Join(*dir, "test_summary.html")
	} else {
		// 在当前目录创建
		outputPath = "test_summary.html"
	}

	// 解析测试输出
	results := parseTestOutput(inputContent)

	// 生成HTML摘要
	html := generateHTMLSummary(results)

	// 写入文件
	err := os.WriteFile(outputPath, []byte(html), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法写入输出文件 %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("测试摘要已生成: %s\n", outputPath)
}

// -- 报告服务器功能 --

// 获取所有测试报告目录
func getReportDirs(baseDir string) ([]string, error) {
	var dirs []string
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// 生成首页HTML
func generateIndexHTML(reportDirs []string, baseDir string) string {
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>测试报告中心</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }
        .report-list {
            list-style: none;
            padding: 0;
        }
        .report-item {
            border-bottom: 1px solid #eee;
            padding: 10px 0;
        }
        .report-date {
            font-weight: bold;
            color: #555;
        }
        .report-link {
            display: inline-block;
            margin: 5px 10px 5px 0;
            padding: 8px 15px;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
            border-radius: 3px;
        }
        .report-link:hover {
            background-color: #45a049;
        }
        .no-reports {
            color: #999;
            font-style: italic;
            padding: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>KY-Admin 测试报告中心</h1>
`

	if len(reportDirs) == 0 {
		html += `<p class="no-reports">目前没有测试报告</p>`
	} else {
		html += `<ul class="report-list">`

		// 按时间排序（最新的在前面）
		for i := len(reportDirs) - 1; i >= 0; i-- {
			dir := reportDirs[i]
			// 从目录名解析时间（假设格式为 YYYYMMDD_HHMMSS）
			date := "未知日期"
			if len(dir) == 15 && strings.Contains(dir, "_") {
				t, err := time.Parse("20060102_150405", dir)
				if err == nil {
					date = t.Format("2006-01-02 15:04:05")
				}
			}

			html += fmt.Sprintf(`
            <li class="report-item">
                <div class="report-date">%s</div>
                <div>
                    <a href="/%s/coverage.html" class="report-link">覆盖率报告</a>
                    <a href="/%s/test_summary.html" class="report-link">测试摘要</a>
                </div>
            </li>`, date, dir, dir)
		}

		html += `</ul>`
	}

	html += `
    </div>
</body>
</html>
`
	return html
}

// 服务器启动函数
func startServer() {
	// 如果没有指定目录，使用默认的test_reports
	if *dir == "" {
		*dir = "test_reports"
	}

	// 获取项目根目录
	projectRoot, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))
	if err != nil {
		log.Fatalf("无法确定项目根目录: %v", err)
	}

	reportsPath := filepath.Join(projectRoot, *dir)

	// 确保报告目录存在
	if _, err := os.Stat(reportsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(reportsPath, 0755); err != nil {
			log.Fatalf("无法创建报告目录: %v", err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			dirs, err := getReportDirs(reportsPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("无法读取报告目录: %v", err), http.StatusInternalServerError)
				return
			}

			html := generateIndexHTML(dirs, reportsPath)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
			return
		}

		// 为测试报告文件提供服务
		http.FileServer(http.Dir(reportsPath)).ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("测试报告服务器启动在 http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// 程序入口
func main() {
	flag.Parse()

	// 根据模式选择不同的功能
	switch *mode {
	case "summary":
		generateSummary()
	case "server":
		startServer()
	default:
		fmt.Fprintf(os.Stderr, "未知的模式: %s\n", *mode)
		fmt.Fprintf(os.Stderr, "有效的模式: summary, server\n")
		flag.Usage()
		os.Exit(1)
	}
}
