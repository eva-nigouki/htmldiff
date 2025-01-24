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
 @Time    : 2025/1/24 -- 11:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: htmldiff /diff.go
*/

package htmldiff

import (
	"strings"
	"unicode"
)

type DiffType int

const (
	NoDiff DiffType = 0
	Del    DiffType = 1
	Insert DiffType = 2
)

type DiffImpl struct{}

func NewDiff() *DiffImpl {
	return &DiffImpl{}
}

type DiffByToken struct {
	Start, End int
	Typ        DiffType
}

type DiffItem struct {
	Content string
	Typ     DiffType
}

// MaxCommonSubSeq 最长公共子序列
func (d *DiffImpl) MaxCommonSubSeq(prevText, latestText string) []DiffItem {
	// 将文本分解为词元
	prevTokens := d.tokenize(prevText)
	latestTokens := d.tokenize(latestText)

	// 使用动态规划计算最长公共子序列
	m, n := len(prevTokens), len(latestTokens)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 填充DP表
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if prevTokens[i-1] == latestTokens[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// 根据DP表回溯找出差异区间
	var diffs []DiffByToken // typ: 0-相同-NoDiff, 1-删除-Insert, 2-插入-Del
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && prevTokens[i-1] == latestTokens[j-1] {
			if len(diffs) == 0 || diffs[len(diffs)-1].Typ != 0 {
				diffs = append(diffs, DiffByToken{i - 1, i, NoDiff})
			} else {
				diffs[len(diffs)-1].Start = i - 1
			}
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			if len(diffs) == 0 || diffs[len(diffs)-1].Typ != 2 {
				diffs = append(diffs, DiffByToken{j - 1, j, Insert})
			} else {
				diffs[len(diffs)-1].Start = j - 1
			}
			j--
		} else {
			if len(diffs) == 0 || diffs[len(diffs)-1].Typ != 1 {
				diffs = append(diffs, DiffByToken{i - 1, i, Del})
			} else {
				diffs[len(diffs)-1].Start = i - 1
			}
			i--
		}
	}

	// 反转差异数组，使其按照文本顺序输出
	for i, j := 0, len(diffs)-1; i < j; i, j = i+1, j-1 {
		diffs[i], diffs[j] = diffs[j], diffs[i]
	}

	// 合并相近的差异区间
	var mergedDiffs []DiffByToken
	for i := 0; i < len(diffs); i++ {
		current := diffs[i]
		// 查找可以合并的区间
		for j := i + 1; j < len(diffs); j++ {
			next := diffs[j]
			// 如果两个区间类型不同，则停止合并
			if next.Typ != current.Typ {
				break
			}
			// 检查是否为HTML转义字符
			isCurrentEscape := isHTMLEscapeSequence(prevTokens[current.Start:current.End])
			isNextEscape := isHTMLEscapeSequence(prevTokens[next.Start:next.End])
			// 如果两个区间都是HTML转义字符，或者距离小于10个token，则合并
			if (isCurrentEscape && isNextEscape) || next.Start-current.End <= 10 {
				// 合并区间
				current.End = next.End
				i = j
			} else {
				break
			}
		}
		mergedDiffs = append(mergedDiffs, current)
	}

	// 输出差异结果
	var res []DiffItem
	for _, d := range mergedDiffs {
		switch d.Typ {
		case NoDiff: // 相同
			for i := d.Start; i < d.End; i++ {
				res = append(res, DiffItem{Typ: d.Typ, Content: prevTokens[i]})
			}
		case 1: // 删除
			for i := d.Start; i < d.End; i++ {
				res = append(res, DiffItem{Typ: d.Typ, Content: prevTokens[i]})
			}
		case 2: // 插入
			for i := d.Start; i < d.End; i++ {
				res = append(res, DiffItem{Typ: d.Typ, Content: latestTokens[i]})
			}
		}
	}
	return res
}

