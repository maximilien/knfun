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
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/maximilien/knfun/funcs/common"

	yaml "gopkg.in/yaml.v2"
)

type ClassifiedTweet struct {
	Text             string            `json:"text"`
	ClassifiedImages []ClassifiedImage `json:"classified-images"`
}

type Tweet struct {
	Text      string   `json:"text"`
	ImageURLs []string `json:"image-urls"`
}

type ClassifiedImage struct {
	SourceURL   string       `json:"source_url"`
	ResolvedURL string       `json:"resolved_url"`
	Classifiers []Classifier `json:"classifiers"`
}

type Classifier struct {
	Name         string  `json:"name"`
	ClassifierID string  `json:"classifier_id"`
	Classes      []Class `json:"classes"`
}

type Class struct {
	Name          string  `json:"class"`
	Score         float32 `json:"score"`
	TypeHierarchy string  `json:"type_hierarchy,omitempty"`
}

type SummaryFn struct {
	common.CommonFn

	TwitterFnURL string
	WatsonFnURL  string
}

func (summaryFn *SummaryFn) Summary() ([]ClassifiedTweet, error) {
	return summaryFn.collectClassifiedTweets()
}

func (summaryFn *SummaryFn) SummaryHandler(writer http.ResponseWriter, request *http.Request) {
	//classifiedTweets, err := summaryFn.collectClassifiedTweets()
	_, err := summaryFn.collectClassifiedTweets()
	if err != nil {
		log.Fatal("Error collecting classified tweets: %s\n", err.Error())
		return
	}

	//TODO: use correct template and pass correct data
	// tmpl := template.Must(template.ParseFiles("./funcs/summary/layout.html"))
	// data := TodoPageData{
	// 	PageTitle: "My TODO list",
	// 	Todos: []Todo{
	// 		{Title: "Task 1", Done: false},
	// 		{Title: "Task 2", Done: true},
	// 		{Title: "Task 3", Done: true},
	// 	},
	// }
	// tmpl.Execute(w, data)
	//TODO
}

// Private SummaryFn

func (summaryFn *SummaryFn) searchTweets(searchString string, count int) ([]Tweet, error) {
	var err error

	twitterFnClient := http.Client{
		Timeout: time.Second * time.Duration(summaryFn.Timeout),
	}

	url := fmt.Sprintf("%s?q=%s&c=%d&o=json", summaryFn.TwitterFnURL, searchString, count)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []Tweet{}, err
	}

	req.Header.Set("User-Agent", "application/json")

	res, err := twitterFnClient.Do(req)
	if err != nil {
		return []Tweet{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Tweet{}, err
	}

	tweets := []Tweet{}
	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return []Tweet{}, err
	}

	return tweets, nil
}

func (summaryFn *SummaryFn) classifyImage(imageURL string) (ClassifiedImage, error) {
	var err error

	watsonFnClient := http.Client{
		Timeout: time.Second * time.Duration(summaryFn.Timeout),
	}

	url := fmt.Sprintf("%s?q=%s&o=json", summaryFn.WatsonFnURL, imageURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ClassifiedImage{}, err
	}

	req.Header.Set("User-Agent", "application/json")

	res, err := watsonFnClient.Do(req)
	if err != nil {
		return ClassifiedImage{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ClassifiedImage{}, err
	}

	classifiedImage := ClassifiedImage{}
	err = json.Unmarshal(body, &classifiedImage)
	if err != nil {
		return ClassifiedImage{}, err
	}

	return classifiedImage, nil
}

func (summaryFn *SummaryFn) collectTweetsWithImages(tweets []Tweet) []Tweet {
	tweetsWithImages := []Tweet{}
	for _, tweet := range tweets {
		if len(tweet.ImageURLs) > 0 {
			tweetsWithImages = append(tweetsWithImages, tweet)
		}
	}
	return tweetsWithImages
}

func (summaryFn *SummaryFn) collectClassifiedTweets() ([]ClassifiedTweet, error) {
	tweets, err := summaryFn.searchTweets(summaryFn.SearchString, summaryFn.Count)
	if err != nil {
		return []ClassifiedTweet{}, err
	}

	tweetsWithImages := summaryFn.collectTweetsWithImages(tweets)
	classifiedTweets := []ClassifiedTweet{}
	for _, tweet := range tweetsWithImages {
		classifiedImages := []ClassifiedImage{}
		for _, imageURL := range tweet.ImageURLs {
			classifiedImage, err := summaryFn.classifyImage(imageURL)
			if err != nil {
				return []ClassifiedTweet{}, err
			}
			classifiedImages = append(classifiedImages, classifiedImage)
		}
		classifiedTweet := ClassifiedTweet{
			Text:             tweet.Text,
			ClassifiedImages: classifiedImages,
		}
		classifiedTweets = append(classifiedTweets, classifiedTweet)
	}
	return classifiedTweets, nil
}

// Private ClassifiedTweet

func (cTweet ClassifiedTweet) ToText() string {
	sb := bytes.NewBufferString("")
	sb.WriteString(fmt.Sprintf("\n🐦 %s\n", cTweet.Text))
	for i, cImage := range cTweet.ClassifiedImages {
		sb.WriteString(fmt.Sprintf("\n%d.  📸 URL: `%s`\n", i, cImage.ResolvedURL))
		for _, classifier := range cImage.Classifiers {
			sb.WriteString(fmt.Sprintf("\n%d.  📚 classifier: `%s`\n", i, classifier.Name))
			for _, class := range classifier.Classes {
				sb.WriteString(fmt.Sprintf("\n.   📸 is a `%s` with `%1.3f` confidence\n", class.Name, class.Score))
			}
			sb.WriteString("------\n")
		}
	}
	return sb.String()
}

func (cTweet ClassifiedTweet) ToYAML() string {
	data, err := yaml.Marshal(&cTweet)
	if err != nil {
		panic("Error YAML marshalling ClassifiedTweet")
	}

	return string(data)
}

func (cTweet ClassifiedTweet) ToJSON() string {
	data, err := json.MarshalIndent(&cTweet, "", "  ")
	if err != nil {
		panic("Error JSON marshalling ClassifiedTweet")
	}

	return string(data)
}

func (cTweet ClassifiedTweet) Flatten(output string) string {
	outputData := ""
	switch output {
	case "yaml":
		outputData = cTweet.ToYAML()
		break
	case "json":
		outputData = cTweet.ToJSON()
		break
	default:
		outputData = cTweet.ToText()
	}
	return outputData
}
