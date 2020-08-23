package anitogo

import (
	"strings"
	"testing"
)

func TestErrorTraceError(t *testing.T) {
	err := traceError(indexTooLargeErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), indexTooLargeErr) == -1 {
			t.Errorf("expected %s in error, got %s", indexTooLargeErr, err.Error())
		}
	}
	err = traceError(indexTooSmallErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), indexTooSmallErr) == -1 {
			t.Errorf("expected %s in error, got %s", indexTooSmallErr, err.Error())
		}
	}
	err = traceError(endIndexTooSmallErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), endIndexTooSmallErr) == -1 {
			t.Errorf("expected %s in error, got %s", endIndexTooSmallErr, err.Error())
		}
	}
	err = traceError(tokensEmptyErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), tokensEmptyErr) == -1 {
			t.Errorf("expected %s in error, got %s", tokensEmptyErr, err.Error())
		}
	}
}
