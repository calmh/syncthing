// Copyright (C) 2023 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package sliceutil

import (
	"math"

	"golang.org/x/exp/constraints"
)

// RemoveAndZero removes the element at index i from slice s and returns the
// resulting slice. The slice ordering is preserved; the last slice element
// is zeroed before shrinking.
func RemoveAndZero[E any, S ~[]E](s S, i int) S {
	copy(s[i:], s[i+1:])
	s[len(s)-1] = *new(E)
	return s[:len(s)-1]
}

type Numerical interface {
	constraints.Integer | constraints.Float
}

func Average[E Numerical, S ~[]E](s S) E {
	var total E = 0
	if len(s) == 0 {
		return 0
	}

	for _, value := range s {
		if s, ok := any(value).(float64); ok {
			if math.IsNaN(s) {
				continue
			}
		}
		total += value
	}

	return total / E(len(s))
}
