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

package processor

import (
	"context"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/sunsingerus/pipeline/pkg/controller/packet"
)

type inPacket interface {
	fmt.Stringer
	Slice(...int) []int
	Len() int
}

type OutPacket interface {
	fmt.Stringer
}

type outPacketConstructor func([]int) OutPacket

type Pipes struct {
	In  chan packet.Packet
	Out chan packet.Packet
}

// Options specifies processor options
type Options struct {
	// Size specifies size of result packet
	ResultSize int
}

type Processor struct {
	id                   int
	outPacketConstructor outPacketConstructor
	Pipes
	Options
}

func New(id int, outPacketConstructor outPacketConstructor, pipes Pipes, opts Options) *Processor {
	return &Processor{
		id:                   id,
		outPacketConstructor: outPacketConstructor,
		Pipes:                pipes,
		Options:              opts,
	}
}

func (p *Processor) processPacket(in inPacket) OutPacket {
	if p == nil {
		return nil
	}

	// Sort incoming packet ASC and get N greatest numbers, which will be located on the right side of the slice
	sort.Ints(in.Slice())
	return p.outPacketConstructor(in.Slice(-p.Options.ResultSize))
}

func (p *Processor) deliver(ctx context.Context, pack OutPacket) {
	if p == nil {
		return
	}
	select {
	case <-ctx.Done():
		log.Infof("Processor [%d] - NODELIVERY: %s", p.id, pack)
	case p.Pipes.Out <- pack.(packet.Packet):
		log.Infof("Processor [%d] - delivered : %s", p.id, pack)
	}
}

func (p *Processor) Process(ctx context.Context) {
	log.Infof("Processor [%d] - start", p.id)
	defer log.Infof("Processor [%d] - end", p.id)

	for {
		select {
		case <-ctx.Done():
			log.Infof("Processor [%d] - done", p.id)
			return
		case pack := <-p.Pipes.In:
			log.Infof("Processor [%d] - received  : %s", p.id, pack)
			result := p.processPacket(pack)
			log.Infof("Processor [%d] - prepared  : %s", p.id, result)
			p.deliver(ctx, result)
		}
	}
}
