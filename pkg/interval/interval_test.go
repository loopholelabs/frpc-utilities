// SPDX-License-Identifier: Apache-2.0

package interval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	i := New()

	require.Equal(t, false, i.Contains(1, 2))
	i.Insert(1, 2)
	require.True(t, i.Contains(1, 2))

	require.False(t, i.Contains(0, 1))
	require.False(t, i.Contains(1, 3))

	i.Insert(2, 5)
	require.False(t, i.Contains(0, 1))
	require.True(t, i.Contains(1, 3))
	require.True(t, i.Contains(3, 5))

	i.Insert(1, 5)
	require.False(t, i.Contains(0, 1))

	i.Insert(0, 1)
	require.True(t, i.Contains(0, 5))
	require.False(t, i.Contains(0, 6))

	i.Insert(3, 4)
	require.True(t, i.Contains(0, 5))
	require.False(t, i.Contains(0, 6))

	i.Insert(7, 10)
	require.True(t, i.Contains(0, 5))
	require.False(t, i.Contains(0, 6))
	require.True(t, i.Contains(7, 10))
	require.False(t, i.Contains(7, 11))
	require.False(t, i.Contains(6, 7))

	i.Insert(6, 7)
	require.True(t, i.Contains(0, 5))
	require.True(t, i.Contains(0, 6))
	require.True(t, i.Contains(0, 10))
}
