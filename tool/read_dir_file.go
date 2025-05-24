package tool

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type readDirFile struct {
	server.ServerTool
}

func ReadDirFile() server.ServerTool {
	tool := readDirFile{}
	tool.Tool = tool.tool()
	tool.Handler = tool.handler()
	return tool.ServerTool
}

func (t *readDirFile) tool() mcp.Tool {
	t.Tool = mcp.NewTool("readDirFile",
		mcp.WithDescription("此工具可以读取指定目录以及子目录下的文件有哪些，并返回每个文件的全路径"),
		mcp.WithString("dir",
			mcp.Required(),
			mcp.Description("文件目录"),
		),
	)
	return t.Tool
}

func (t *readDirFile) handler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir, err := request.RequireString("dir")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// 判断目录是否存在
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return mcp.NewToolResultError("目录不存在: " + dir), nil
		}

		// 遍历目录下的所有文件和子目录
		var files []string
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		return mcp.NewToolResultText(strings.Join(files, " \n")), nil
	}
}
