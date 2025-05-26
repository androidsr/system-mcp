package tool

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/androidsr/sc-go/sno"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/playwright-community/playwright-go"
	"rsc.io/pdf"
)

type readHtml struct {
	server.ServerTool
}

func ReadHtml() server.ServerTool {
	tool := readHtml{}
	tool.Tool = tool.tool()
	tool.Handler = tool.handler()
	return tool.ServerTool
}

func (t *readHtml) tool() mcp.Tool {
	t.Tool = mcp.NewTool("readHtml",
		mcp.WithDescription("此工具可以读取指定网页的正文内容"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("网页地址"),
		),
	)
	return t.Tool
}

func (t *readHtml) handler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 1. 获取 URL 参数
		url, err := request.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		fmt.Println("处理 URL:", url)

		// 2. 启动 Playwright

		pw, err := playwright.Run()
		if err != nil {
			err = playwright.Install()
			if err != nil {
				return nil, fmt.Errorf("安装 Playwright 失败: %v", err)
			}
			return nil, fmt.Errorf("无法启动 Playwright: %v", err)
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
			return nil, fmt.Errorf("无法启动 Chromium 浏览器: %v", err)
		}
		defer browser.Close()

		// 4. 创建上下文和页面
		context, err := browser.NewContext(playwright.BrowserNewContextOptions{
			AcceptDownloads: playwright.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("无法创建浏览器上下文: %v", err)
		}

		page, err := context.NewPage()
		if err != nil {
			return mcp.NewToolResultError("打开浏览器页面失败"), nil
		}

		if _, err = page.Goto(url); err != nil {
			return mcp.NewToolResultError("打开指定地址失败"), nil
		}

		if _, err = page.Evaluate(Readability); err != nil {
			return mcp.NewToolResultError("注入 JS 插件失败"), nil
		}
		// 6. 提取正文内容
		content, err := page.Evaluate(`() => {
			try {
				const article = new Readability(document.cloneNode(true)).parse();
				return article?.content || '未提取到内容';
			} catch (e) {
				return 'Readability 提取失败: ' + e.toString();
			}
		}`)
		if err != nil {
			return mcp.NewToolResultError("自动获取文本内容失败"), nil
		}

		// 7. 渲染为 PDF 页面
		newPage, err := context.NewPage()
		if err != nil {
			return mcp.NewToolResultError("创建 PDF 页面失败"), nil
		}
		defer newPage.Close()

		if err = newPage.SetContent(content.(string), playwright.PageSetContentOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); err != nil {
			return mcp.NewToolResultError("设置页面内容失败"), nil
		}

		// 8. 保存 PDF
		savePath := path.Join(WorkDir, "temp", time.Now().Format("2006-01-02"))
		if err = os.MkdirAll(savePath, 0755); err != nil {
			return mcp.NewToolResultError("创建目录失败"), nil
		}
		pdfPath := path.Join(savePath, sno.GetString()+".pdf")
		data, err := newPage.PDF(playwright.PagePdfOptions{
			Path:                playwright.String(pdfPath),
			Format:              playwright.String("A4"),
			DisplayHeaderFooter: playwright.Bool(false),
			PrintBackground:     playwright.Bool(false),
		})
		if err != nil {
			return mcp.NewToolResultError("生成 PDF 文件失败"), nil
		}

		// 9. 读取 PDF 文本
		r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return mcp.NewToolResultError("PDF 读取失败"), nil
		}

		var mergedTexts []string
		for i := 1; i <= r.NumPage(); i++ {
			p := r.Page(i)
			if p.V.IsNull() {
				continue
			}
			var mergedSentence string
			var lastY float64
			for _, text := range p.Content().Text {
				if text.Y == lastY {
					mergedSentence += text.S
				} else {
					if mergedSentence != "" {
						mergedTexts = append(mergedTexts, mergedSentence)
					}
					mergedSentence = text.S
					lastY = text.Y
				}
			}
			if mergedSentence != "" {
				mergedTexts = append(mergedTexts, mergedSentence)
			}
		}
		// 10. 返回合并结果
		mergedText := strings.Join(mergedTexts, "\n")
		return mcp.NewToolResultText(mergedText), nil
	}
}
