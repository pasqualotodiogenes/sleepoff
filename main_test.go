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
		{input: "15", want: 15 * time.Minute},
		{input: "30", want: 30 * time.Minute},
		{input: "45m", want: 45 * time.Minute},
		{input: "90s", want: 90 * time.Second},
		{input: "1h", want: 1 * time.Hour},
		{input: "2h", want: 2 * time.Hour},
		{input: "1h30m", want: 90 * time.Minute},
		{input: "2h15m", want: 135 * time.Minute},
		{input: "300s", want: 300 * time.Second},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := parseDuration(tc.input)
			if err != nil {
				t.Fatalf("parseDuration(%q) retornou erro inesperado: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("parseDuration(%q) = %v, esperado %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseDurationInvalid(t *testing.T) {
	t.Parallel()

	tests := []string{
		"0",
		"-1",
		"-15",
		"0s",
		"0m",
		"0h",
		"-5m",
		"-1h",
		"abc",
		"",
		"   ",
		"10x",
		"hours",
	}
	for _, input := range tests {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := parseDuration(input); err == nil {
				t.Fatalf("parseDuration(%q) deveria falhar, mas retornou sucesso", input)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input time.Duration
		want  string
	}{
		{input: 90 * time.Second, want: "1m"}, // 1m30s -> int(hours)>0 && int(minutes)>0 ou else
		{input: 30 * time.Second, want: "30s"},
		{input: 15 * time.Minute, want: "15m"},
		{input: 60 * time.Minute, want: "1h"},
		{input: 90 * time.Minute, want: "1h30m"},
		{input: 150 * time.Minute, want: "2h30m"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			got := formatDuration(tc.input)
			if got != tc.want {
				t.Fatalf("formatDuration(%v) = %q, esperado %q", tc.input, got, tc.want)
			}
		})
	}
}
