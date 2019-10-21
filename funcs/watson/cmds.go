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
	cfgFile         string
	classifyImageFn *ClassifyImageFn
)

func NewWatsonCmd() *cobra.Command {
	classifyImageFn = &ClassifyImageFn{
		keys: watsonKeys{},
	}

	cobra.OnInitialize(initConfig)

	watsonCmd := &cobra.Command{
		Use:   "watson",
		Short: "watson root function",
		Long:  `Various functions over the Watson API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initWatsonKeys()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	vrCmd := &cobra.Command{
		Use:   "vr",
		Short: "visual recognition",
		Long:  `visual recognition Watson APIs`,
	}

	classifyCmd := &cobra.Command{
		Use:   "classify",
		Short: "classify image",
		Long:  `classify an image using the Watson APIs`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Classify(cmd, args)
		},
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Starts a server to allow search for tweets and extract contents",
		Long: `A server that allows searches twitter for tweets matching some criteria
and responds with with content of that tweet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Server(cmd, args)
		},
	}

	addWatsonCmdFlags(watsonCmd)
	addClassifyCmdFlags(classifyCmd)
	addServerCmdFlags(serverCmd)

	watsonCmd.AddCommand(vrCmd)

	vrCmd.AddCommand(classifyCmd)
	vrCmd.AddCommand(serverCmd)

	return watsonCmd
}

func Execute() error {
	return NewWatsonCmd().Execute()
}

func Classify(cmd *cobra.Command, args []string) error {
	err := initClassifyInput(args)
	if err != nil {
		return err
	}

	classifyData, err := classifyImageFn.ClassifyImage()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", classifyData.Flatten(classifyImageFn.Output))
	return nil
}

func Server(cmd *cobra.Command, args []string) error {
	err := initClassifyInput(args)
	if err != nil {
		return err
	}

	http.HandleFunc("/", classifyImageFn.ClassifyHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", classifyImageFn.Port), nil)
}

// Private

func addWatsonCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.watson.yaml)")

	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.apiKey, "api-key", "", "watson API key")
	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.apiURL, "api-url", "", "watson API URL")
	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.apiURL, "api-version", "", "watson API version")

	viper.BindPFlag("api-key", cmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("api-url", cmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("api-version", cmd.PersistentFlags().Lookup("api-version"))
}

func addClassifyCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&classifyImageFn.ImageURL, "image-url", "u", "", "the URL of the image to classify")
	cmd.Flags().StringVarP(&classifyImageFn.Output, "output", "o", "text", "the output: text, yaml, or json, of results")
}

func addServerCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&classifyImageFn.ImageURL, "image-url", "u", "", "the URL of the image to classify")
	cmd.Flags().StringVarP(&classifyImageFn.Output, "output", "o", "text", "the output: text, yaml, or json, of results")

	cmd.Flags().BoolVarP(&classifyImageFn.Server, "server", "S", false, "start as a server")
	cmd.Flags().IntVarP(&classifyImageFn.Port, "port", "p", 8080, "the port for the server")
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
		viper.SetConfigName(".watson")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initClassifyInput(args []string) error {
	if len(args) == 1 {
		classifyImageFn.ImageURL = args[0]
	}

	if classifyImageFn.ImageURL == "" {
		return errors.New(fmt.Sprintf("You must pass an image URL to classify"))
	}

	return nil
}

func initWatsonKeys() {
	if classifyImageFn.keys == (watsonKeys{}) {
		classifyImageFn.keys.apiKey = viper.GetString("api-key")
		classifyImageFn.keys.apiURL = viper.GetString("api-url")
		classifyImageFn.keys.apiVersion = viper.GetString("api-version")
	}
}
