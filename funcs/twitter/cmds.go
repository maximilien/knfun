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
	"fmt"
	"net/http"
	"os"

	"github.com/maximilien/knfun/funcs/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchFn *SearchFn
)

func NewTwitterCmd() *cobra.Command {
	searchFn = &SearchFn{
		keys: keys{},
	}

	cobra.OnInitialize(searchFn.InitConfig)

	twitterCmd := &cobra.Command{
		Use:   "twitter",
		Short: "Twitter root function",
		Long:  `Various functions over the Twitter API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			searchFn.initTwitterKeysFlags()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	searchCmd := &cobra.Command{
		Use:   "search [SEARCH_STRING]",
		Short: "Search for tweets and extract contents",
		Long: `Searches twitter for tweets matching some string criteria
and responds with with content of that tweet`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			searchFn.initTwitterKeysFlags()
			return searchFn.InitCommonInputFlags(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchFn.search(cmd, args)
		},
	}

	searchFn.AddCommonCmdFlags(searchCmd)
	searchFn.addTwitterCmdFlags(twitterCmd)

	twitterCmd.AddCommand(searchCmd)

	return twitterCmd
}

func Execute() error {
	return NewTwitterCmd().Execute()
}

// Private

func (searchFn *SearchFn) search(cmd *cobra.Command, args []string) error {
	if searchFn.StartServer {
		http.HandleFunc("/", searchFn.SearchHandler)
		return http.ListenAndServe(fmt.Sprintf(":%d", searchFn.Port), nil)
	} else {
		tweetsData, err := searchFn.Search()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", common.Flatten(&tweetsData, searchFn.Output, tweetsData.ToText))
	}

	return nil
}

func (searchFn *SearchFn) addTwitterCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&searchFn.keys.twitterAPIKey, "twitter-api-key", "", "twitter API key")
	cmd.PersistentFlags().StringVar(&searchFn.keys.twitterAPISecretKey, "twitter-api-secret-key", "", "twitter API secret key")
	cmd.PersistentFlags().StringVar(&searchFn.keys.twitterAccessToken, "twitter-access-token", "", "twitter access token")
	cmd.PersistentFlags().StringVar(&searchFn.keys.twitterAccessTokenSecret, "twitter-access-token-secret", "", "twitter access token secret")

	viper.BindPFlag("twitter-api-key", cmd.PersistentFlags().Lookup("twitter-api-key"))
	viper.BindPFlag("twitter-api-secret-key", cmd.PersistentFlags().Lookup("twitter-api-secret-key"))
	viper.BindPFlag("twitter-access-token", cmd.PersistentFlags().Lookup("twitter-access-token"))
	viper.BindPFlag("twitter-access-token-secret", cmd.PersistentFlags().Lookup("twitter-access-token-secret"))
}

func (searchFn *SearchFn) initTwitterKeysFlags() {
	if searchFn.keys.twitterAPIKey == "" {
		searchFn.keys.twitterAPIKey = viper.GetString("twitter-api-key")
	}

	if searchFn.keys.twitterAPISecretKey == "" {
		searchFn.keys.twitterAPISecretKey = viper.GetString("twitter-api-secret-key")
	}

	if searchFn.keys.twitterAccessToken == "" {
		searchFn.keys.twitterAccessToken = viper.GetString("twitter-access-token")
	}

	if searchFn.keys.twitterAccessTokenSecret == "" {
		searchFn.keys.twitterAccessTokenSecret = viper.GetString("twitter-access-token-secret")
	}
}
