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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	summaryFn *SummaryFn
)

func NewSummaryCmd() *cobra.Command {
	summaryFn = &SummaryFn{}

	cobra.OnInitialize(summaryFn.InitConfig)

	summaryCmd := &cobra.Command{
		Use:   "summary [SEARCH_STRING]",
		Short: "Summary root function",
		Long: `Summarizes the TwitterFn and WatsonFn APIs funcs
searching for tweets (with images) that contains SEARCH_STRING`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			summaryFn.initInputFlags(args)
			return summaryFn.InitCommonInputFlags(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}

			return summaryFn.summary(cmd, args)
		},
	}

	summaryFn.AddCommonCmdFlags(summaryCmd)
	summaryFn.addSummaryCmdFlags(summaryCmd)

	return summaryCmd
}

func Execute() error {
	return NewSummaryCmd().Execute()
}

// Private

func (summaryFn *SummaryFn) summary(cmd *cobra.Command, args []string) error {
	if summaryFn.StartServer {
		if os.Getenv("ASYNC") != "" {
			http.HandleFunc("/", summaryFn.SummaryAsyncHandler)
		} else {
			http.HandleFunc("/", summaryFn.SummaryHandler)
		}

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

func (summaryFn *SummaryFn) addSummaryCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&summaryFn.TwitterFnURL, "twitter-fn-url", "", "twitter API func URL")
	cmd.PersistentFlags().StringVar(&summaryFn.WatsonFnURL, "watson-fn-url", "", "watson API func URL")

	viper.BindPFlag("twitter-fn-url", cmd.PersistentFlags().Lookup("twitter-fn-url"))
	viper.BindPFlag("watson-fn-url", cmd.PersistentFlags().Lookup("watson-fn-url"))
}

func (summaryFn *SummaryFn) initInputFlags(args []string) {
	if summaryFn.TwitterFnURL == "" {
		summaryFn.TwitterFnURL = viper.GetString("twitter-fn-url")
	}

	if summaryFn.WatsonFnURL == "" {
		summaryFn.WatsonFnURL = viper.GetString("watson-fn-url")
	}
}
