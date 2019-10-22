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
	cfgFile   string
	summaryFn *SummaryFn
)

func NewSummaryCmd() *cobra.Command {
	summaryFn = &SummaryFn{}

	cobra.OnInitialize(initConfig)

	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Summary root function",
		Long:  `Summarizes the TwitterFn and WatsonFn APIs funcs`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}

			return summary(cmd, args)
		},
	}

	addSummaryCmdFlags(summaryCmd)

	return summaryCmd
}

func Execute() error {
	return NewSummaryCmd().Execute()
}

// Private

func summary(cmd *cobra.Command, args []string) error {
	err := initInput(args)
	if err != nil {
		return err
	}

	if summaryFn.StartServer {
		http.HandleFunc("/", summaryFn.SummaryHandler)
		return http.ListenAndServe(fmt.Sprintf(":%d", summaryFn.Port), nil)
	} else {
		classifiedTweets, err := summaryFn.Summary()
		if err != nil {
			return err
		}

		for _, cTweet := range classifiedTweets {
			fmt.Printf("%s\n", cTweet.Flatten(summaryFn.Output))
			fmt.Printf("=======\n\n")
		}
	}

	return nil
}

func addSummaryCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.twiter.yaml)")

	cmd.PersistentFlags().StringVar(&summaryFn.TwitterFnURL, "twitter-fn-url", "", "twitter API func URL")
	cmd.PersistentFlags().StringVar(&summaryFn.WatsonFnURL, "watson-fn-url", "", "watson API func URL")

	viper.BindPFlag("twitter-fn-url", cmd.PersistentFlags().Lookup("twitter-fn-url"))
	viper.BindPFlag("watson-fn-url", cmd.PersistentFlags().Lookup("watson-fn-url"))

	cmd.Flags().StringVarP(&summaryFn.SearchString, "search-string", "s", "", "the string to search for")
	cmd.Flags().IntVarP(&summaryFn.Count, "count", "c", 10, "the max number of results")
	cmd.Flags().IntVarP(&summaryFn.Timeout, "timeout", "t", 60, "the timeout in seconds for requests to functions")

	cmd.Flags().StringVarP(&summaryFn.Output, "output", "o", "text", "the output: text, yaml, or json, of results")

	cmd.Flags().BoolVarP(&summaryFn.StartServer, "start-server", "S", false, "start as a server")
	cmd.Flags().IntVarP(&summaryFn.Port, "port", "p", 8080, "the port for the server")
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
		viper.SetConfigName(".summary")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initInput(args []string) error {
	if len(args) == 1 {
		summaryFn.SearchString = args[0]
	}

	if summaryFn.SearchString == "" {
		return errors.New(fmt.Sprintf("You must pass a search string"))
	}

	if summaryFn.TwitterFnURL == "" {
		summaryFn.TwitterFnURL = viper.GetString("twitter-fn-url")
	}

	if summaryFn.WatsonFnURL == "" {
		summaryFn.WatsonFnURL = viper.GetString("watson-fn-url")
	}

	return nil
}
