// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"errors"
)

var (
	Closed     = errors.New("queue is closed")
	FullError  = errors.New("queue is full")
	EmptyError = errors.New("queue is empty")
)

// round takes an uint64 value and rounds up to the nearest power of 2
func round(value uint64) uint64 {
	value--
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	value++
	return value
}
