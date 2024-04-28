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
	"bytes"
	"strconv"
)

type Packet []int

func New(what any) *Packet {
	switch _typed := what.(type) {
	case int:
		return newFromLen(_typed)
	case []int:
		return newFromSlice(_typed)
	}
	return nil
}

func newFromLen(_len int) *Packet {
	packet := make(Packet, _len)
	return &packet
}

func newFromSlice(slice []int) *Packet {
	packet := Packet(slice)
	return &packet
}

func (p *Packet) Len() int {
	if p == nil {
		return 0
	}
	return len(*p)
}

func (p *Packet) Set(i int, value int) {
	if p == nil {
		return
	}
	(*p)[i] = value
}

func (p *Packet) Get(i int) int {
	if p == nil {
		return 0
	}
	return (*p)[i]
}

func (p *Packet) Slice(boundaries ...int) []int {
	if p == nil {
		return nil
	}
	switch len(boundaries) {
	// Whole packet
	case 0:
		return (*p)[:]
	// Left or Right side of the slice
	case 1:
		boundary := boundaries[0]
		if boundary < 0 {
			return (*p)[p.Len()+boundary:]
		} else {
			return (*p)[boundary:]
		}

	// Slice by specified boundaries
	default:
		return (*p)[boundaries[0]:boundaries[1]]
	}
}

func (p *Packet) String() string {
	if p == nil {
		return ""
	}

	var str bytes.Buffer
	str.WriteString("[")
	for i := 0; i < len(*p); i++ {
		if i > 0 {
			str.WriteString(",")
		}
		str.WriteString(strconv.Itoa((*p)[i]))
	}
	str.WriteString("]")
	return str.String()
}
