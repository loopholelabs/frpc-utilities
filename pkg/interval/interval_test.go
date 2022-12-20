/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package interval

import (
	"github.com/stretchr/testify/require"
	"testing"
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
