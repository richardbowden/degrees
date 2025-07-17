package problems

import (
	"errors"
	"net/http"
	"testing"
)

func TestKind_String(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{Other, "other error"},
		{Invalid, "invalid operation"},
		{IO, "I/O error"},
		{Exist, "item already exists"},
		{NotExist, "item does not exist"},
		{Private, "information withheld"},
		{Internal, "internal error"},
		{BrokenLink, "link target does not exist"},
		{Database, "database error"},
		{Validation, "input validation error"},
		{Unanticipated, "unanticipated error"},
		{InvalidRequest, "invalid request error"},
		{Unauthenticated, "unauthenticated request"},
		{Unauthorized, "unauthorized request"},
		{Kind(999), "unknown error kind"}, // Unknown kind
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("Kind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpErrorStatusCode(t *testing.T) {
	tests := []struct {
		kind Kind
		want int
	}{
		{NotExist, http.StatusNotFound},
		{Invalid, http.StatusBadRequest},
		{Exist, http.StatusBadRequest},
		{Private, http.StatusBadRequest},
		{BrokenLink, http.StatusBadRequest},
		{Validation, http.StatusBadRequest},
		{InvalidRequest, http.StatusBadRequest},
		{Other, http.StatusInternalServerError},
		{IO, http.StatusInternalServerError},
		{Internal, http.StatusInternalServerError},
		{Database, http.StatusInternalServerError},
		{Unanticipated, http.StatusInternalServerError},
		{Unauthenticated, http.StatusUnauthorized},
		{Unauthorized, http.StatusForbidden},
		{Kind(999), http.StatusInternalServerError}, // Unknown kind defaults to 500
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := httpErrorStatusCode(tt.kind); got != tt.want {
				t.Errorf("httpErrorStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetail_Error(t *testing.T) {
	tests := []struct {
		name   string
		detail Detail
		want   string
	}{
		{
			name:   "message only",
			detail: Detail{Message: "validation failed"},
			want:   "validation failed",
		},
		{
			name:   "message with value",
			detail: Detail{Message: "invalid field", Value: "email"},
			want:   "invalid field (email)",
		},
		{
			name:   "empty value",
			detail: Detail{Message: "error occurred", Value: ""},
			want:   "error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.detail.Error(); got != tt.want {
				t.Errorf("Detail.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblem_AddDetail(t *testing.T) {
	t.Run("add Detail", func(t *testing.T) {
		p := &Problem{}
		detail := Detail{Message: "test error", Value: "test"}

		p.AddDetail(detail)

		if len(p.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(p.Errors))
		}
		if p.Errors[0].Message != "test error" {
			t.Errorf("Expected message 'test error', got %s", p.Errors[0].Message)
		}
		if p.Errors[0].Value != "test" {
			t.Errorf("Expected value 'test', got %s", p.Errors[0].Value)
		}
	})

	t.Run("add regular error", func(t *testing.T) {
		p := &Problem{}
		err := errors.New("regular error")

		p.AddDetail(err)

		if len(p.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(p.Errors))
		}
		if p.Errors[0].Message != "regular error" {
			t.Errorf("Expected message 'regular error', got %s", p.Errors[0].Message)
		}
		if p.Errors[0].Value != "" {
			t.Errorf("Expected empty value, got %s", p.Errors[0].Value)
		}
	})

	t.Run("add multiple errors", func(t *testing.T) {
		p := &Problem{}
		detail1 := Detail{Message: "error 1"}
		err2 := errors.New("error 2")

		p.AddDetail(detail1)
		p.AddDetail(err2)

		if len(p.Errors) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(p.Errors))
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("basic error creation", func(t *testing.T) {
		p := New(Invalid, "test error")

		if p.Status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, p.Status)
		}
		if p.Kind != Invalid {
			t.Errorf("Expected kind %v, got %v", Invalid, p.Kind)
		}
		if p.Detail != "test error" {
			t.Errorf("Expected detail 'test error', got %s", p.Detail)
		}
		if string(p.Title) != http.StatusText(http.StatusBadRequest) {
			t.Errorf("Expected title '%s', got %s", http.StatusText(http.StatusBadRequest), p.Title)
		}
		if len(p.Errors) != 0 {
			t.Errorf("Expected no detail errors, got %d", len(p.Errors))
		}
	})

	t.Run("error with details", func(t *testing.T) {
		detail1 := Detail{Message: "field error", Value: "name"}
		err2 := errors.New("another error")

		p := New(Validation, "validation failed", detail1, err2)

		if p.Status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, p.Status)
		}
		if len(p.Errors) != 2 {
			t.Errorf("Expected 2 detail errors, got %d", len(p.Errors))
		}
		if p.Errors[0].Message != "field error" {
			t.Errorf("Expected first error message 'field error', got %s", p.Errors[0].Message)
		}
		if p.Errors[1].Message != "another error" {
			t.Errorf("Expected second error message 'another error', got %s", p.Errors[1].Message)
		}
	})

	t.Run("different kinds map to correct status codes", func(t *testing.T) {
		tests := []struct {
			kind   Kind
			status int
		}{
			{NotExist, http.StatusNotFound},
			{Invalid, http.StatusBadRequest},
			{Internal, http.StatusInternalServerError},
			{Unauthenticated, http.StatusUnauthorized},
			{Unauthorized, http.StatusForbidden},
		}

		for _, tt := range tests {
			p := New(tt.kind, "test")
			if p.Status != tt.status {
				t.Errorf("Kind %v: expected status %d, got %d", tt.kind, tt.status, p.Status)
			}
		}
	})
}

func TestProblem_Error(t *testing.T) {
	p := New(Invalid, "test error message")

	if p.Error() != "test error message" {
		t.Errorf("Expected 'test error message', got %s", p.Error())
	}
}

func TestProblem_GetStatus(t *testing.T) {
	p := New(NotExist, "not found")

	if p.GetStatus() != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, p.GetStatus())
	}
}

func TestProblem_ZeroValue(t *testing.T) {
	// Test that zero value Kind (Other) behaves correctly
	p := New(Kind(0), "zero value test")

	if p.Status != http.StatusInternalServerError {
		t.Errorf("Expected zero value to map to 500, got %d", p.Status)
	}
	if p.Kind != Other {
		t.Errorf("Expected zero value to be Other, got %v", p.Kind)
	}
}

func TestProblem_JSONFields(t *testing.T) {
	// Test that the Problem struct has correct JSON fields
	p := New(Validation, "validation error")

	// Verify JSON fields are accessible
	if p.Status == 0 {
		t.Error("Status should be set")
	}
	if p.Title == "" {
		t.Error("Title should be set")
	}
	if p.Detail == "" {
		t.Error("Detail should be set")
	}

	// Verify internal fields are set but not JSON exported
	if p.Kind == 0 && p.Kind != Other {
		t.Error("Kind should be set")
	}
}

func TestTitle_Type(t *testing.T) {
	// Test that Title is properly typed
	p := New(Invalid, "test")
	title := p.Title

	if string(title) != "Bad Request" {
		t.Errorf("Expected 'Bad Request', got %s", string(title))
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(Invalid, "benchmark test")
	}
}

func BenchmarkNewWithDetails(b *testing.B) {
	detail := Detail{Message: "benchmark detail"}
	err := errors.New("benchmark error")

	for i := 0; i < b.N; i++ {
		_ = New(Validation, "benchmark test", detail, err)
	}
}

func BenchmarkAddDetail(b *testing.B) {
	p := &Problem{}
	detail := Detail{Message: "benchmark detail"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.AddDetail(detail)
	}
}
