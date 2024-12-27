package utils_test

import (
	"testing"

	"github.com/tdevsin/keyforge/internal/utils"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Should be true if string contains nothing",
			s:    "",
			want: true,
		},
		{
			name: "Should be true if string contains many empty spaces",
			s:    "    ",
			want: true,
		},
		{
			name: "Should be true if string contains only 1 empty space",
			s:    " ",
			want: true,
		},
		{
			name: "Should be false if string contains some value",
			s:    " v ",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsEmpty(tt.s)
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
