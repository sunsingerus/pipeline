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

package publisher

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type accum interface {
	Get() int
}

// Options specifies generator options
type Options struct {
	// Interval specifies interval between packet publications
	Interval time.Duration
}

// Publisher specifies publisher
type Publisher struct {
	accum accum
	Options
}

// New creates new publisher from options
func New(accum accum, opts Options) *Publisher {
	return &Publisher{
		accum:   accum,
		Options: opts,
	}
}

// Run runs generator until context is done
func (p *Publisher) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	if p == nil {
		return
	}
	log.Infof("Publisher - start")
	defer log.Infof("Publisher - end")

	ticker := time.NewTicker(p.Options.Interval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Infof("Publisher - done")
			return
		case at := <-ticker.C:
			log.Infof("Publisher: %d @[%s]", p.accum.Get(), at)
		}
	}
}
