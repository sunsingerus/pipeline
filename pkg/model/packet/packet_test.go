// Copyright 2024 Vladislav Klimenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package packet

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPacketFromSlice(t *testing.T) {
	tests := []struct {
		input []int
	}{
		{
			input: []int{},
		},
		{
			input: []int{0},
		},
		{
			input: []int{0, 1},
		},
	}
	for _, tt := range tests {
		pack := New(tt.input)
		require.ElementsMatch(t, tt.input, pack.Slice(), "Check packet from slice: %s", pack)
	}
}

func TestPacketFromSize(t *testing.T) {
	tests := []struct {
		size int
	}{
		{
			size: 0,
		},
		{
			size: 1,
		},
		{
			size: 10,
		},
	}
	for _, tt := range tests {
		pack := New(tt.size)
		require.Equal(t, tt.size, pack.Len(), "Check packet from size: %s", pack)
	}
}
