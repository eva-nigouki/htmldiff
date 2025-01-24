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
 @Time    : 2025/1/24 -- 16:31
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: htmldiff /compare.go
*/

package htmldiff

import (
	"fmt"
	"regexp"
	"strings"
)

// CompareByWord 使用更细粒度的分词方式比较两个 string 的内容并生成差异结果
// skipPatterns 为需要跳过的正则表达式模式列表，匹配这些模式的token将不进行差异标记
func CompareByWord(prevText, latestText string, skipPatterns ...string) string {
	// 编译正则表达式
	skipRegexps := make([]*regexp.Regexp, 0, len(skipPatterns))
	for _, pattern := range skipPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			skipRegexps = append(skipRegexps, re)
		}
	}

	// 检查token是否需要跳过
	shouldSkip := func(token string) bool {
		for _, re := range skipRegexps {
			if re.MatchString(token) {
				return true
			}
		}
		return false
	}

	mergedDiffs := NewDiff().MaxCommonSubSeq(prevText, latestText)

	// 输出差异结果
	var result strings.Builder
	for _, d := range mergedDiffs {
		switch d.Typ {
		case NoDiff: // 相同
			result.WriteString(d.Content)
		case Del: // 删除
			if shouldSkip(d.Content) {
				result.WriteString(d.Content)
			} else {
				result.WriteString(fmt.Sprintf("<span class='diff-delete'>%s</span>", d.Content))
			}
		case Insert: // 插入
			if shouldSkip(d.Content) {
				result.WriteString(d.Content)
			} else {
				result.WriteString(fmt.Sprintf("<span class='diff-insert'>%s</span>", d.Content))
			}
		}
	}

	return result.String()
}
