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
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sunsingerus/pipeline/pkg/controller/generator/packet_builder"
	"github.com/sunsingerus/pipeline/pkg/controller/packet"
)

type PacketBuilder interface {
	Build() packetbuilder.Packet
}

// Options specifies generator options
type Options struct {
	// Interval specifies interval between packet generations (in milliseconds)
	Interval time.Duration
}

// Generator specifies generator
type Generator struct {
	// out specifies chan where generator puts generated packet
	out           chan packet.Packet
	packetBuilder PacketBuilder
	Options
}

// New creates new generator from options
func New(out chan packet.Packet, packetBuilder PacketBuilder, opts Options) *Generator {
	return &Generator{
		out:           out,
		packetBuilder: packetBuilder,
		Options:       opts,
	}
}

func (g *Generator) deliver(ctx context.Context, pack packetbuilder.Packet) {
	if g == nil {
		return
	}
	select {
	case <-ctx.Done():
		log.Infof("Generator - NODELIVERY: %s", pack)
	case g.out <- pack.(packet.Packet):
		log.Infof("Generator - delivered : %s", pack)
	}
}

// Run runs generator until context is done
func (g *Generator) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if g == nil {
		return
	}
	log.Infof("Generator start")
	defer log.Infof("Generator end")

	ticker := time.NewTicker(g.Options.Interval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Infof("Generator - done")
			return
		case at := <-ticker.C:
			pack := g.packetBuilder.Build()
			log.Infof("Generator - new packet: %s @[%s]", pack, at)
			g.deliver(ctx, pack)
		}
	}
}
