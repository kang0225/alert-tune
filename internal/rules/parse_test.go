package rules

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// TestParse checks that a valid rule file fills every field with the expected values.
func TestParse(t *testing.T) {
	input := `groups:
  - name: orders
    interval: 30s
    limit: 10
    query_offset: 1m
    labels:
      team: checkout
    rules:
      - alert: HighErrors
        expr: |-
          sum(rate(errors_total[5m]))
            / sum(rate(requests_total[5m])) > 0.05
        for: 1d
        keep_firing_for: 2w
        labels:
          severity: page
        annotations:
          summary: 'Errors on {{ $labels.instance }}: {{ $value }}'
      - alert: HighErrors
        expr: up == 0
        labels:
          severity: warning
  - name: latency
    rules:
      - alert: HighLatency
        expr: latency_seconds > 0.5
`
	want := &File{Groups: []Group{
		{
			Name:        "orders",
			Interval:    "30s",
			Limit:       10,
			QueryOffset: "1m",
			Labels:      map[string]string{"team": "checkout"},
			Rules: []Rule{
				{
					Alert:         "HighErrors",
					Expr:          "sum(rate(errors_total[5m]))\n  / sum(rate(requests_total[5m])) > 0.05",
					For:           "1d",
					KeepFiringFor: "2w",
					Labels:        map[string]string{"severity": "page"},
					Annotations:   map[string]string{"summary": "Errors on {{ $labels.instance }}: {{ $value }}"},
				},
				{
					Alert:  "HighErrors",
					Expr:   "up == 0",
					Labels: map[string]string{"severity": "warning"},
				},
			},
		},
		{
			Name: "latency",
			Rules: []Rule{
				{Alert: "HighLatency", Expr: "latency_seconds > 0.5"},
			},
		},
	}}

	got, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

// TestParseRejectsInvalidInput checks that each invalid rule file returns an error.
func TestParseRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		context string
	}{
		{name: "empty input", input: ""},
		{name: "comments only", input: "# no rules\n"},
		{name: "missing groups", input: "{}"},
		{name: "empty groups", input: "groups: []"},
		{name: "malformed YAML", input: "groups: ["},
		{name: "groups must be a sequence", input: "groups: {name: orders}"},
		{name: "missing group name", input: "groups: [{rules: [{alert: HighErrors, expr: up}]}]", context: "name"},
		{name: "blank group name", input: "groups: [{name: '  ', rules: [{alert: HighErrors, expr: up}]}]", context: "name"},
		{name: "missing rules", input: "groups: [{name: orders}]", context: "orders"},
		{name: "empty rules", input: "groups: [{name: orders, rules: []}]", context: "orders"},
		{name: "missing alert", input: "groups: [{name: orders, rules: [{expr: up}]}]", context: "orders"},
		{name: "blank alert", input: "groups: [{name: orders, rules: [{alert: '  ', expr: up}]}]"},
		{name: "missing expression", input: "groups: [{name: orders, rules: [{alert: HighErrors}]}]", context: "HighErrors"},
		{name: "blank expression", input: "groups: [{name: orders, rules: [{alert: HighErrors, expr: '  '}]}]"},
		{name: "expression must be scalar", input: "groups: [{name: orders, rules: [{alert: HighErrors, expr: [up]}]}]"},
		{name: "limit must be an integer", input: "groups: [{name: orders, limit: many, rules: [{alert: HighErrors, expr: up}]}]"},
		{name: "unknown top-level field", input: "other: true", context: "other"},
		{name: "unknown rule field", input: "groups: [{name: orders, rules: [{alert: HighErrors, expr: up, fro: 2m}]}]", context: "fro"},
		{name: "recording rule", input: "groups: [{name: orders, rules: [{record: instance:up, expr: up}]}]", context: "record"},
		{name: "duplicate YAML key", input: "groups: [{name: orders, rules: [{alert: HighErrors, expr: up, expr: down}]}]", context: "expr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(strings.NewReader(tt.input))
			if err == nil {
				t.Fatal("Parse() succeeded, want an error")
			}
			if tt.context != "" && !strings.Contains(err.Error(), tt.context) {
				t.Errorf("Parse() error = %q, want context %q", err, tt.context)
			}
		})
	}
}

// TestParseReaderError checks that a read failure is returned as a wrapped error.
func TestParseReaderError(t *testing.T) {
	want := errors.New("read failed")
	_, err := Parse(errorReader{err: want})
	if !errors.Is(err, want) {
		t.Errorf("Parse() error = %v, want wrapped reader error %v", err, want)
	}
}

// errorReader is a test reader whose Read always fails with err.
type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
