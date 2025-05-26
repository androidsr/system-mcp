package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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
	_, err := tool.getHtmlContent("https://products.groupdocs.app/zh/conversion/html-to-md#:~:text=%E8%BD%BB%E6%9D%BE%E5%9C%B0%E5%9C%A8%E7%BA%BF%E8%BD%AC%E6%8D%A2%E6%82%A8%E7%9A%84HTML%E6%96%87%E4%BB%B6%EF%BC%8C%E6%97%A0%E9%9C%80%E6%B3%A8%E5%86%8C%E3%80%82%20%E4%BB%8EWindows%EF%BC%8CLinux%E6%88%96MacOS%E8%BD%AC%E6%8D%A2%E6%82%A8%E7%9A%84HTML%E6%96%87%E6%A1%A3%E3%80%82%20%E5%8F%AA%E9%9C%80%E5%B0%86%E6%96%87%E4%BB%B6%E6%8B%96%E6%94%BE%E5%88%B0%E4%B8%8A%E4%BC%A0%E8%A1%A8%E5%8D%95%E4%B8%AD%EF%BC%8C%E9%80%89%E6%8B%A9%E6%89%80%E9%9C%80%E7%9A%84%E8%BE%93%E5%87%BA%E6%A0%BC%E5%BC%8F%EF%BC%8C%E7%84%B6%E5%90%8E%E5%8D%95%E5%87%BB%E2%80%9C%E8%BD%AC%E6%8D%A2%E2%80%9D%E6%8C%89%E9%92%AE%E3%80%82,%E8%BD%AC%E6%8D%A2%E5%AE%8C%E6%88%90%E5%90%8E%EF%BC%8C%E6%82%A8%E5%8F%AF%E4%BB%A5%E4%B8%8B%E8%BD%BDMD%E6%96%87%E4%BB%B6%E3%80%82%20%E6%A0%B9%E6%8D%AEHTML%E6%96%87%E4%BB%B6%E7%9A%84%E5%A4%A7%E5%B0%8F%E5%92%8C%E6%A0%BC%E5%BC%8F%EF%BC%8C%E8%BD%AC%E6%8D%A2%E5%8F%AF%E8%83%BD%E9%9C%80%E8%A6%81%E4%B8%80%E4%BA%9B%E7%89%87%E5%88%BB%EF%BC%8C%E4%BD%86%E6%98%AF%E5%A4%A7%E5%A4%9A%E6%95%B0%E6%96%87%E4%BB%B6%E4%BC%9A%E8%BF%85%E9%80%9F%E8%BD%AC%E6%8D%A2%E3%80%82%20%E5%80%9F%E5%8A%A9%E9%AB%98%E7%BA%A7%E5%8A%9F%E8%83%BD%EF%BC%8C%E6%82%A8%E5%8F%AF%E4%BB%A5%E5%B0%86%E5%8F%97%E5%AF%86%E7%A0%81%E4%BF%9D%E6%8A%A4%E7%9A%84%E6%96%87%E6%A1%A3%E8%BD%AC%E6%8D%A2%E5%B9%B6%E5%B0%86%E7%BB%93%E6%9E%9C%E5%8F%91%E9%80%81%E5%88%B0%E7%94%B5%E5%AD%90%E9%82%AE%E4%BB%B6%E3%80%82")
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(10000000)
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
