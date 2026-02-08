package aggregator

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// --- fakes for testing ---

type fakeProcessor struct {
	fn  func(MetricsStore)
	err error
}

func (f *fakeProcessor) Process(r io.Reader, store MetricsStore) error {
	if f.err != nil {
		return f.err
	}
	if f.fn != nil {
		f.fn(store)
	}
	return nil
}

type fakeWriter struct {
	called bool
	store  MetricsStore
	err    error
}

func (f *fakeWriter) WriteReports(store MetricsStore) error {
	f.called = true
	f.store = store
	return f.err
}

// --- tests ---

func TestService_Run(t *testing.T) {
	proc := &fakeProcessor{
		fn: func(store MetricsStore) {
			store.Add("camp1", 1000, 0, 0, 0)
		},
	}
	writer := &fakeWriter{}
	svc := NewService(proc, writer, false)

	if err := svc.Run(strings.NewReader("")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !writer.called {
		t.Error("expected writer to be called")
	}
	if writer.store == nil {
		t.Fatal("expected writer to receive store")
	}
}

func TestService_ProcessError(t *testing.T) {
	proc := &fakeProcessor{err: errors.New("parse failed")}
	writer := &fakeWriter{}
	svc := NewService(proc, writer, false)

	err := svc.Run(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error from processor")
	}
	if writer.called {
		t.Error("writer should not be called when processor fails")
	}
}

func TestService_WriterError(t *testing.T) {
	proc := &fakeProcessor{
		fn: func(store MetricsStore) {
			store.Add("camp1", 0, 0, 0, 0)
		},
	}
	writer := &fakeWriter{err: errors.New("disk full")}
	svc := NewService(proc, writer, false)

	err := svc.Run(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error from writer")
	}
}