// tokenize 将文本分解为词元（英文单词、标点符号之间的英文字符序列或中文字符序列）
func (*DiffImpl) tokenize(text string) []string {
	var tokens []string
	var currentToken strings.Builder
	var lastType rune
	var escapeBuffer strings.Builder

	// 辅助函数：获取字符类型
	getCharType := func(r rune) rune {
		switch {
		case unicode.Is(unicode.Han, r):
			return 'c' // 中文字符
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			return 'e' // 英文字母和数字
		case unicode.IsSpace(r):
			return 's' // 空白字符
		default:
			return 'p' // 标点符号和其他字符
		}
	}

	// 辅助函数：检查是否为HTML转义字符的开始
	isEscapeStart := func(r rune) bool {
		return r == '&'
	}

	// 辅助函数：检查是否为HTML转义字符的结束
	isEscapeEnd := func(r rune) bool {
		return r == ';'
	}

	// 辅助函数：检查当前缓冲区是否构成有效的HTML转义字符
	isValidEscapeSequence := func(s string) bool {
		escapeSequences := []string{
			// 基本HTML实体
			"&amp;", "&lt;", "&gt;", "&quot;", "&apos;", "&#",
			// 常用符号
			"&nbsp;", "&copy;", "&reg;", "&trade;", "&sect;", "&deg;",
			// 货币符号
			"&cent;", "&pound;", "&euro;", "&yen;",
			// 数学符号
			"&plusmn;", "&times;", "&divide;", "&ne;", "&le;", "&ge;", "&infin;",
			// 箭头符号
			"&larr;", "&uarr;", "&rarr;", "&darr;",
			// 其他特殊字符
			"&middot;", "&bull;", "&hellip;", "&prime;", "&Prime;",
			// 常用变音符号
			"&acute;", "&cedil;", "&uml;", "&macr;",
		}
		for _, seq := range escapeSequences {
			if strings.HasPrefix(s, seq) {
				return true
			}
		}
		return false
	}

	// 辅助函数：添加当前token到结果中
	addToken := func() {
		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		}
	}

	// 辅助函数：添加HTML转义字符到结果中
	addEscapeToken := func() {
		if escapeBuffer.Len() > 0 {
			tokens = append(tokens, escapeBuffer.String())
			escapeBuffer.Reset()
		}
	}

	for _, r := range text {
		currentType := getCharType(r)

		// 如果正在处理HTML转义字符
		if escapeBuffer.Len() > 0 {
			escapeBuffer.WriteRune(r)
			if isEscapeEnd(r) {
				// 如果是有效的HTML转义字符，添加为独立的token
				if isValidEscapeSequence(escapeBuffer.String()) {
					addToken() // 先添加之前的token
					addEscapeToken()
					lastType = 'p'
					continue
				}
			}
			// 如果转义字符缓冲区过长，说明不是有效的转义字符
			if escapeBuffer.Len() > 10 {
				// 将缓冲区内容作为普通字符处理
				currentToken.WriteString(escapeBuffer.String())
				escapeBuffer.Reset()
			}
			continue
		}

		// 检查是否开始新的HTML转义字符
		if isEscapeStart(r) {
			escapeBuffer.WriteRune(r)
			continue
		}

		// 处理空白字符
		if currentType == 's' {
			// 如果前一个字符是字母或数字，添加当前token
			if lastType == 'e' {
				addToken()
			}
			// 将空白字符作为独立的token
			tokens = append(tokens, string(r))
		} else if currentType == 'p' {
			// 如果前面有未处理的token，先添加
			if currentToken.Len() > 0 {
				addToken()
			}
			// 将标点符号作为独立的token
			tokens = append(tokens, string(r))
		} else {
			// 如果当前是字母或数字，但前一个是中文，先添加当前token
			if currentType == 'e' && lastType == 'c' {
				addToken()
			}
			// 如果当前是中文，但前一个是字母或数字，先添加当前token
			if currentType == 'c' && lastType == 'e' {
				addToken()
			}
			// 添加当前字符到token
			currentToken.WriteRune(r)
		}

		lastType = currentType
	}

	// 添加最后一个token
	addToken()

	// 处理未完成的转义字符
	if escapeBuffer.Len() > 0 {
		currentToken.WriteString(escapeBuffer.String())
		addToken()
	}

	return tokens
}

// isHTMLEscapeSequence 检查给定的token序列是否为HTML转义字符
func isHTMLEscapeSequence(tokens []string) bool {
	if len(tokens) == 0 {
		return false
	}
	// 将tokens合并为一个字符串
	content := strings.Join(tokens, "")
	// 检查是否为常见的HTML转义字符
	escapeSequences := []string{"&amp;", "&lt;", "&gt;", "&quot;", "&apos;", "&#"}
	for _, seq := range escapeSequences {
		if strings.Contains(content, seq) {
			return true
		}
	}
	return false
}
