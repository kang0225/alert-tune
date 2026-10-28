// Package rules reads Prometheus alerting rule files.
package rules

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.yaml.in/yaml/v3"
)

// File contains the groups from a Prometheus rule file.
type File struct {
	Groups []Group `yaml:"groups"`
}

// Group contains alerting rules and their shared evaluation settings.
// Duration fields retain their original Prometheus duration notation.
type Group struct {
	Name        string            `yaml:"name"`
	Interval    string            `yaml:"interval,omitempty"`
	Limit       int               `yaml:"limit,omitempty"`
	QueryOffset string            `yaml:"query_offset,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Rules       []Rule            `yaml:"rules"`
}

// Rule contains an alert's expression, timing, and metadata.
// Expressions, duration strings, and templates are preserved without evaluation.
type Rule struct {
	Alert         string            `yaml:"alert"`
	Expr          string            `yaml:"expr"`
	For           string            `yaml:"for,omitempty"`
	KeepFiringFor string            `yaml:"keep_firing_for,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
}

// Parse reads alerting rules from the first YAML document.
// Unknown fields, including recording rules, are rejected. PromQL expressions
// and duration syntax are left for later validation when they are used.
func Parse(reader io.Reader) (*File, error) {
	// Read all input into memory.
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read rules YAML: %w", err)
	}
	// Prepare the YAML decoder.
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	// Reject unknown YAML fields. If a field is missing from the struct, it will be reported as an unknown field.
	decoder.KnownFields(true)

	// Fill the struct using YAML tags.
	var file File
	if err := decoder.Decode(&file); err != nil {
		// If err = io.EOF, the decoder has not read any data, which means the yaml document has no content.
		if errors.Is(err, io.EOF) {
			return nil, errors.New("parse rules YAML: input is empty")
		}
		return nil, fmt.Errorf("parse rules YAML: %w", err)
	}

	// Require at least one group.
	if len(file.Groups) == 0 {
		return nil, errors.New("parse rules YAML: at least one group is required")
	}

	for groupIndex, group := range file.Groups {
		// Trimspace removes any leading or trailing whitespace.
		// If the name is empty or only whitespace, it causes an error.
		if strings.TrimSpace(group.Name) == "" {
			return nil, fmt.Errorf("rule group %d: name is required", groupIndex+1)
		}

		if len(group.Rules) == 0 {
			return nil, fmt.Errorf("rule group %q: at least one rule is required", group.Name)
		}
		// Check each alert's required fields.
		for ruleIndex, rule := range group.Rules {
			if strings.TrimSpace(rule.Alert) == "" {
				return nil, fmt.Errorf("rule group %q, rule %d: alert is required", group.Name, ruleIndex+1)
			}
			if strings.TrimSpace(rule.Expr) == "" {
				return nil, fmt.Errorf("rule group %q, alert %q: expr is required", group.Name, rule.Alert)
			}
		}
	}

	return &file, nil
}
