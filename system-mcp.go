package main

import (
	"fmt"
	"os"
	"system-mcp/tool"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/playwright-community/playwright-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("此工具不支持命令行参数")
		return
	}
	tool.Init(os.Args[1])

	// 2. 启动 Playwright
	pw, err := playwright.Run()
	if err != nil {
		fmt.Printf("无法启动 Playwright: %v\n", err)
		return
	}
	defer pw.Stop()

	browserPath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	showBrowser := true
	interval := 700.00 // 可设置延时毫秒

	// 3. 启动浏览器
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(browserPath),
		Headless:       playwright.Bool(!showBrowser),
		SlowMo:         playwright.Float(interval),
	})
	if err != nil {
		return
	}
	defer browser.Close()
	fmt.Println("启动浏览器成功")
	time.Sleep(20 * time.Second)
	s := server.NewMCPServer(
		"system-mcp-go",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s.AddTools(tool.CreateFile(), tool.ReadDirFile(), tool.CopyFile(), tool.ReadFile(), tool.ReadHtml())

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("MCP Server错误: %v\n", err)
	}
}
