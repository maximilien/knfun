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
	"fmt"
	"log"
	"net/http"

	"github.com/maximilien/knfun/funcs/common"
)

type keys struct {
	gVisionAPIJSON string
}

type ClassifyImageData struct {
}

type DetectLabelsFn struct {
	common.CommonFn

	ImageURL string

	keys keys
}

func (detectLabelsFn *DetectLabelsFn) ClassifyImage() (ClassifyImageData, error) {
	// 	vr, err := detectLabelsFn.createWatsonClient()
	// 	if err != nil {
	// 		return ClassifyImageData{}, err
	// 	}

	// 	classifiedImages, _, err := vr.Classify(
	// 		&vr3.ClassifyOptions{
	// 			URL: core.StringPtr(detectLabelsFn.ImageURL),
	// 		},
	// 	)
	// 	if err != nil {
	// 		return ClassifyImageData{}, errors.New(fmt.Sprintf("Error classifying image: %s\n", err.Error()))
	// 	}

	// 	return detectLabelsFn.collectClassifyImageData(classifiedImages), nil
	return ClassifyImageData{}, nil
}

func (detectLabelsFn *DetectLabelsFn) ClassifyHandler(writer http.ResponseWriter, request *http.Request) {
	detectLabelsFn.initQueryParams(request)
	log.Printf("WatsonFn.Classify: q=\"%s\", o=\"%s\"", detectLabelsFn.ImageURL, detectLabelsFn.Output)

	classifiedImageData, err := detectLabelsFn.ClassifyImage()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Headers", "x-requested-with")

	writer.Header().Add("Content-Type", detectLabelsFn.OutputContentType(detectLabelsFn.Output))
	fmt.Fprintf(writer, "%s\n", common.Flatten(&classifiedImageData, detectLabelsFn.Output, classifiedImageData.ToText))
}

// Private classifyImageFn

func (classifyImageFn *DetectLabelsFn) initQueryParams(request *http.Request) {
	// 	detectLabelsFn.ImageURL = classifyImageFn.ExtractQueryStringParam(request, []string{"query", "q", "image-url", "u"}, classifyImageFn.ImageURL)
	// 	detectLabelsFn.Output = classifyImageFn.ExtractQueryStringParam(request, []string{"o", "output"}, classifyImageFn.Output)
}

// func (detectLabelsFn *ClassifyImageFn) createWatsonClient() (*vr3.VisualRecognitionV3, error) {
// 	return vr3.NewVisualRecognitionV3(&vr3.VisualRecognitionV3Options{
// 		URL:     detectLabelsFn.keys.watsonAPIURL,
// 		Version: detectLabelsFn.keys.watsonAPIVersion,
// 		Authenticator: &core.IamAuthenticator{
// 			ApiKey: detectLabelsFn.keys.watsonAPIKey,
// 		},
// 	})
// }

// func (detectLabelsFn *ClassifyImageFn) collectClassifyImageData(classifiedImages *vr3.ClassifiedImages) ClassifyImageData {
// 	cIData := ClassifyImageData{}
// 	if *classifiedImages.ImagesProcessed >= 1 {
// 		cIData.ClassifiedImage = classifiedImages.Images[0]
// 		cIData.Warnings = classifiedImages.Warnings
// 	}

// 	return cIData
// }

// Public ClassifyImageData

func (cIData ClassifyImageData) ToText(in interface{}) string {
	sb := bytes.NewBufferString("")

	// 	sb.WriteString(fmt.Sprintf("source URL: %s\n", *cIData.ClassifiedImage.SourceURL))
	// 	sb.WriteString(fmt.Sprintf("resolved URL: %s\n", *cIData.ClassifiedImage.ResolvedURL))
	// 	for _, classifier := range cIData.ClassifiedImage.Classifiers {
	// 		sb.WriteString("----\n")
	// 		sb.WriteString(fmt.Sprintf("name: %s\n", *classifier.Name))
	// 		sb.WriteString(fmt.Sprintf("ID: %s\n", *classifier.ClassifierID))
	// 		for _, class := range classifier.Classes {
	// 			if class.Class != nil {
	// 				sb.WriteString(fmt.Sprintf("- class: %s\n", *class.Class))
	// 			}
	// 			if class.Score != nil {
	// 				sb.WriteString(fmt.Sprintf("- score: %1.3f\n", *class.Score))
	// 			}
	// 			if class.TypeHierarchy != nil {
	// 				sb.WriteString(fmt.Sprintf("- type hierarchy: %s\n", *class.TypeHierarchy))
	// 			}
	// 			sb.WriteString("-\n")
	// 		}
	// 		sb.WriteString("----\n")
	// 	}

	return sb.String()
}
