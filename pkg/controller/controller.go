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

package controller

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/sunsingerus/pipeline/pkg/controller/accum"
	"github.com/sunsingerus/pipeline/pkg/controller/generator"
	"github.com/sunsingerus/pipeline/pkg/controller/generator/packet_builder"
	"github.com/sunsingerus/pipeline/pkg/controller/packet"
	"github.com/sunsingerus/pipeline/pkg/controller/pool"
	"github.com/sunsingerus/pipeline/pkg/controller/processor"
	"github.com/sunsingerus/pipeline/pkg/controller/publisher"
	mpacket "github.com/sunsingerus/pipeline/pkg/model/packet"
)

type Config struct {
	GeneratorIntervalMillisecond int
	PublisherIntervalSecond      int
	PacketSizeIn                 int
	PacketSizeOut                int
	WorkersNum                   int
}

type Controller struct {
	generatorInterval   time.Duration
	publisherInterval   time.Duration
	generatorPacketSize int
	processorPacketSize int
	workersNum          int
}

func New(conf Config) *Controller {
	return &Controller{
		generatorInterval:   time.Duration(conf.GeneratorIntervalMillisecond) * time.Millisecond,
		publisherInterval:   time.Duration(conf.PublisherIntervalSecond) * time.Second,
		generatorPacketSize: conf.PacketSizeIn,
		processorPacketSize: conf.PacketSizeOut,
		workersNum:          conf.WorkersNum,
	}
}

func (c *Controller) buildPacketBuilder() generator.PacketBuilder {
	log.Info("Building packet builder")
	return packetbuilder.New(
		func(_len int) packetbuilder.Packet {
			return mpacket.New(_len)
		},
		packetbuilder.Options{
			Size: c.generatorPacketSize,
		},
	)
}

func (c *Controller) buildGenerator() (*generator.Generator, chan packet.Packet) {
	log.Info("Building generator")
	log.Info("Making generator channel")
	ch := make(chan packet.Packet)
	gen := generator.New(
		ch,
		c.buildPacketBuilder(),
		generator.Options{
			Interval: c.generatorInterval,
		},
	)

	return gen, ch
}

func (c *Controller) buildPool(in chan packet.Packet, out chan packet.Packet) *pool.Pool {
	log.Info("Building pool")
	_pool := pool.New(
		c.workersNum,
		func(id int) pool.Processor {
			return processor.New(
				id,
				func(slice []int) processor.OutPacket {
					return mpacket.New(slice)
				},
				processor.Pipes{
					In:  in,
					Out: out,
				},
				processor.Options{
					ResultSize: c.processorPacketSize,
				},
			)
		})
	return _pool
}

func (c *Controller) buildAccum() (*accum.Accum, chan packet.Packet) {
	log.Info("Building accum")
	log.Info("Making accum channel")
	ch := make(chan packet.Packet)
	return accum.New(ch), ch

}

func (c *Controller) buildPublisher(accum *accum.Accum) *publisher.Publisher {
	log.Info("Building publilsher")
	return publisher.New(accum, publisher.Options{
		Interval: c.publisherInterval,
	})
}

func (c *Controller) build() (*generator.Generator, *pool.Pool, *accum.Accum, *publisher.Publisher, func()) {
	gen, genCh := c.buildGenerator()
	acc, accCh := c.buildAccum()
	_pool := c.buildPool(genCh, accCh)
	pub := c.buildPublisher(acc)
	return gen, _pool, acc, pub, func() {
		log.Info("Closing generator channel")
		close(genCh)
		log.Info("Closing accum channel")
		close(accCh)
	}
}

func (c *Controller) Run(ctx context.Context) (*sync.WaitGroup, func()) {
	gen, _pool, acc, pub, cancel := c.build()

	log.Info("Launching components")

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go gen.Run(ctx, wg)
	go acc.Run(ctx, wg)
	go pub.Run(ctx, wg)

	_pool.Launch(ctx, wg)
	return wg, cancel
}
