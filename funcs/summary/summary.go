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
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
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
	ImageURL string  `json:"ImageURL"`
	Labels   []Label `json:"labels"`
}

type Label struct {
	Name  string  `json:"name"`
	Score float32 `json:"score"`
}

type SummaryFn struct {
	common.CommonFn

	TwitterFnURL string
	WatsonFnURL  string
}

type SummaryPageData struct {
	PageTitle string

	Tweets           []Tweet
	ClassifiedTweets []ClassifiedTweet

	WatsonFnURL string
	Timeout     int
}

func (summaryFn *SummaryFn) Summary() ([]ClassifiedTweet, error) {
	return summaryFn.collectClassifiedTweets()
}

func (summaryFn *SummaryFn) SummaryHandler(writer http.ResponseWriter, request *http.Request) {
	summaryFn.InitCommonQueryParams(request)
	log.Printf("SummaryFn.Summary: s=\"%s\", c=\"%d\", o=\"%s\"", summaryFn.SearchString, summaryFn.Count, summaryFn.Output)

	classifiedTweets, err := summaryFn.Summary()
	if err != nil {
		log.Printf("Error collecting classified tweets: %s\n", err.Error())
		return
	}

	tmpl := template.Must(template.ParseFiles("./funcs/summary/layout.html"))
	data := SummaryPageData{
		PageTitle:        fmt.Sprintf("Recent tweets with images for search `%s`", summaryFn.SearchString),
		ClassifiedTweets: classifiedTweets,
	}

	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Printf("Error executing template with classified tweets: %s\n", err.Error())
		return
	}
}

func (summaryFn *SummaryFn) SummaryAsyncHandler(writer http.ResponseWriter, request *http.Request) {
	summaryFn.InitCommonQueryParams(request)
	log.Printf("SummaryFn.Summary: s=\"%s\", c=\"%d\", o=\"%s\"", summaryFn.SearchString, summaryFn.Count, summaryFn.Output)

	tweets, err := summaryFn.searchTweets(summaryFn.SearchString, summaryFn.Count)
	if err != nil {
		log.Printf("Error collecting tweets: %s\n", err.Error())
		return
	}

	tmplName := path.Base("./funcs/summary/async_layout.html")
	tmpl := template.New(tmplName)
	tmpl.Funcs(template.FuncMap{
		"ClassifyImage": func(watsonFnURL string, imageURL string, timeout int) (ClassifiedImage, error) {
			return classifyImage(watsonFnURL, imageURL, timeout)
		},
	})

	tmpl, err = tmpl.ParseFiles("./funcs/summary/async_layout.html")
	if err != nil {
		log.Printf("Error parsing Golang template for tweets: %s\n", err.Error())
		return
	}

	data := SummaryPageData{
		PageTitle: fmt.Sprintf("Recent tweets for search `%s`", summaryFn.SearchString),
		Tweets:    tweets,

		WatsonFnURL: summaryFn.WatsonFnURL,
		Timeout:     summaryFn.Timeout,
	}

	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Printf("Error executing template with tweets: %s\n", err.Error())
		return
	}
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
			classifiedImage, err := classifyImage(summaryFn.WatsonFnURL, imageURL, summaryFn.Timeout)
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

// Private function

func classifyImage(watsonFnURL string, imageURL string, timeout int) (ClassifiedImage, error) {
	var err error

	watsonFnClient := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	url := fmt.Sprintf("%s?q=%s&o=json", watsonFnURL, imageURL)
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

	fmt.Println(string(body[:])) //DEBUG

	classifiedImage := ClassifiedImage{}
	err = json.Unmarshal(body, &classifiedImage)
	if err != nil {
		return ClassifiedImage{}, err
	}

	return classifiedImage, nil
}

// ClassifiedTweet

func (cTweet ClassifiedTweet) ToText() string {
	sb := bytes.NewBufferString("")
	sb.WriteString(fmt.Sprintf("\n🐦 %s\n", cTweet.Text))
	for i, cImage := range cTweet.ClassifiedImages {
		sb.WriteString(fmt.Sprintf("\n%d.  📸 URL: `%s`\n", i, cImage.ImageURL))
		for _, label := range cImage.Labels {
			sb.WriteString(fmt.Sprintf("\n.   📸 is a `%s` with `%1.3f` confidence\n", label.Name, label.Score))
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
