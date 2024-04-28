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

package generator

import (
	"context"
	"github.com/stretchr/testify/require"
	packetbuilder "github.com/sunsingerus/pipeline/pkg/controller/generator/packet_builder"
	"sync"
	"testing"

	"github.com/sunsingerus/pipeline/pkg/controller/packet"
	model "github.com/sunsingerus/pipeline/pkg/model/packet"
)

func TestGenerator(t *testing.T) {

	tests := []struct {
		size   int
		expect int
	}{
		{
			size:   10,
			expect: 10,
		},
		{
			size:   20,
			expect: 20,
		},
	}
	for _, tt := range tests {
		ch := make(chan packet.Packet)
		builder := packetbuilder.New(func(size int) packetbuilder.Packet { return model.New(size) }, packetbuilder.Options{Size: tt.size})
		gen := New(ch, builder, Options{Interval: 100})

		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		wg.Add(1)
		go gen.Run(ctx, wg)

		pack := <-ch
		require.Equal(t, tt.expect, pack.Len(), "Check generator: %s", pack)

		cancel()
		wg.Wait()
		close(ch)
		ch = nil
	}
}
