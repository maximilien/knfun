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
	detectLabelsFn *DetectLabelsFn
)

func NewGVisionCmd() *cobra.Command {
	detectLabelsFn = &DetectLabelsFn{
		keys: keys{},
	}

	cobra.OnInitialize(detectLabelsFn.InitConfig)

	gVisionCmd := &cobra.Command{
		Use:   "gvision",
		Short: "Google vision root function",
		Long:  `Various functions over the Google Vision API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			detectLabelsFn.initGVisionKeysFlags()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	detectLabelsCmd := &cobra.Command{
		Use:   "dl [IMAGE_URL]",
		Short: "detect labels image",
		Long:  `detect labels (classify) an image (via its URL) using the GVision APIs`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			detectLabelsFn.initGVisionKeysFlags()
			return detectLabelsFn.initDetectLabelsCmdInputFlags(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return detectLabelsFn.detectLabels(cmd, args)
		},
	}

	detectLabelsFn.AddCommonCmdFlags(detectLabelsCmd)
	detectLabelsFn.addGVisionCmdFlags(gVisionCmd)
	detectLabelsFn.addDetectLabelsCmdFlags(detectLabelsCmd)

	gVisionCmd.AddCommand(detectLabelsCmd)

	return gVisionCmd
}

func Execute() error {
	return NewGVisionCmd().Execute()
}

// Private

func (detectLabelsFn *DetectLabelsFn) detectLabels(cmd *cobra.Command, args []string) error {
	if detectLabelsFn.StartServer {
		http.HandleFunc("/", detectLabelsFn.ClassifyHandler)
		return http.ListenAndServe(fmt.Sprintf(":%d", detectLabelsFn.Port), nil)
	} else {
		classifyData, err := detectLabelsFn.ClassifyImage()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", common.Flatten(&classifyData, detectLabelsFn.Output, classifyData.ToText))
	}

	return nil
}

func (detectLabelsFn *DetectLabelsFn) addGVisionCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&detectLabelsFn.keys.gVisionAPIJSON, "watson-api-json", "", "watson API JSON")

	viper.BindPFlag("watson-api-json", cmd.PersistentFlags().Lookup("watson-api-json"))
}

func (detectLabelsFn *DetectLabelsFn) addDetectLabelsCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&detectLabelsFn.ImageURL, "image-url", "u", "", "the URL of the image to detect labels")
}

func (detectLabelsFn *DetectLabelsFn) initDetectLabelsCmdInputFlags(args []string) error {
	if len(args) == 1 {
		detectLabelsFn.ImageURL = args[0]
	}

	if detectLabelsFn.ImageURL == "" {
		return errors.New(fmt.Sprintf("You must pass an image URL to detect labels"))
	}

	return nil
}

func (detectLabelsFn *DetectLabelsFn) initGVisionKeysFlags() {
	// if detectLabelsFn.keys.watsonAPIKey == "" {
	// 	detectLabelsFn.keys.watsonAPIKey = viper.GetString("watson-api-key")
	// }

	// if detectLabelsFn.keys.watsonAPIURL == "" {
	// 	detectLabelsFn.keys.watsonAPIURL = viper.GetString("watson-api-url")
	// }

	// if detectLabelsFn.keys.watsonAPIVersion == "" {
	// 	detectLabelsFn.keys.watsonAPIVersion = viper.GetString("watson-api-version")
	// }
}
