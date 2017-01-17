// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	srv "github.com/trumanw/findpro/server"
)

// Flags can be setup by command or os.ENV
var host string
var port int
var etcdns []string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start findpro server node",
	Long: `Launch the gRPC server of the findpro and register itself to etcd cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		defer glog.Flush()

		if err := srv.Run(host, port, etcdns); err != nil {
			glog.Fatal(err)
		}
	},
}

// init adds the serverCmd to RoodCmd
func init() {
	RootCmd.AddCommand(serverCmd)

	// Configuration settings.
	serverCmd.Flags().StringVarP(&host, "host", "", "", "host of the gRPC server.")
	serverCmd.Flags().IntVarP(&port, "port", "", 9090, "port of the gRPC server.")
	serverCmd.Flags().StringArrayVar(&etcdns, "etcdns", []string{"http://localhost:2379"}, "endpoints of etcd cluster.")
}
