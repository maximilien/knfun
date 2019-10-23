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

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type keys struct {
	twitterAPIKey            string
	twitterAPISecretKey      string
	twitterAccessToken       string
	twitterAccessTokenSecret string
}

type TweetData struct {
	Text      string   `yaml:"text" json:"text"`
	ImageURLs []string `yaml:"image-urls" json:"image-urls"`
}

type TweetsData []TweetData

type SearchFn struct {
	common.CommonFn

	keys keys
}

func (searchFn *SearchFn) Search() (TweetsData, error) {
	client := searchFn.createTwitterClient()
	results, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: searchFn.SearchString,
		Count: searchFn.Count,
	})
	if err != nil {
		return []TweetData{}, err
	}

	return searchFn.collectTweetsData(results.Statuses), nil
}

func (searchFn *SearchFn) SearchHandler(writer http.ResponseWriter, request *http.Request) {
	searchString := searchFn.ExtractQueryStringParam(request, []string{"q", "query", "search-string", "s"}, searchFn.SearchString)
	output := searchFn.ExtractQueryStringParam(request, []string{"o", "output"}, searchFn.Output)
	count := searchFn.ExtractQueryIntParam(request, []string{"c", "count"}, searchFn.Count)

	if searchString == "" {
		log.Printf("Must pass a query string using 'q' or 'query' parameter")
		return
	}

	client := searchFn.createTwitterClient()
	results, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: searchString,
		Count: count,
	})

	if err != nil {
		log.Printf("Error searching for tweets: %s\n", err.Error())
		return
	}

	tweetsData := searchFn.collectTweetsData(results.Statuses)
	writer.Header().Add("Content-Type", searchFn.OutputContentType(output))
	fmt.Fprintf(writer, "%s\n", tweetsData.Flatten(output))
}

// Private SearchFn

func (searchFn *SearchFn) createTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(searchFn.keys.twitterAPIKey, searchFn.keys.twitterAPISecretKey)
	token := oauth1.NewToken(searchFn.keys.twitterAccessToken, searchFn.keys.twitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

func (searchFn *SearchFn) collectTweetsData(tweets []twitter.Tweet) TweetsData {
	tweetsData := TweetsData{}
	for _, tweet := range tweets {
		tweetData := TweetData{Text: tweet.Text}
		imageURLs := []string{}
		for _, media := range tweet.Entities.Media {
			if media.MediaURL != "photo" {
				imageURLs = append(imageURLs, media.MediaURL)
			}
		}
		tweetData.ImageURLs = imageURLs
		tweetsData = append(tweetsData, tweetData)
	}
	return tweetsData
}

// Private TweetsData

func (tweets TweetsData) ToText() string {
	sb := bytes.NewBufferString("")
	for _, tweetData := range tweets {
		sb.WriteString(fmt.Sprintf("\n🐦 %s\n", tweetData.Text))
		if len(tweetData.ImageURLs) > 0 {
			for _, imageUrl := range tweetData.ImageURLs {
				sb.WriteString(fmt.Sprintf("- 📸 %s\n", imageUrl))
			}
		}
		sb.WriteString("------\n")
	}
	return sb.String()
}

func (tweets TweetsData) ToYAML() string {
	data, err := yaml.Marshal(&tweets)
	if err != nil {
		panic("Error YAML marshalling TweetData")
	}

	return string(data)
}

func (tweets TweetsData) ToJSON() string {
	data, err := json.MarshalIndent(&tweets, "", "  ")
	if err != nil {
		panic("Error JSON marshalling TweetData")
	}

	return string(data)
}

func (tweets TweetsData) Flatten(output string) string {
	outputData := ""
	switch output {
	case "yaml":
		outputData = tweets.ToYAML()
		break
	case "json":
		outputData = tweets.ToJSON()
		break
	default:
		outputData = tweets.ToText()
	}
	return outputData
}
