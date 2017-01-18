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
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	gw "github.com/trumanw/findpro/gateway"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start findpro gRPC gateway.",
	Long: `Launch the gRPC gateway and parse the endpoints of gRPC servers from the etcd cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		defer glog.Flush()

		if err := gw.Run(etcdns); err != nil {
			glog.Fatal(err)
		}
	},
}

// init adds the proxyCmd to RootCmd
func init() {
	RootCmd.AddCommand(proxyCmd)
}
