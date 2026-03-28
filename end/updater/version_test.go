package updater

import "testing"

func TestCompareVersions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		left     string
		right    string
		expected int
	}{
		{name: "equal", left: "v1.2.3", right: "1.2.3", expected: 0},
		{name: "greater", left: "v1.3.0", right: "1.2.9", expected: 1},
		{name: "less", left: "1.2.3", right: "v2.0.0", expected: -1},
		{name: "runtime naming", left: "runtime-v2026.03.2", right: "runtime-v2026.03.1", expected: 1},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if actual := CompareVersions(tc.left, tc.right); actual != tc.expected {
				t.Fatalf("CompareVersions(%q, %q)=%d, want %d", tc.left, tc.right, actual, tc.expected)
			}
		})
	}
}
