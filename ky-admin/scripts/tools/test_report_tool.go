package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	reportRoot = "/Volumes/DataDocument/CodeSpace/AI Solution/kysion.com/backend/temp/test-reports"
)

func main() {
	// 定义命令行参数
	action := flag.String("action", "view", "操作类型: view, serve, clean")
	flag.Parse()

	switch *action {
	case "view":
		viewLatestReport()
	case "serve":
		serveReports()
	case "clean":
		cleanReports()
	default:
		fmt.Printf("未知操作: %s\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}

// 查看最新的测试报告
func viewLatestReport() {
	latestDir := filepath.Join(reportRoot, "latest")
	coverageHTML := filepath.Join(latestDir, "coverage.html")

	if _, err := os.Stat(coverageHTML); os.IsNotExist(err) {
		fmt.Println("错误: 找不到最新的覆盖率报告")
		os.Exit(1)
	}

	cmd := exec.Command("open", coverageHTML)
	if err := cmd.Run(); err != nil {
		fmt.Printf("打开覆盖率报告失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已打开覆盖率报告: %s\n", coverageHTML)
}

// 启动HTTP服务器提供报告访问
func serveReports() {
	latestDir := filepath.Join(reportRoot, "latest")
	if _, err := os.Stat(latestDir); os.IsNotExist(err) {
		fmt.Println("错误: 找不到测试报告目录")
		os.Exit(1)
	}

	fmt.Println("启动测试报告服务器...")
	fmt.Println("访问地址: http://localhost:8089/")

	cmd := exec.Command("python3", "-m", "http.server", "8089")
	cmd.Dir = latestDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("启动服务器失败: %v\n", err)
		os.Exit(1)
	}
}

// 清理测试报告
func cleanReports() {
	dirs, err := os.ReadDir(reportRoot)
	if err != nil {
		fmt.Printf("读取报告目录失败: %v\n", err)
		os.Exit(1)
	}

	for _, dir := range dirs {
		if dir.Name() == "latest" && dir.Type().IsDir() {
			// 跳过latest符号链接
			continue
		}
		path := filepath.Join(reportRoot, dir.Name())
		fmt.Printf("删除: %s\n", path)
		os.RemoveAll(path)
	}

	fmt.Println("测试报告已清理")
}
