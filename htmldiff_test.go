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
 @Time    : 2025/1/24 -- 16:36
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: htmldiff /htmldiff_test.go
*/

package htmldiff

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewHTMLDiff(t *testing.T) {
	prev := "<p>Old</p>"
	latest := "<p>New</p>"
	diff := NewHtmlDiff(prev, latest)

	if diff.Previous != prev {
		t.Errorf("Previous content not set correctly, got: %v, want: %v", diff.Previous, prev)
	}
	if diff.Latest != latest {
		t.Errorf("Latest content not set correctly, got: %v, want: %v", diff.Latest, latest)
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"HTTP URL", "http://example.com", true},
		{"HTTPS URL", "https://example.com", true},
		{"Non URL", "<p>Hello</p>", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUrl(tt.url); got != tt.want {
				t.Errorf("isURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Test Content</body></html>"))
	}))
	defer server.Close()

	html, err := FetchHTML(server.URL)
	if err != nil {
		t.Errorf("FetchHTML() error = %v", err)
	}
	if !strings.Contains(html, "Test Content") {
		t.Errorf("FetchHTML() content not correct, got: %v", html)
	}

	_, err = FetchHTML("http://invalid-url")
	if err == nil {
		t.Error("FetchHTML() should return error for invalid URL")
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		prev     string
		latest   string
		contains []string
	}{
		{
			name:   "Simple text change",
			prev:   "<html><body><p>Old text</p></body></html>",
			latest: "<html><body><p>New text</p></body></html>",
			contains: []string{
				"<del class='diff-delete'>Old text</del>",
				"<ins class='diff-insert'>New text</ins>",
			},
		},
		{
			name:   "No changes",
			prev:   "<html><body><p>Same content</p></body></html>",
			latest: "<html><body><p>Same content</p></body></html>",
			contains: []string{
				"<p>Same content</p>",
			},
		},
		{
			name:   "Multiple changes",
			prev:   "<html><body><h1>Title</h1><p>First</p><p>Second</p></body></html>",
			latest: "<html><body><h1>New Title</h1><p>First</p><p>Changed</p></body></html>",
			contains: []string{
				"<del class='diff-delete'>Title</del>",
				"<ins class='diff-insert'>New Title</ins>",
				"<p>First</p>",
				"<del class='diff-delete'>Second</del>",
				"<ins class='diff-insert'>Changed</ins>",
			},
		},
		{
			name:   "Attribute changes",
			prev:   "<html><body><div class=\"old\" id=\"test\">Content</div></body></html>",
			latest: "<html><body><div class=\"new\" id=\"test\">Content</div></body></html>",
			contains: []string{
				"<del class='diff-delete'>",
				"class=\"old\"",
				"<ins class='diff-insert'>",
				"class=\"new\"",
			},
		},
		{
			name:   "Node addition and deletion",
			prev:   "<html><body><div><p>First</p></div></body></html>",
			latest: "<html><body><div><p>First</p><p>New</p></div></body></html>",
			contains: []string{
				"<p>First</p>",
				"<ins class='diff-insert'><p>New</p></ins>",
			},
		},
		{
			name:   "Script tag handling",
			prev:   "<html><body><script>var old = 1;</script></body></html>",
			latest: "<html><body><script>var new = 2;</script></body></html>",
			contains: []string{
				"<del class='diff-delete'>",
				"var old = 1;",
				"<ins class='diff-insert'>",
				"var new = 2;",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := NewHtmlDiff(tt.prev, tt.latest)
			result, err := diff.Compare()
			if err != nil {
				t.Errorf("Compare() error = %v", err)
				return
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Compare() result does not contain %q", substr)
				}
			}
		})
	}
}

func TestCompareWithURLs(t *testing.T) {
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<p>Old content</p>"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<p>New content</p>"))
	}))
	defer server2.Close()

	diff := NewHtmlDiff(server1.URL, server2.URL)
	result, err := diff.Compare()
	if err != nil {
		t.Errorf("Compare() with URLs error = %v", err)
		return
	}

	expectedSubstrings := []string{
		"<del class='diff-delete'>Old content</del>",
		"<ins class='diff-insert'>New content</ins>",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(result, substr) {
			t.Errorf("Compare() with URLs result does not contain %q", substr)
		}
	}
}

func TestCompareUrl(t *testing.T) {
	// 从URL获取HTML内容
	ttt, err := FetchHTML("https://fresh.qianlht.com/wapi/fast-app.html?pkg=com.ss.fresh&channel_id=mjznh5ht&link_id=mjznh5ht-174-h16&utm_ad_id=__bundle__&landingPageCode=860895&page=pages/spa&pageExt2=shakeRedPage&automatic=1&autocookie=0&mediumType=mj")
	if err != nil {
		t.Fatalf("获取URL内容失败: %v", err)
	}

	// 读取old.html文件内容
	www, err := os.ReadFile("old.html")
	if err != nil {
		t.Fatalf("读取old.html文件失败: %v", err)
	}
	wwwStr := string(www)

	// 创建HTMLDirectDiff实例并比较
	// diff := NewHTMLDirectDiff(wwwStr, ttt)
	diff := NewHtmlDiff(wwwStr, ttt)
	result, err := diff.Compare()
	if err != nil {
		t.Fatalf("比较HTML失败: %v", err)
	}

	// 保存差异比对结果
	err = diff.SaveToFile("res_direct.html")
	if err != nil {
		t.Fatalf("保存HTML失败: %v", err)
	}

	// 验证结果中包含预期的HTML标记和样式
	expectedSubstrings := []string{
		"<div class='diff-container'>",
		"<pre>",
		"<code>",
		"<span class='diff-delete'>",
		"<span class='diff-insert'>",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(result, substr) {
			t.Errorf("结果中缺少预期的HTML标记: %s", substr)
		}
	}
}

func TestSaveToFile(t *testing.T) {
	// 创建测试数据
	prev := "<html><body><p>Old content</p></body></html>"
	latest := "<html><body><p>New content</p></body></html>"
	diff := NewHtmlDiff(prev, latest)

	// 创建临时文件路径
	outputPath := "test_diff_result.html"
	defer os.Remove(outputPath) // 测试完成后清理文件

	// 保存差异比对结果
	err := diff.SaveToFile(outputPath)
	if err != nil {
		t.Errorf("SaveToFile() error = %v", err)
	}

	// 验证文件是否创建成功
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("SaveToFile() 未能成功创建文件")
	}

	// 读取生成的文件内容
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Errorf("读取生成的文件失败: %v", err)
	}

	// 验证文件内容
	expectedSubstrings := []string{
		"<!DOCTYPE html>",
		"<title>HTML Diff Result</title>",
		"<del class='diff-delete'>Old content</del>",
		"<ins class='diff-insert'>New content</ins>",
		".diff-delete { background-color: #ffcdd2; color: #b71c1c;",
		".diff-insert { background-color: #c8e6c9; color: #1b5e20;",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(string(content), substr) {
			t.Errorf("生成的文件缺少预期内容: %s", substr)
		}
	}
}
