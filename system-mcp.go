package main

import (
	"fmt"
	"os"
	"system-mcp/tool"

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
