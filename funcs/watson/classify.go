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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/maximilien/knfun/funcs/common"

	"gopkg.in/yaml.v2"

	"github.com/IBM/go-sdk-core/core"
	vr3 "github.com/watson-developer-cloud/go-sdk/visualrecognitionv3"
)

type keys struct {
	watsonAPIKey     string
	watsonAPIURL     string
	watsonAPIVersion string
}

type ClassifyImageData struct {
	vr3.ClassifiedImage
	Warnings []vr3.WarningInfo
}

type ClassifyImageFn struct {
	common.CommonFn

	ImageURL string

	keys keys
}

func (classifyImageFn *ClassifyImageFn) ClassifyImage() (ClassifyImageData, error) {
	vr, err := classifyImageFn.createWatsonClient()
	if err != nil {
		return ClassifyImageData{}, err
	}

	classifiedImages, _, err := vr.Classify(
		&vr3.ClassifyOptions{
			URL: core.StringPtr(classifyImageFn.ImageURL),
		},
	)
	if err != nil {
		log.Fatal("Error classifying image: %s\n", err.Error())
		return ClassifyImageData{}, err
	}

	return classifyImageFn.collectClassifyImageData(classifiedImages), nil
}

func (classifyImageFn *ClassifyImageFn) ClassifyHandler(writer http.ResponseWriter, request *http.Request) {
	imageURL, output := classifyImageFn.extractQueryParams(request)

	if imageURL == "" {
		log.Printf("Must pass an image URL string using 'q' or 'query' parameter")
		return
	}

	vr, err := classifyImageFn.createWatsonClient()
	if err != nil {
		log.Printf("Error creating visual recognition client: %s\n", err.Error())
		return
	}

	classifiedImages, _, err := vr.Classify(
		&vr3.ClassifyOptions{
			URL: core.StringPtr(imageURL),
		},
	)
	if err != nil {
		log.Printf("Error classifying image: %s\n", err.Error())
		return
	}

	writer.Header().Add("Content-Type", classifyImageFn.outputContentType(output))
	fmt.Fprintf(writer, "%s\n", classifyImageFn.collectClassifyImageData(classifiedImages).Flatten(output))
}

// Private classifyImageFn

func (classifyImageFn *ClassifyImageFn) createWatsonClient() (*vr3.VisualRecognitionV3, error) {
	return vr3.NewVisualRecognitionV3(&vr3.VisualRecognitionV3Options{
		URL:     classifyImageFn.keys.watsonAPIURL,
		Version: classifyImageFn.keys.watsonAPIVersion,
		Authenticator: &core.IamAuthenticator{
			ApiKey: classifyImageFn.keys.watsonAPIKey,
		},
	})
}

func (classifyImageFn *ClassifyImageFn) collectClassifyImageData(classifiedImages *vr3.ClassifiedImages) ClassifyImageData {
	cIData := ClassifyImageData{}
	if *classifiedImages.ImagesProcessed >= 1 {
		cIData.ClassifiedImage = classifiedImages.Images[0]
		cIData.Warnings = classifiedImages.Warnings
	}

	return cIData
}

func (classifyImageFn *ClassifyImageFn) extractQueryParams(request *http.Request) (string, string) {
	imageURL, output := classifyImageFn.ImageURL, classifyImageFn.Output

	query := request.URL.Query()

	if query.Get("q") != "" {
		imageURL = query.Get("q")
	} else {
		if query.Get("query") != "" {
			imageURL = query.Get("query")
		}
	}

	if query.Get("o") != "" {
		output = query.Get("o")
	} else {
		if query.Get("ouput") != "" {
			output = query.Get("output")
		}
	}

	return imageURL, output
}

func (classifyImageFn *ClassifyImageFn) outputContentType(output string) string {
	switch output {
	case "yaml":
		return "application/yaml"
	case "json":
		return "application/json"
	}
	return "text/html"
}

// Private ClassifyImageData

func (cIData ClassifyImageData) ToText() string {
	sb := bytes.NewBufferString("")

	sb.WriteString(fmt.Sprintf("source URL: %s\n", *cIData.ClassifiedImage.SourceURL))
	sb.WriteString(fmt.Sprintf("resolved URL: %s\n", *cIData.ClassifiedImage.ResolvedURL))
	for _, classifier := range cIData.ClassifiedImage.Classifiers {
		sb.WriteString("----\n")
		sb.WriteString(fmt.Sprintf("name: %s\n", *classifier.Name))
		sb.WriteString(fmt.Sprintf("ID: %s\n", *classifier.ClassifierID))
		for _, class := range classifier.Classes {
			if class.Class != nil {
				sb.WriteString(fmt.Sprintf("- class: %s\n", *class.Class))
			}
			if class.Score != nil {
				sb.WriteString(fmt.Sprintf("- score: %1.3f\n", *class.Score))
			}
			if class.TypeHierarchy != nil {
				sb.WriteString(fmt.Sprintf("- type hierarchy: %s\n", *class.TypeHierarchy))
			}
			sb.WriteString("-\n")
		}
		sb.WriteString("----\n")
	}

	return sb.String()
}

func (cIData ClassifyImageData) ToYAML() string {
	data, err := yaml.Marshal(&cIData)
	if err != nil {
		panic("Error YAML marshalling ClassifyImageData")
	}

	return string(data)
}

func (cIData ClassifyImageData) ToJSON() string {
	data, err := json.MarshalIndent(&cIData, "", "  ")
	if err != nil {
		panic("Error JSON marshalling ClassifyImageData")
	}

	return string(data)
}

func (cIData ClassifyImageData) Flatten(output string) string {
	outputData := ""
	switch output {
	case "yaml":
		outputData = cIData.ToYAML()
		break
	case "json":
		outputData = cIData.ToJSON()
		break
	default:
		outputData = cIData.ToText()
	}
	return outputData
}
