package tool

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type systemFile struct {
	server.ServerTool
}

func CreateFile() server.ServerTool {
	tool := systemFile{}
	tool.Tool = tool.tool()
	tool.Handler = tool.handler()
	return tool.ServerTool
}

func (t *systemFile) tool() mcp.Tool {
	t.Tool = mcp.NewTool("createFile",
		mcp.WithDescription("此工具可以将文件写入本地磁盘"),
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Description("文件名称"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("文件内容"),
			mcp.Required(),
		),
	)
	return t.Tool
}

func (t *systemFile) handler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		filename = strings.ReplaceAll(filename, "\\", "/")

		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// 拼接完整路径
		savePath := filepath.Join(WorkDir, filename)

		// 确保目录存在
		dirPath := filepath.Dir(savePath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return mcp.NewToolResultError("无法创建目录: " + err.Error()), nil
		}

		// 创建文件
		f, err := os.Create(savePath)
		if err != nil {
			return mcp.NewToolResultError("无法创建文件: " + err.Error()), nil
		}
		f.WriteString(content)
		defer f.Close()

		return mcp.NewToolResultText("文件已成功创建: " + savePath), nil
	}
}
