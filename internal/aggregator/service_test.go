package aggregator

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// --- fakes for testing ---

type fakeProcessor struct {
	result map[string]*CampaignMetrics
	err    error
}

func (f *fakeProcessor) Process(r io.Reader) (map[string]*CampaignMetrics, error) {
	return f.result, f.err
}

type fakeWriter struct {
	called   bool
	received map[string]*CampaignMetrics
	err      error
}

func (f *fakeWriter) WriteReports(metrics map[string]*CampaignMetrics) error {
	f.called = true
	f.received = metrics
	return f.err
}

// --- tests ---

func TestService_Run(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"camp1": {CampaignID: "camp1", TotalImpressions: 1000},
	}

	proc := &fakeProcessor{result: metrics}
	writer := &fakeWriter{}
	svc := NewService(proc, writer)

	got, err := svc.Run(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !writer.called {
		t.Error("expected writer to be called")
	}
	if got["camp1"] == nil {
		t.Error("expected metrics to be returned")
	}
	if writer.received["camp1"] == nil {
		t.Error("expected writer to receive metrics")
	}
}

func TestService_ProcessError(t *testing.T) {
	proc := &fakeProcessor{err: errors.New("parse failed")}
	writer := &fakeWriter{}
	svc := NewService(proc, writer)

	_, err := svc.Run(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error from processor")
	}
	if writer.called {
		t.Error("writer should not be called when processor fails")
	}
}

func TestService_WriterError(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"camp1": {CampaignID: "camp1"},
	}

	proc := &fakeProcessor{result: metrics}
	writer := &fakeWriter{err: errors.New("disk full")}
	svc := NewService(proc, writer)

	_, err := svc.Run(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error from writer")
	}
}
