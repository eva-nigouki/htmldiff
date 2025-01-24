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
 @Time    : 2025/1/24 -- 11:26
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: htmldiff /utils.go
*/

package htmldiff

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// FetchHTML 从URL获取HTML内容
func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取URL内容失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %v", err)
	}

	return string(body), nil
}

// attributesEqual 比较两个节点的属性是否相同
func attributesEqual(attr1, attr2 []html.Attribute) bool {
	if len(attr1) != len(attr2) {
		return false
	}

	for i := range attr1 {
		if attr1[i].Key != attr2[i].Key || attr1[i].Val != attr2[i].Val {
			return false
		}
	}

	return true
}

// renderNode 将节点渲染为HTML字符串
func renderNode(n *html.Node) string {
	var result strings.Builder
	if n.Type == html.ElementNode {
		result.WriteString(fmt.Sprintf("<%s%s>", n.Data, renderAttributes(n.Attr)))
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result.WriteString(renderNode(c))
		}
		result.WriteString(fmt.Sprintf("</%s>", n.Data))
	} else if n.Type == html.TextNode {
		// 处理文本节点中的空白字符
		trimmed := strings.TrimSpace(n.Data)
		if trimmed != "" {
			result.WriteString(trimmed)
		} else if n.Data != "" {
			// 保留一个空格，以保持节点间的分隔
			result.WriteString(" ")
		}
	}
	return result.String()
}

// renderAttributes 将属性列表渲染为字符串
func renderAttributes(attrs []html.Attribute) string {
	var result strings.Builder
	for _, attr := range attrs {
		result.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Key, attr.Val))
	}
	return result.String()
}

// hasOnlyTextContent 检查节点是否只包含文本内容
func hasOnlyTextContent(n *html.Node) bool {
	if n.FirstChild == nil {
		return true
	}
	return n.FirstChild.Type == html.TextNode && n.FirstChild.NextSibling == nil
}

// getTextContent 获取节点的文本内容
func getTextContent(n *html.Node) string {
	if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
		return n.FirstChild.Data
	}
	return ""

}

// escapeHTML 将HTML特殊符号转义为实体符号
func escapeHTML(content string) string {
	return html.EscapeString(content)
}

func isUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
