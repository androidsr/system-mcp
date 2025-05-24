package tool

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type readFile struct {
	server.ServerTool
}

func ReadFile() server.ServerTool {
	tool := readFile{}
	tool.Tool = tool.tool()
	tool.Handler = tool.handler()
	return tool.ServerTool
}

func (t *readFile) tool() mcp.Tool {
	t.Tool = mcp.NewTool("readFile",
		mcp.WithDescription("此工具可以读取指定文件内容"),
		mcp.WithString("filePath",
			mcp.Required(),
			mcp.Description("文件名称"),
		),
	)
	return t.Tool
}

func (t *readFile) handler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, err := request.RequireString("filePath")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		bs, err := os.ReadFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(string(bs)), nil
	}
}
