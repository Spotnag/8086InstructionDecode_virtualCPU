package main

import (
	"strings"
	"testing"
)

func BenchmarkStringConcatenation(b *testing.B) {
	// Using the + operator
	result := ""
	for i := 0; i < b.N; i++ {
		result += "foo"
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	// Using strings.Builder
	builder := strings.Builder{}
	for i := 0; i < b.N; i++ {
		builder.WriteString("foo")
	}
	_ = builder.String() // Ignoring the final string to match the previous benchmark
}
