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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/maximilien/knfun/funcs/common"

	vision "cloud.google.com/go/vision/apiv1"
)

type keys struct {
	gVisionAPIJSON string
}

type Label struct {
	Name  string
	Score float32
}

type ClassifyImageData struct {
	ImageURL string
	Labels   []Label
}

type DetectLabelsFn struct {
	common.CommonFn

	client *vision.ImageAnnotatorClient

	ImageURL string

	keys keys
}

func (detectLabelsFn *DetectLabelsFn) ClassifyImage() (ClassifyImageData, error) {
	ctx := context.Background()

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", detectLabelsFn.keys.gVisionAPIJSON)
		if err != nil {
			return ClassifyImageData{}, fmt.Errorf("error finding GOOGLE_APPLICATION_CREDENTIALS enviroment variable: %s", err.Error())
		}
	}

	var err error
	if detectLabelsFn.client == nil {
		gVisionClient, err := vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return ClassifyImageData{}, fmt.Errorf("error creating client: %s", err.Error())
		}
		detectLabelsFn.client = gVisionClient
	}

	filepath := detectLabelsFn.ImageURL
	if strings.HasPrefix(detectLabelsFn.ImageURL, "http") {
		filepath, err = common.DownloadTmpFile(detectLabelsFn.ImageURL)
		if err != nil {
			return ClassifyImageData{}, fmt.Errorf("error creating tmp file for image: %s", err.Error())
		}
	}

	file, err := os.Open(filepath)
	if err != nil {
		return ClassifyImageData{}, fmt.Errorf("error loading image: %s", err.Error())
	}
	defer file.Close()
	defer os.Remove(filepath)

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return ClassifyImageData{}, fmt.Errorf("error reading image: %s", err.Error())
	}

	labels, err := detectLabelsFn.client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return ClassifyImageData{}, fmt.Errorf("error detecting labels for image: %s", err.Error())
	}

	cImageData := ClassifyImageData{
		ImageURL: detectLabelsFn.ImageURL,
	}
	for _, label := range labels {
		l := Label{
			Name:  label.Description,
			Score: label.Score,
		}
		cImageData.Labels = append(cImageData.Labels, l)
	}

	return cImageData, nil
}

func (detectLabelsFn *DetectLabelsFn) ClassifyHandler(writer http.ResponseWriter, request *http.Request) {
	detectLabelsFn.initQueryParams(request)
	log.Printf("GVisionFn.DetectLabels: q=\"%s\", o=\"%s\"", detectLabelsFn.ImageURL, detectLabelsFn.Output)

	classifiedImageData, err := detectLabelsFn.ClassifyImage()
	if err != nil {
		log.Print(err.Error())
		return
	}

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Headers", "x-requested-with")

	writer.Header().Add("Content-Type", detectLabelsFn.OutputContentType(detectLabelsFn.Output))
	fmt.Fprintf(writer, "%s\n", common.Flatten(&classifiedImageData, detectLabelsFn.Output, classifiedImageData.ToText))
}

// Private classifyImageFn

func (classifyImageFn *DetectLabelsFn) initQueryParams(request *http.Request) {
	detectLabelsFn.ImageURL = classifyImageFn.ExtractQueryStringParam(request, []string{"query", "q", "image-url", "u"}, classifyImageFn.ImageURL)
	detectLabelsFn.Output = classifyImageFn.ExtractQueryStringParam(request, []string{"o", "output"}, classifyImageFn.Output)
}

// Public ClassifyImageData

func (cIData ClassifyImageData) ToText(in interface{}) string {
	sb := bytes.NewBufferString("")

	sb.WriteString(fmt.Sprintf("image URL: %s\n", cIData.ImageURL))
	sb.WriteString("----\n")
	for _, label := range cIData.Labels {
		sb.WriteString(fmt.Sprintf("name: %s\n", label.Name))
		sb.WriteString(fmt.Sprintf("score: %f2\n", label.Score))
	}
	sb.WriteString("----\n")

	return sb.String()
}
