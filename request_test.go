package main

import (
	"reflect"
	"testing"
)

func TestParseRequestTarge(t *testing.T) {
	tests := []struct {
		input string
		path  string
		args  map[string][]string
	}{
		{"test?a=1", "test", map[string][]string{"a": {"1"}}},
		{"test?x=hello=world", "test", map[string][]string{"x": {"hello=world"}}},
		{"test?encoded=base64==", "test", map[string][]string{"encoded": {"base64=="}}},
		{"test?key=value1&key=value2", "test", map[string][]string{"key": {"value1", "value2"}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			resultPath, resultArgs := ParseRequestTarget(tt.input)

			if resultPath != tt.path {
				t.Errorf("Path = %v; want %v", resultPath, tt.path)
			}

			if !reflect.DeepEqual(resultArgs, tt.args) {
				t.Errorf("Args = %v; want %v", resultArgs, tt.args)
			}
		})
	}
}
