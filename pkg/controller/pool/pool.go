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

package pool

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Pool struct {
	size                 int
	processorConstructor processorConstructor
}

type Processor interface {
	Process(ctx context.Context)
}

type processorConstructor func(id int) Processor

func New(size int, processorConstructor processorConstructor) *Pool {
	return &Pool{
		size:                 size,
		processorConstructor: processorConstructor,
	}
}

func (p *Pool) launch(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	if p == nil {
		return
	}
	log.Infof("Launcher  [%d] - start", id)
	defer log.Infof("Launcher  [%d] - end", id)
	p.processorConstructor(id).Process(ctx)
}

func (p *Pool) Launch(ctx context.Context, wg *sync.WaitGroup) {
	if p == nil {
		return
	}
	log.Infof("Launch - start")
	defer log.Infof("Launch - end")

	for i := 0; i < p.size; i++ {
		wg.Add(1)
		go p.launch(ctx, wg, i)
	}
}
