package mgp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tiancheng92/mgp/errors"
	"github.com/tiancheng92/mgp/errors/default_error_code"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestResponse_Error(t *testing.T) {
	c, w := newTestContext()
	err := errors.WithCode(default_error_code.ErrClientParam, "bad param")
	Response(c, nil, err)

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var result Result[any]
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Code != 400000 {
		t.Errorf("Code = %d, want 400000", result.Code)
	}
}

func TestResponse_Success(t *testing.T) {
	c, w := newTestContext()
	Response(c, "hello", nil)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var result Result[any]
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Code != successCode {
		t.Errorf("Code = %d, want %d", result.Code, successCode)
	}
	if result.Msg != successMsg {
		t.Errorf("Msg = %q, want %q", result.Msg, successMsg)
	}
}

func TestResponse_NilData(t *testing.T) {
	c, w := newTestContext()
	Response(c, nil, nil)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}
	var result Result[any]
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if result.Data != nil {
		t.Errorf("Data should be nil, got %v", result.Data)
	}
}

func TestResponse_PaginateData(t *testing.T) {
	n := 1
	pd := &PaginateData[int]{
		PaginateInfo: &PaginateInfo{Total: 1, Page: 1, PageSize: 20},
		Items:        []*int{&n},
	}
	c, w := newTestContext()
	Response(c, pd, nil)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}

	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatalf("data field is not an object, got %T", result["data"])
	}
	if _, ok := data["paginate"]; !ok {
		t.Error("paginate field missing from response data")
	}
	if _, ok := data["items"]; !ok {
		t.Error("items field missing from response data")
	}
}

func TestResponse_Aborted(t *testing.T) {
	c, w := newTestContext()
	c.Abort()
	Response(c, "should not render", nil)

	if w.Body.Len() != 0 {
		t.Errorf("aborted context should produce no body, got: %s", w.Body.String())
	}
}

func TestResponse_ErrorSetsLogLevel(t *testing.T) {
	t.Run("4xx sets warn level", func(t *testing.T) {
		c, _ := newTestContext()
		err := errors.WithCode(default_error_code.ErrClientParam, "bad")
		Response(c, nil, err)
		if _, exists := c.Get(ErrorLogLevelWarn); !exists {
			t.Error("warn log level key should be set for 4xx error")
		}
	})

	t.Run("5xx sets error level", func(t *testing.T) {
		c, _ := newTestContext()
		err := errors.WithCode(default_error_code.ErrServer, "internal")
		Response(c, nil, err)
		if _, exists := c.Get(ErrorLogLevelError); !exists {
			t.Error("error log level key should be set for 5xx error")
		}
	})
}

func TestResponseDownloadStream(t *testing.T) {
	t.Run("normal download", func(t *testing.T) {
		c, w := newTestContext()
		buf := bytes.NewBufferString("file content")
		ResponseDownloadStream(c, buf, "test.txt")

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status = %d, want 200", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/octet-stream" {
			t.Errorf("Content-Type = %q, want application/octet-stream", ct)
		}
		if cd := w.Header().Get("Content-Disposition"); cd != `attachment; filename="test.txt"` {
			t.Errorf("Content-Disposition = %q", cd)
		}
		if w.Body.String() != "file content" {
			t.Errorf("body = %q, want %q", w.Body.String(), "file content")
		}
	})

	t.Run("aborted context skips write", func(t *testing.T) {
		c, w := newTestContext()
		c.Abort()
		ResponseDownloadStream(c, bytes.NewBufferString("data"), "file.bin")
		if w.Body.Len() != 0 {
			t.Error("aborted context should produce no body")
		}
	})
}
