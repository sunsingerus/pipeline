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

package accum

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sunsingerus/pipeline/pkg/controller/packet"
	model "github.com/sunsingerus/pipeline/pkg/model/packet"
)

func TestAccum(t *testing.T) {

	tests := []struct {
		input  []int
		expect int
	}{
		{
			input:  []int{},
			expect: 0,
		},
		{
			input:  []int{0},
			expect: 0,
		},
		{
			input:  []int{0, 1},
			expect: 1,
		},
		{
			input:  []int{0, 1},
			expect: 2,
		},
	}
	ch := make(chan packet.Packet)
	accum := New(ch)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go accum.Run(ctx, wg)

	for _, tt := range tests {
		pack := model.New(tt.input)
		ch <- pack
		time.Sleep(time.Second)
		require.Equal(t, tt.expect, accum.Get(), "Check accumulation: %s", pack)
	}

	accum.Get()
	cancel()
	wg.Wait()
	close(ch)
}
