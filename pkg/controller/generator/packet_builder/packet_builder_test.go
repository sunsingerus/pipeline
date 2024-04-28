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

package packetbuilder

import (
	"github.com/stretchr/testify/require"
	"testing"

	model "github.com/sunsingerus/pipeline/pkg/model/packet"
)

func TestPacketBuilder(t *testing.T) {

	tests := []struct {
		constructor packetConstructor
		options     Options
		expect      int
	}{
		{
			constructor: func(size int) Packet { return model.New(size) },
			options:     Options{Size: 30},
			expect:      30,
		},
	}
	for _, tt := range tests {
		builder := New(tt.constructor, tt.options)
		// Check two packets are of expected size but not equal
		a := builder.Build()
		b := builder.Build()
		require.Equal(t, a.Len(), tt.expect, "Check packet builder: %s", tt.expect)
		require.Equal(t, b.Len(), tt.expect, "Check packet builder: %s", tt.expect)
		require.NotEqual(t, a.String(), b.String(), "Check packet builder: %s", tt.expect)
	}
}
