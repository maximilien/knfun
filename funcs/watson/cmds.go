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

	"github.com/maximilien/knfun/funcs/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	classifyImageFn *ClassifyImageFn
)

func NewWatsonCmd() *cobra.Command {
	classifyImageFn = &ClassifyImageFn{
		keys: keys{},
	}

	cobra.OnInitialize(classifyImageFn.InitConfig)

	watsonCmd := &cobra.Command{
		Use:   "watson",
		Short: "watson root function",
		Long:  `Various functions over the Watson API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			classifyImageFn.initWatsonKeysFlags()
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
		Use:   "classify [IMAGE_URL]",
		Short: "classify image",
		Long:  `classify an image (via its URL) using the Watson APIs`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			classifyImageFn.initWatsonKeysFlags()
			return classifyImageFn.initClassifyCmdInputFlags(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return classifyImageFn.classify(cmd, args)
		},
	}

	classifyImageFn.AddCommonCmdFlags(classifyCmd)
	classifyImageFn.addWatsonCmdFlags(watsonCmd)
	classifyImageFn.addClassifyCmdFlags(classifyCmd)

	watsonCmd.AddCommand(vrCmd)
	vrCmd.AddCommand(classifyCmd)

	return watsonCmd
}

func Execute() error {
	return NewWatsonCmd().Execute()
}

// Private

func (classifyImageFn *ClassifyImageFn) classify(cmd *cobra.Command, args []string) error {
	if classifyImageFn.StartServer {
		http.HandleFunc("/", classifyImageFn.ClassifyHandler)
		return http.ListenAndServe(fmt.Sprintf(":%d", classifyImageFn.Port), nil)
	} else {
		classifyData, err := classifyImageFn.ClassifyImage()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", common.Flatten(&classifyData, classifyImageFn.Output, classifyData.ToText))
	}

	return nil
}

func (classifyImageFn *ClassifyImageFn) addWatsonCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.watsonAPIKey, "watson-api-key", "", "watson API key")
	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.watsonAPIURL, "watson-api-url", "", "watson API URL")
	cmd.PersistentFlags().StringVar(&classifyImageFn.keys.watsonAPIVersion, "watson-api-version", "", "watson API version")

	viper.BindPFlag("watson-api-key", cmd.PersistentFlags().Lookup("watson-api-key"))
	viper.BindPFlag("watson-api-url", cmd.PersistentFlags().Lookup("watson-api-url"))
	viper.BindPFlag("watson-api-version", cmd.PersistentFlags().Lookup("watson-api-version"))
}

func (classifyImageFn *ClassifyImageFn) addClassifyCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&classifyImageFn.ImageURL, "image-url", "u", "", "the URL of the image to classify")
}

func (classifyImageFn *ClassifyImageFn) initClassifyCmdInputFlags(args []string) error {
	if len(args) == 1 {
		classifyImageFn.ImageURL = args[0]
	}

	if classifyImageFn.ImageURL == "" {
		return errors.New(fmt.Sprintf("You must pass an image URL to classify"))
	}

	return nil
}

func (classifyImageFn *ClassifyImageFn) initWatsonKeysFlags() {
	if classifyImageFn.keys.watsonAPIKey == "" {
		classifyImageFn.keys.watsonAPIKey = viper.GetString("watson-api-key")
	}

	if classifyImageFn.keys.watsonAPIURL == "" {
		classifyImageFn.keys.watsonAPIURL = viper.GetString("watson-api-url")
	}

	if classifyImageFn.keys.watsonAPIVersion == "" {
		classifyImageFn.keys.watsonAPIVersion = viper.GetString("watson-api-version")
	}
}
