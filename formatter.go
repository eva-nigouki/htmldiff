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
 @Description: htmldiff /formatter.go
*/

package htmldiff

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

// HTMLFormatter 用于格式化HTML文本
type HTMLFormatter struct {
	IndentSize int // 缩进大小
}

// NewHTMLFormatter 创建一个新的HTMLFormatter实例
func NewHTMLFormatter(indentSize int) *HTMLFormatter {
	return &HTMLFormatter{
		IndentSize: indentSize,
	}
}

// Format 格式化HTML文本
func (f *HTMLFormatter) Format(content string) (string, error) {
	// 解析HTML文档
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	// 创建格式化输出缓冲区
	var buf bytes.Buffer

	// 递归格式化节点
	f.formatNode(&buf, doc, 0)

	return buf.String(), nil
}

// formatNode 递归格式化节点
func (f *HTMLFormatter) formatNode(buf *bytes.Buffer, n *html.Node, depth int) {
	switch n.Type {
	case html.DocumentNode:
		// 处理文档节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f.formatNode(buf, c, depth)
		}

	case html.ElementNode:
		// 添加缩进
		f.writeIndent(buf, depth)

		// 写入开始标签
		buf.WriteByte('<')
		buf.WriteString(n.Data)
		f.writeAttributes(buf, n.Attr)
		buf.WriteByte('>')
		buf.WriteByte('\n')

		// 处理子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f.formatNode(buf, c, depth+1)
		}

		// 如果有子节点，添加结束标签的缩进
		if n.FirstChild != nil {
			f.writeIndent(buf, depth)
		}

		// 写入结束标签
		buf.WriteString("</")
		buf.WriteString(n.Data)
		buf.WriteByte('>')
		buf.WriteByte('\n')

	case html.TextNode:
		// 处理文本节点，去除多余的空白字符
		text := strings.TrimSpace(n.Data)
		if text != "" {
			f.writeIndent(buf, depth)
			buf.WriteString(text)
			buf.WriteByte('\n')
		}

	case html.CommentNode:
		// 处理注释节点
		f.writeIndent(buf, depth)
		buf.WriteString("<!--")
		buf.WriteString(n.Data)
		buf.WriteString("-->\n")
	}
}

// writeIndent 写入缩进
func (f *HTMLFormatter) writeIndent(buf *bytes.Buffer, depth int) {
	for i := 0; i < depth*f.IndentSize; i++ {
		buf.WriteByte(' ')
	}
}

// writeAttributes 写入属性
func (f *HTMLFormatter) writeAttributes(buf *bytes.Buffer, attrs []html.Attribute) {
	for _, attr := range attrs {
		buf.WriteByte(' ')
		buf.WriteString(attr.Key)
		buf.WriteString(`="`)
		buf.WriteString(attr.Val)
		buf.WriteByte('"')
	}
}
