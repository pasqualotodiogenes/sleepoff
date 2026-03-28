package main

import (
	"testing"
	"time"
)

func TestParseDurationValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  time.Duration
	}{
		{input: "30", want: 30 * time.Minute},
		{input: "90s", want: 90 * time.Second},
		{input: "1h30m", want: 90 * time.Minute},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := parseDuration(tc.input)
			if err != nil {
				t.Fatalf("parseDuration(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("parseDuration(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseDurationInvalid(t *testing.T) {
	t.Parallel()

	tests := []string{"0", "-1", "0s", "-5m", "abc"}
	for _, input := range tests {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := parseDuration(input); err == nil {
				t.Fatalf("parseDuration(%q) expected error", input)
			}
		})
	}
}
