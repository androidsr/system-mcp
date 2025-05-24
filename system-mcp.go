package main

import (
	"fmt"
	"os"
	"system-mcp/tool"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("此工具不支持命令行参数")
		return
	}
	tool.Init(os.Args[1])

	s := server.NewMCPServer(
		"system-mcp-go",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s.AddTools(tool.CreateFile())

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("MCP Server错误: %v\n", err)
	}
}
