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
	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	cmd "github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sunsingerus/pipeline/pkg/logger"
)

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cmd.Command{
	Use:   "pipeline [COMMAND]",
	Short: "Pipeline service.",
	Long:  heredoc.Docf(`Pipeline service is used to serve pipelines. Pipeline with caution.`),
	PersistentPreRun: func(cmd *cmd.Command, args []string) {
		logger.Init()
		log.Infof(heredoc.Docf(`
				Starting root:
				log-level: %s
				log-format: %s
			`, logger.Level, logger.Formatter))
	},
}

func init() {
	flagInit()
	// Options (CLI+ENV)
	pFlagString(rootCmd, "log-level", "l", "log level, one of: panic,fatal,error,warn,warning,info,debug,trace", "info", &logger.Level)
	pFlagString(rootCmd, "log-format", "f", "log format, one of: text,json", "text", &logger.Formatter)

	// Bind full flag set to the configuration
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
