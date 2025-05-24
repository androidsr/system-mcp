package tool

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type copyFile struct {
	server.ServerTool
}

func CopyFile() server.ServerTool {
	tool := copyFile{}
	tool.Tool = tool.tool()
	tool.Handler = tool.handler()
	return tool.ServerTool
}

func (t *copyFile) tool() mcp.Tool {
	t.Tool = mcp.NewTool("copyFile",
		mcp.WithDescription("此工具可以将文件复制到指定目录"),
		mcp.WithString("sourceFile",
			mcp.Required(),
			mcp.Description("源文件路径"),
		),
		mcp.WithString("targetFile",
			mcp.Required(),
			mcp.Description("目标文件路径"),
		),
	)
	return t.Tool
}

func (t *copyFile) handler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceFile, err := request.RequireString("sourceFile")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		sourceFile = strings.ReplaceAll(sourceFile, "\\", "/")

		targetFile, err := request.RequireString("targetFile")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		targetFile = strings.ReplaceAll(targetFile, "\\", "/")
		fmt.Println("sourceFile: " + sourceFile)
		fmt.Println("targetFile: " + targetFile)
		time.Sleep(1000)
		//如果文件已经存在，则忽略该文件
		if _, err = os.Stat(targetFile); err == nil {
			return mcp.NewToolResultText("文件已经存在: " + targetFile), nil
		}

		// 判断目录是否存在, 如果不存在则创建
		dirPath := filepath.Dir(targetFile)
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return mcp.NewToolResultError("无法创建目录: " + err.Error()), nil
		}

		// 复制文件
		err = Copy(sourceFile, targetFile)
		if err != nil {
			return mcp.NewToolResultError("复制文件失败: " + err.Error()), nil
		}
		return mcp.NewToolResultText("文件复制成功: " + targetFile + " \n"), nil
	}
}

func Copy(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}

	return nil
}
