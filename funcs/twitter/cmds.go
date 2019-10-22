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

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	cfgFile  string
	searchFn *SearchFn
)

func NewTwitterCmd() *cobra.Command {
	searchFn = &SearchFn{
		keys: twitterKeys{},
	}

	cobra.OnInitialize(initConfig)

	twitterCmd := &cobra.Command{
		Use:   "twitter",
		Short: "Twitter root function",
		Long:  `Various functions over the Twitter API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initTwitterKeys()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for tweets and extract contents",
		Long: `Searches twitter for tweets matching some criteria
and responds with with content of that tweet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(cmd, args)
		},
	}

	addTwitterCmdFlags(twitterCmd)
	addSearchCmdFlags(searchCmd)

	twitterCmd.AddCommand(searchCmd)

	return twitterCmd
}

func Execute() error {
	return NewTwitterCmd().Execute()
}

// Private

func search(cmd *cobra.Command, args []string) error {
	err := initSearchInput(args)
	if err != nil {
		return err
	}

	if searchFn.StartServer {
		http.HandleFunc("/", searchFn.SearchHandler)
		return http.ListenAndServe(fmt.Sprintf(":%d", searchFn.Port), nil)
	} else {
		tweetsData, err := searchFn.Search()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", tweetsData.Flatten(searchFn.Output))
	}

	return nil
}

func addTwitterCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.twiter.yaml)")

	cmd.PersistentFlags().StringVar(&searchFn.keys.apiKey, "api-key", "", "twitter API key")
	cmd.PersistentFlags().StringVar(&searchFn.keys.apiSecretKey, "api-secret-key", "", "twitter API secret key")
	cmd.PersistentFlags().StringVar(&searchFn.keys.accessToken, "access-token", "", "twitter access token")
	cmd.PersistentFlags().StringVar(&searchFn.keys.accessTokenSecret, "access-token-secret", "", "twitter access token secret")

	viper.BindPFlag("api-key", cmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("api-secret-key", cmd.PersistentFlags().Lookup("api-secret-key"))
	viper.BindPFlag("access-token", cmd.PersistentFlags().Lookup("access-token"))
	viper.BindPFlag("access-token-secret", cmd.PersistentFlags().Lookup("access-token-secret"))
}

func addSearchCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&searchFn.String, "string", "s", "", "the string to search for")
	cmd.Flags().IntVarP(&searchFn.Count, "count", "c", 10, "the max number of results")

	cmd.Flags().StringVarP(&searchFn.Output, "output", "o", "text", "the output: text, yaml, or json, of results")

	cmd.Flags().BoolVarP(&searchFn.StartServer, "start-server", "S", false, "start as a server")
	cmd.Flags().IntVarP(&searchFn.Port, "port", "p", 8080, "the port for the server")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			os.Exit(-1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".twitter")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initSearchInput(args []string) error {
	if len(args) == 1 {
		searchFn.String = args[0]
	}

	if searchFn.String == "" {
		return errors.New(fmt.Sprintf("You must pass a search string"))
	}

	return nil
}

func initTwitterKeys() {
	if searchFn.keys == (twitterKeys{}) {
		searchFn.keys.apiKey = viper.GetString("api-key")
		searchFn.keys.apiSecretKey = viper.GetString("api-secret-key")
		searchFn.keys.accessToken = viper.GetString("access-token")
		searchFn.keys.accessTokenSecret = viper.GetString("access-token-secret")
	}
}
