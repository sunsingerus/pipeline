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
	"fmt"
	"math/rand"
)

type Packet interface {
	fmt.Stringer
	Len() int
	Set(int, int)
}

type packetConstructor func(int) Packet

// Options specifies generator options
type Options struct {
	// Size specifies size of generated packet
	Size int
}

// PacketBuilder specifies packet builder
type PacketBuilder struct {
	packetConstructor packetConstructor
	Options
}

// New creates new packet builder from options
func New(packetConstructor packetConstructor, opts Options) *PacketBuilder {
	return &PacketBuilder{
		packetConstructor: packetConstructor,
		Options:           opts,
	}
}

// Build builds one packet
func (b *PacketBuilder) Build() Packet {
	if b == nil {
		return nil
	}

	// Create randomly filled packet
	packet := b.packetConstructor(b.Options.Size)
	for i := 0; i < packet.Len(); i++ {
		packet.Set(i, rand.Intn(20))
	}
	return packet
}
