package catalog

import "testing"

func TestOffset(t *testing.T) {
	cases := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{"first page", 1, 20, 0},
		{"second page", 2, 20, 20},
		{"page zero clamps to first", 0, 20, 0},
		{"negative page clamps to first", -5, 20, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Offset(c.page, c.pageSize); got != c.want {
				t.Errorf("Offset(%d, %d) = %d, want %d", c.page, c.pageSize, got, c.want)
			}
		})
	}
}

func TestTotalPages(t *testing.T) {
	cases := []struct {
		name       string
		totalCount int
		pageSize   int
		want       int
	}{
		{"zero total", 0, 20, 0},
		{"exact multiple", 40, 20, 2},
		{"non-divisible rounds up", 41, 20, 3},
		{"less than one page", 5, 20, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := TotalPages(c.totalCount, c.pageSize); got != c.want {
				t.Errorf("TotalPages(%d, %d) = %d, want %d", c.totalCount, c.pageSize, got, c.want)
			}
		})
	}
}
