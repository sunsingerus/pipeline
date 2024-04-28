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
	"strings"

	cmd "github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pFlagString creates persistent flag with string value
func pFlagString(cmd *cmd.Command, name, short, description string, defaultValue string, variable *string) {
	viper.SetDefault(name, defaultValue)
	cmd.PersistentFlags().StringVarP(variable, name, short, viper.GetString(name), description)
}

// pFlagInt creates persistent flag with int value
func pFlagInt(cmd *cmd.Command, name, short, description string, defaultValue int, variable *int) {
	viper.SetDefault(name, defaultValue)
	cmd.PersistentFlags().IntVarP(variable, name, short, viper.GetInt(name), description)
}

// flagInit initializes flag components
func flagInit() {
	// By default, empty environment variables are considered unset and will fall back to the next configuration source.
	// To treat empty environment variables as set, use the AllowEmptyEnv method.
	viper.AllowEmptyEnv(false)
	// Check for an env var with a name matching the key upper-cased and prefixed with the EnvPrefix
	// Prefix has "_" added automatically, so no need to say 'TBOX_'
	// viper.SetEnvPrefix()
	// SetEnvKeyReplacer allows you to use a strings.Replacer object to rewrite Env keys to an extent.
	// This is useful if you want to use - or something in your Get() calls, but want your environmental variables to use _ delimiters.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Check ENV variables for all keys set in config, default & flags
	viper.AutomaticEnv()
}
