package coname

import (
	"testing"
)

func TestNormalizeWords(t *testing.T) {
	for _, tt := range tests {
		got := NormalizeWords(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeWords(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func BenchmarkNormalizeWords(b *testing.B) {
	for b.Loop() {
		NormalizeWords("Gunderson Dettmer Stough Villeneuve Franklin & Hachigian, LLP")
	}
}
