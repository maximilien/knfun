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
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/maximilien/knfun/funcs/common"

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
		return ClassifyImageData{}, errors.New(fmt.Sprintf("Error classifying image: %s\n", err.Error()))
	}

	return classifyImageFn.collectClassifyImageData(classifiedImages), nil
}

func (classifyImageFn *ClassifyImageFn) ClassifyHandler(writer http.ResponseWriter, request *http.Request) {
	classifyImageFn.initQueryParams(request)
	log.Printf("WatsonFn.Classify: q=\"%s\", o=\"%s\"", classifyImageFn.ImageURL, classifyImageFn.Output)

	classifiedImageData, err := classifyImageFn.ClassifyImage()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	writer.Header().Add("Content-Type", classifyImageFn.OutputContentType(classifyImageFn.Output))
	fmt.Fprintf(writer, "%s\n", common.Flatten(&classifiedImageData, classifyImageFn.Output, classifiedImageData.ToText))
}

// Private classifyImageFn

func (classifyImageFn *ClassifyImageFn) initQueryParams(request *http.Request) {
	classifyImageFn.ImageURL = classifyImageFn.ExtractQueryStringParam(request, []string{"query", "q", "image-url", "u"}, classifyImageFn.ImageURL)
	classifyImageFn.Output = classifyImageFn.ExtractQueryStringParam(request, []string{"o", "output"}, classifyImageFn.Output)
}

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

// Private ClassifyImageData

func (cIData ClassifyImageData) ToText(in interface{}) string {
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
