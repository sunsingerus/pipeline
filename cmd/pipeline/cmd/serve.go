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

package cmd

import (
	"context"
	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	cmd "github.com/spf13/cobra"
	vprConfig "github.com/spf13/viper"
	"github.com/sunsingerus/pipeline/pkg/controller"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	generatorIntervalMillisecond int
	publisherIntervalSecond      int
	packetSizeIn                 int
	packetSizeOut                int
	workersNum                   int
	runTimeoutSecond             int
)

var serveCmd = &cmd.Command{
	Use:   "serve [OPTION(s)] ",
	Short: "Serve pipelines",
	Long:  heredoc.Docf(`Serve pipeline service`),
	Args: func(cmd *cmd.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cmd.Command, args []string) {
		// Init termination context
		ctx := contextInit()

		log.Infof("Starting service")
		log.Infof(heredoc.Docf(`
			Options:
			----------------------------
			generator-interval (ms)    : %d
			publisher-interval (s)     : %d
			packet-size-in     (items) : %d
			packet-size-out    (items) : %d
			workers            (num)   : %d
			timeout            (s)     : %d
			----------------------------
`, generatorIntervalMillisecond, publisherIntervalSecond, packetSizeIn, packetSizeOut, workersNum, runTimeoutSecond))

		if runTimeoutSecond > 0 {
			log.Infof("Will run for %d sec", runTimeoutSecond)
			var fn context.CancelFunc
			ctx, fn = context.WithTimeout(ctx, time.Duration(runTimeoutSecond)*time.Second)
			defer fn()
		}
		wg, cancel := run(ctx)
		contextWait(ctx)
		wg.Wait()
		cancel()
		log.Info("Shut down")
	},
}

func init() {
	flagInit()
	// Options (CLI+ENV)
	pFlagInt(serveCmd, "generator-interval", "g", "interval in microseconds between packets produced by the generator", 1000, &generatorIntervalMillisecond)
	pFlagInt(serveCmd, "publisher-interval", "p", "interval in seconds between publisher reports", 1, &publisherIntervalSecond)
	pFlagInt(serveCmd, "packet-size-in", "s", "size of the generated packet", 10, &packetSizeIn)
	pFlagInt(serveCmd, "packet-size-out", "o", "size of the processed packet", 3, &packetSizeOut)
	pFlagInt(serveCmd, "workers", "w", "workers number", 3, &workersNum)
	pFlagInt(serveCmd, "timeout", "t", "run timeout in seconds (default unlimited)", 0, &runTimeoutSecond)

	// Bind full flag set to the configuration
	if err := vprConfig.BindPFlags(serveCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(serveCmd)
}

// run runs service
func run(ctx context.Context) (*sync.WaitGroup, func()) {
	log.Infof("run() - start")
	defer log.Infof("run() - end")

	return controller.New(controller.Config{
		GeneratorIntervalMillisecond: generatorIntervalMillisecond,
		PublisherIntervalSecond:      publisherIntervalSecond,
		PacketSizeIn:                 packetSizeIn,
		PacketSizeOut:                packetSizeOut,
		WorkersNum:                   workersNum,
	}).Run(ctx)
}

func contextInit() context.Context {
	// Set OS signals and termination context
	ctx, cancelFunc := context.WithCancel(context.Background())
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		cancelFunc()
		<-stopChan
	}()

	return ctx
}

// contextWait
func contextWait(ctx context.Context) {
	<-ctx.Done()
	log.Info("Shutting down...")
}
