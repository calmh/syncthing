package model

import (
	"bytes"
	"testing"
)

func TestMergeLists(t *testing.T) {
	t.Parallel()

	cases := []struct{ src, dst, res string }{
		{"abc", "", "abc"},
		{"", "abc", "abc"},
		{"a", "a", "a"},
		{"a", "b", "ab"},
		{"a", "ab", "ab"},
		{"ab", "a", "ab"},
		{"ab", "b", "ab"},
		{"ab", "ab", "ab"},
		{"ace", "bdf", "abcdef"},
		{"accceezz", "bdf", "abcdefz"},
	}

	for _, tc := range cases {
		src := bytes.Split([]byte(tc.src), nil)
		dst := bytes.Split([]byte(tc.dst), nil)
		res := mergeLists(dst, src)
		resStr := string(bytes.Join(res, nil))
		if resStr != tc.res {
			t.Errorf("join %q into %q => %q, expected %q", tc.src, tc.dst, resStr, tc.res)
		}
	}
}
