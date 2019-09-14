// Copyright © 2019 NAME HERE scbizu@gmail.com
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
	"fmt"
	"os"

	"github.com/CNSC2Events/feeds/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var port int32
var debugMode bool

func init() {
	if debugMode {
		EnableDebug()
	}
}

func EnableDebug() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "feeds",
	Short: "tl.net feed service",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		log.Info().Msg("cache: cache will be registered")
		service.RegisterCache(ctx)
		log.Info().Msgf("service: feed service will start at port %d", port)
		if err := service.NewFeedService(port).Serve(ctx); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	RootCmd.Flags().Int32VarP(&port, "port", "p", 8888, "service port")
	RootCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "open debug mode")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
