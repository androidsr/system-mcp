package tool

import (
	"context"
	"errors"
	"os"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/androidsr/sc-go/sno"
	"github.com/androidsr/sc-go/syaml"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/playwright-community/playwright-go"
)

type readHtml struct {
	server.ServerTool
}

func ReadHtml() server.ServerTool {
	sno.New(syaml.SnowflakeInfo{
		WorkerId: 2,
	})
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
		data, err := t.getHtmlContent(url)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(data), nil
	}
}

func (t *readHtml) getHtmlContent(url string) (string, error) {
	// 2. 启动 Playwright
	pw, err := playwright.Run()
	if err != nil {
		err = playwright.Install()
		if err != nil {
			return "", errors.New("安装 Playwright 失败")
		}
		return "", errors.New("无法启动 Playwright")

	}
	defer pw.Stop()

	browserPath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	showBrowser := true
	interval := 700.00 // 可设置延时毫秒

	// 3. 启动浏览器
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(browserPath),
		Headless:       playwright.Bool(!showBrowser),
		SlowMo:         playwright.Float(interval),
	})
	if err != nil {
		return "", errors.New("无法启动 Chromium 浏览器")
	}
	defer browser.Close()

	// 4. 创建上下文和页面
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads: playwright.Bool(true),
	})
	if err != nil {
		return "", errors.New("无法创建浏览器上下文")
	}

	page, err := context.NewPage()
	if err != nil {
		return "", errors.New("打开浏览器页面失败")
	}

	if _, err = page.Goto(url); err != nil {
		return "", errors.New("打开指定地址失败")
	}

	if _, err = page.Evaluate(Readability); err != nil {
		return "", errors.New("注入 JS 插件失败")
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
		return "", errors.New("自动获取文本内容失败")
	}
	// 7. 渲染为 PDF 页面
	newPage, err := context.NewPage()
	if err != nil {
		return "", errors.New("创建 PDF 页面失败")
	}
	defer newPage.Close()

	if err = newPage.SetContent(content.(string), playwright.PageSetContentOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		return "", errors.New("设置页面内容失败")
	}
	htmlContent, err := newPage.Content()
	if err != nil {
		return "", errors.New("获取html页面失败")
	}
	markdown, err := htmltomarkdown.ConvertString(htmlContent)
	if err != nil {
		return "", errors.New("转换成 Markdown 失败")
	}
	os.WriteFile("test.md", []byte(markdown), 0644)
	return markdown, nil
}
