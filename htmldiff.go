/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2025/1/24 -- 16:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: htmldiff /htmldiff.go
*/

package htmldiff

import (
	"fmt"
	"html"
	"os"
	"strings"
)

// HtmlDiff 结构体用于存储HTML直接差异比对的配置和结果
type HtmlDiff struct {
	Previous  string
	Latest    string
	formatter *HTMLFormatter
}

// NewHtmlDiff 创建一个新的HtmlDiff实例
func NewHtmlDiff(previous, latest string) *HtmlDiff {
	return &HtmlDiff{
		Previous:  previous,
		Latest:    latest,
		formatter: NewHTMLFormatter(2),
	}
}

// getHTMLContent 处理HTML内容，如果是URL则获取内容，否则直接返回
func (h *HtmlDiff) getHTMLContent(content string) (string, error) {
	if strings.HasPrefix(strings.ToLower(content), "http://") ||
		strings.HasPrefix(strings.ToLower(content), "https://") {
		html, err := FetchHTML(content)
		if err != nil {
			return "", err
		}
		return html, nil
	}
	return content, nil
}

// escapeHTML 将HTML特殊符号转义为实体符号
func (h *HtmlDiff) escapeHTML(content string) string {
	return html.EscapeString(content)
}

// Compare 比较两个HTML内容并生成差异结果
func (h *HtmlDiff) Compare() (string, error) {
	// 处理Previous内容
	prevHTML, err := h.getHTMLContent(h.Previous)
	if err != nil {
		return "", err
	}

	// 处理Latest内容
	latestHTML, err := h.getHTMLContent(h.Latest)
	if err != nil {
		return "", err
	}

	formatter := NewHTMLFormatter(2)
	// 格式化HTML内容
	prevHTML, err = formatter.Format(prevHTML)
	if err != nil {
		return "", fmt.Errorf("格式化Previous HTML失败: %v", err)
	}

	latestHTML, err = formatter.Format(latestHTML)
	if err != nil {
		return "", fmt.Errorf("格式化Latest HTML失败: %v", err)
	}

	// HTML特殊符号转义
	prevHTML = h.escapeHTML(prevHTML)
	latestHTML = h.escapeHTML(latestHTML)

	// 使用CompareScriptContentByWord进行比对
	diffResult := CompareByWord(prevHTML, latestHTML, "0x[0-9a-fA-F]+")

	// 构建HTML格式的差异结果
	result := fmt.Sprintf("<div class='diff-container'><pre><code>%s</code></pre></div>", diffResult)
	return result, nil
}

// SaveToFile 将HTML差异比对结果保存到文件
func (h *HtmlDiff) SaveToFile(filePath string) error {
	// 获取差异比对结果
	diffResult, err := h.Compare()
	if err != nil {
		return fmt.Errorf("生成差异比对结果失败: %v", err)
	}

	// 构建完整的HTML文档
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>HTML Diff Result</title>
	<style>
		.diff-container {
			padding: 20px;
			font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
			line-height: 1.5;
			font-size: 14px;
		}
		.diff-delete {
			background-color: #ffeef0;
			color: #b31d28;
			text-decoration: line-through;
			padding: 2px 0;
		}
		.diff-insert {
			background-color: #e6ffed;
			color: #22863a;
			padding: 2px 0;
		}
		pre {
			margin: 0;
			padding: 16px;
			background-color: #f6f8fa;
			border-radius: 6px;
			overflow-x: auto;
		}
		code {
			display: inline-block;
			min-width: 100%%;
			white-space: pre-wrap;
			word-break: break-all;
		}
	</style>
</head>
<body>
%s
</body>
</html>`, diffResult)

	// 写入文件
	err = os.WriteFile(filePath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
