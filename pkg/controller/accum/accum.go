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
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/sunsingerus/pipeline/pkg/controller/packet"
)

type inPacket interface {
	fmt.Stringer
	Len() int
	Get(int) int
}

// Accum specifies accumulator
type Accum struct {
	// in specifies chan where accum reads packets
	in    chan packet.Packet
	accum int
	mux   sync.RWMutex
}

// New creates new accumulator
func New(in chan packet.Packet) *Accum {
	return &Accum{
		in: in,
	}
}

func (a *Accum) Get() int {
	if a == nil {
		return 0
	}
	a.mux.RLock()
	defer a.mux.RUnlock()

	return a.accum
}

func (a *Accum) accumulate(i int) {
	if a == nil {
		return
	}
	a.mux.Lock()
	defer a.mux.Unlock()

	a.accum += i
}

func (a *Accum) processPacket(in inPacket) {
	if a == nil {
		return
	}

	// Accumulate all values from the packet
	for i := 0; i < in.Len(); i++ {
		a.accumulate(in.Get(i))
	}
}

// Run runs accum until context is done
func (a *Accum) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if a == nil {
		return
	}
	log.Infof("Accum start")
	defer log.Infof("Accum end")

	for {
		select {
		case <-ctx.Done():
			log.Infof("Accum done")
			return
		case pack := <-a.in:
			log.Infof("Accum got packet: %s", pack)
			a.processPacket(pack)
		}
	}
}
