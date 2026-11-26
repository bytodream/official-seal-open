package dto

import "testing"

// TestIsBlank 空白字符与 utils.IsBlank 语义一致（空白串视为空）。
func TestIsBlank(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"a", false},
		{" CIRCULAR_42 ", false},
	}
	for _, tc := range cases {
		if got := IsBlank(tc.in); got != tc.want {
			t.Fatalf("IsBlank(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
