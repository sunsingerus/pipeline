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

package logger

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	// Level specifies verbosity level in string form
	// Available levels are:
	// "panic"
	// "fatal"
	// "error"
	// "warn", "warning"
	// "info"
	// "debug"
	// "trace"
	Level string

	// Formatter specifies log format
	// Available formatters are:
	// "json"
	// "text", "txt"
	Formatter string
)

// Init sets logging options
func Init() {
	if formatter, err := parseFormatter(Formatter); err == nil {
		log.SetFormatter(formatter)
		log.Infof("Set formatter: %s", Formatter)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.Infof("Set default formatter - text")
	}

	if level, err := log.ParseLevel(Level); err == nil {
		log.SetLevel(level)
		log.Infof("Set log level: %s", level)
	} else {
		log.SetLevel(log.InfoLevel)
		log.Infof("Set default log level - Info")
	}
}

// parseFormatter makes Formatter out of its string name
func parseFormatter(str string) (log.Formatter, error) {
	switch strings.ToLower(str) {
	case "json":
		return &log.JSONFormatter{}, nil
	case "txt", "text":
		return &log.TextFormatter{}, nil
	}

	return nil, fmt.Errorf("not a valid logrus formatter: %q", str)
}
