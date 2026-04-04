package errors

import (
	"net/http"
	"testing"
)

func TestRegisterAndParseCoder(t *testing.T) {
	t.Run("registered code is parseable", func(t *testing.T) {
		Register("TestErr", 999001, http.StatusConflict, "conflict error")
		err := WithCode("TestErr", New("something went wrong"))
		coder := ParseCoder(err)

		if coder.Code() != 999001 {
			t.Errorf("Code() = %d, want 999001", coder.Code())
		}
		if coder.HTTPStatus() != http.StatusConflict {
			t.Errorf("HTTPStatus() = %d, want %d", coder.HTTPStatus(), http.StatusConflict)
		}
		if coder.String() != "conflict error" {
			t.Errorf("String() = %q, want %q", coder.String(), "conflict error")
		}
	})

	t.Run("re-register overwrites", func(t *testing.T) {
		Register("ReRegErr", 999002, http.StatusBadRequest, "original")
		Register("ReRegErr", 999002, http.StatusUnauthorized, "updated")
		err := WithCode("ReRegErr", New("x"))
		coder := ParseCoder(err)
		if coder.HTTPStatus() != http.StatusUnauthorized {
			t.Errorf("HTTPStatus() = %d, want %d", coder.HTTPStatus(), http.StatusUnauthorized)
		}
		if coder.String() != "updated" {
			t.Errorf("String() = %q, want %q", coder.String(), "updated")
		}
	})

	t.Run("unregistered key falls back to 500000", func(t *testing.T) {
		err := WithCode("NoSuchKey", New("oops"))
		coder := ParseCoder(err)
		if coder.Code() != 500000 {
			t.Errorf("Code() = %d, want 500000", coder.Code())
		}
	})

	t.Run("nil error returns nil coder", func(t *testing.T) {
		if ParseCoder(nil) != nil {
			t.Error("ParseCoder(nil) should return nil")
		}
	})

	t.Run("plain error returns unknown coder", func(t *testing.T) {
		coder := ParseCoder(New("plain error"))
		if coder.HTTPStatus() != http.StatusInternalServerError {
			t.Errorf("HTTPStatus() = %d, want 500", coder.HTTPStatus())
		}
	})
}

func TestErrCode_HTTPStatus_ZeroFallback(t *testing.T) {
	Register("ZeroHTTP", 999003, 0, "zero http status")
	err := WithCode("ZeroHTTP", New("x"))
	coder := ParseCoder(err)
	if coder.HTTPStatus() != http.StatusInternalServerError {
		t.Errorf("HTTPStatus() = %d, want 500 for zero httpStatus", coder.HTTPStatus())
	}
}

func TestWithCode(t *testing.T) {
	Register("WCErr", 999010, http.StatusBadRequest, "wc error")

	t.Run("nil errInfo returns nil", func(t *testing.T) {
		if WithCode("WCErr", nil) != nil {
			t.Error("WithCode with nil errInfo should return nil")
		}
	})

	t.Run("wraps plain error", func(t *testing.T) {
		err := WithCode("WCErr", New("original"))
		if err == nil {
			t.Fatal("should not be nil")
		}
		if err.Error() != "original" {
			t.Errorf("Error() = %q, want %q", err.Error(), "original")
		}
	})

	t.Run("wraps string via default", func(t *testing.T) {
		err := WithCode("WCErr", "string message")
		if err == nil {
			t.Fatal("should not be nil")
		}
	})

	t.Run("wraps withCode preserves cause", func(t *testing.T) {
		Register("WCErr2", 999011, http.StatusBadRequest, "wc2")
		inner := WithCode("WCErr", New("cause")).(*withCode)
		outer := WithCode("WCErr2", inner)
		coder := ParseCoder(outer)
		if coder.Code() != 999011 {
			t.Errorf("Code() = %d, want 999011", coder.Code())
		}
		if outer.Error() != "cause" {
			t.Errorf("Error() = %q, want %q", outer.Error(), "cause")
		}
	})
}
