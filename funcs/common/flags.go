// Copyright © 2019 The Knative Authors
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

package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

type CommonFn struct {
	CfgFile string

	SearchString string
	Count        int

	Output string

	Timeout     int
	StartServer bool
	Port        int
}

func (commonFn *CommonFn) AddCommonCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&commonFn.CfgFile, "config", "", "config file (default is $HOME/.twiter.yaml)")

	cmd.Flags().StringVarP(&commonFn.SearchString, "search-string", "s", "", "the string to search for")
	cmd.Flags().IntVarP(&commonFn.Count, "count", "c", 10, "the max number of results")

	cmd.Flags().StringVarP(&commonFn.Output, "output", "o", "text", "the output: text, yaml, or json, of results")

	cmd.Flags().BoolVarP(&commonFn.StartServer, "start-server", "S", false, "start as a server")
	cmd.Flags().IntVarP(&commonFn.Port, "port", "p", 8080, "the port for the server")
}

func (commonFn *CommonFn) InitCommonInputFlags(args []string) error {
	if len(args) == 1 {
		commonFn.SearchString = args[0]
	}

	if commonFn.SearchString == "" {
		return errors.New("you must pass a search string")
	}

	return nil
}

func (commonFn *CommonFn) InitConfig() {
	if commonFn.CfgFile != "" {
		viper.SetConfigFile(commonFn.CfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			os.Exit(-1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".knfun")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
