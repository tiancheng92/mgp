package mgp

import (
	"strings"
	"testing"
)

type swagTestStruct struct{ Name string }

func TestGetStructAndPackageName(t *testing.T) {
	t.Run("non-pointer struct", func(t *testing.T) {
		name := getStructAndPackageName(swagTestStruct{})
		if !strings.Contains(name, "swagTestStruct") {
			t.Errorf("got %q, want name containing swagTestStruct", name)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		name := getStructAndPackageName(&swagTestStruct{})
		if !strings.Contains(name, "swagTestStruct") {
			t.Errorf("got %q, want name containing swagTestStruct", name)
		}
	})

	t.Run("typed nil pointer does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("getStructAndPackageName panicked: %v", r)
			}
		}()
		var p *swagTestStruct
		name := getStructAndPackageName(p)
		if !strings.Contains(name, "swagTestStruct") {
			t.Errorf("got %q, want name containing swagTestStruct", name)
		}
	})

	t.Run("pointer and non-pointer yield same name", func(t *testing.T) {
		a := getStructAndPackageName(swagTestStruct{})
		b := getStructAndPackageName(&swagTestStruct{})
		if a != b {
			t.Errorf("pointer=%q, non-pointer=%q, should be equal", b, a)
		}
	})
}

func TestAddTextIfNotEmptyOrDefault(t *testing.T) {
	format := "// @Accept %s\n"

	t.Run("nil text uses default", func(t *testing.T) {
		var s strings.Builder
		addTextIfNotEmptyOrDefault(&s, "json", format)
		if s.String() != "// @Accept json\n" {
			t.Errorf("got %q", s.String())
		}
	})

	t.Run("empty text uses default", func(t *testing.T) {
		var s strings.Builder
		addTextIfNotEmptyOrDefault(&s, "json", format, "")
		if s.String() != "// @Accept json\n" {
			t.Errorf("got %q", s.String())
		}
	})

	t.Run("provided text overrides default", func(t *testing.T) {
		var s strings.Builder
		addTextIfNotEmptyOrDefault(&s, "json", format, "multipart/form-data")
		if s.String() != "// @Accept multipart/form-data\n" {
			t.Errorf("got %q", s.String())
		}
	})

	t.Run("multiple values joined by comma", func(t *testing.T) {
		var s strings.Builder
		addTextIfNotEmptyOrDefault(&s, "json", format, "json", "xml")
		if s.String() != "// @Accept json,xml\n" {
			t.Errorf("got %q", s.String())
		}
	})

	t.Run("empty default with no text writes nothing", func(t *testing.T) {
		var s strings.Builder
		addTextIfNotEmptyOrDefault(&s, "", format)
		if s.String() != "" {
			t.Errorf("expected empty output, got %q", s.String())
		}
	})
}

func TestAddLineIfNotEmpty(t *testing.T) {
	t.Run("empty data writes nothing", func(t *testing.T) {
		var s strings.Builder
		addLineIfNotEmpty(&s, "", "// @Summary %s\n")
		if s.String() != "" {
			t.Errorf("expected empty, got %q", s.String())
		}
	})

	t.Run("non-empty data writes line", func(t *testing.T) {
		var s strings.Builder
		addLineIfNotEmpty(&s, "List users", "// @Summary %s\n")
		if s.String() != "// @Summary List users\n" {
			t.Errorf("got %q", s.String())
		}
	})
}

func TestWriteRoutes_ProduceAlwaysEmitted(t *testing.T) {
	// 验证没有显式设置 Returns 的路由也能生成 @Produce 注解 (Bug 3 修复验证)
	routes := []*Route{
		{
			Path:     "/test",
			Method:   "GET",
			FuncName: "TestHandler",
			Summary:  "Test",
			// Returns 为 nil，由 writeRoutes 赋默认值
		},
	}
	var s strings.Builder
	pkgs := make(map[string]bool)
	writeRoutes("", routes, &s, pkgs)

	output := s.String()
	if !strings.Contains(output, "@Produce") {
		t.Errorf("expected @Produce in output, got:\n%s", output)
	}
}

func TestWriteRoutes_HiddenSkipped(t *testing.T) {
	routes := []*Route{
		{Path: "/visible", Method: "GET", FuncName: "Visible"},
		{Path: "/hidden", Method: "GET", FuncName: "Hidden", Hidden: true},
	}
	var s strings.Builder
	writeRoutes("", routes, &s, make(map[string]bool))

	output := s.String()
	if strings.Contains(output, "/hidden") {
		t.Error("hidden route should not appear in output")
	}
	if !strings.Contains(output, "/visible") {
		t.Error("visible route should appear in output")
	}
}

func TestWriteRoutes_PathParamConverted(t *testing.T) {
	routes := []*Route{
		{Path: "/users/:id/posts/:postId", Method: "GET", FuncName: "GetPost"},
	}
	var s strings.Builder
	writeRoutes("", routes, &s, make(map[string]bool))

	output := s.String()
	if !strings.Contains(output, "/users/{id}/posts/{postId}") {
		t.Errorf("path params not converted, output:\n%s", output)
	}
}

func TestWriteRoutes_AcceptOnlyForMutatingMethods(t *testing.T) {
	methods := []struct {
		method      string
		wantAccept  bool
	}{
		{"POST", true},
		{"PUT", true},
		{"PATCH", true},
		{"GET", false},
		{"DELETE", false},
	}

	for _, m := range methods {
		var s strings.Builder
		routes := []*Route{{Path: "/x", Method: m.method, FuncName: "F"}}
		writeRoutes("", routes, &s, make(map[string]bool))
		hasAccept := strings.Contains(s.String(), "@Accept")
		if hasAccept != m.wantAccept {
			t.Errorf("method %s: @Accept present=%v, want %v", m.method, hasAccept, m.wantAccept)
		}
	}
}
