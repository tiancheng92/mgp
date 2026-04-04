package mgp

import (
	"testing"
)

func TestGetFullPath(t *testing.T) {
	tests := []struct {
		groupName    string
		relativePath string
		want         string
	}{
		{"", "/users", "/users"},
		{"/api", "/users", "/api/users"},
		{"/api/v1", "/users", "/api/v1/users"},
		{"/api", "/users/", "/api/users/"},
		{"/api/v1", "", "/api/v1"},
	}
	for _, tt := range tests {
		got := getFullPath(tt.groupName, tt.relativePath)
		if got != tt.want {
			t.Errorf("getFullPath(%q, %q) = %q, want %q", tt.groupName, tt.relativePath, got, tt.want)
		}
	}
}

func TestCamelString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "Hello"},
		{"hello_world", "HelloWorld"},
		{"api_v1", "ApiV1"},
		{"already_Camel", "Already_Camel"}, // _ followed by uppercase is not transformed
		{"", ""},
		{"v1", "V1"},
		{"my_api_key", "MyApiKey"},
	}
	for _, tt := range tests {
		got := camelString(tt.input)
		if got != tt.want {
			t.Errorf("camelString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestStringToBytes(t *testing.T) {
	s := "hello world"
	b := StringToBytes(s)
	if string(b) != s {
		t.Errorf("StringToBytes(%q) roundtrip failed, got %q", s, string(b))
	}
	if len(b) != len(s) {
		t.Errorf("StringToBytes(%q) len = %d, want %d", s, len(b), len(s))
	}

	empty := StringToBytes("")
	if len(empty) != 0 {
		t.Errorf("StringToBytes(\"\") should be empty")
	}
}

func TestBytesToString(t *testing.T) {
	b := []byte("hello world")
	s := BytesToString(b)
	if s != string(b) {
		t.Errorf("BytesToString roundtrip failed, got %q", s)
	}

	empty := BytesToString([]byte{})
	if empty != "" {
		t.Errorf("BytesToString([]byte{}) should be empty string")
	}
}

func TestGetFuncName(t *testing.T) {
	// 空参数
	if got := getFuncName(); got != "" {
		t.Errorf("getFuncName() with no args = %q, want \"\"", got)
	}

	// nil 函数
	if got := getFuncName(nil); got != "" {
		t.Errorf("getFuncName(nil) = %q, want \"\"", got)
	}

	// 正常函数
	f := func(c *Context) {}
	got := getFuncName(f)
	if got == "" {
		t.Errorf("getFuncName(func) should return non-empty name")
	}
}
