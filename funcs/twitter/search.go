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
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type twitterKeys struct {
	apiKey            string
	apiSecretKey      string
	accessToken       string
	accessTokenSecret string
}

type TweetData struct {
	Text      string   `yaml:"text" json:"text"`
	ImageURLs []string `yaml:"image-urls" json:"image-urls"`
}

type TweetsData []TweetData

type SearchFn struct {
	SearchString string
	Count        int
	Output       string

	StartServer bool
	Port        int

	keys twitterKeys
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
	searchString, count, output := searchFn.extractQueryParams(request)

	if searchString == "" {
		log.Fatal("Must pass a query string using 'q' or 'query' parameter")
		return
	}

	client := searchFn.createTwitterClient()
	results, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: searchString,
		Count: count,
	})

	if err != nil {
		log.Fatal("Error searching for tweets: %s\n", err.Error())
		return
	}

	tweetsData := searchFn.collectTweetsData(results.Statuses)
	writer.Header().Add("Content-Type", searchFn.outputContentType(output))
	fmt.Fprintf(writer, "%s\n", tweetsData.Flatten(output))
}

// Private SearchFn

func (searchFn *SearchFn) createTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(searchFn.keys.apiKey, searchFn.keys.apiSecretKey)
	token := oauth1.NewToken(searchFn.keys.accessToken, searchFn.keys.accessTokenSecret)
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

func (searchFn *SearchFn) extractQueryParams(request *http.Request) (string, int, string) {
	var err error
	searchString, count, output := searchFn.SearchString, searchFn.Count, searchFn.Output

	query := request.URL.Query()

	if query.Get("q") != "" {
		searchString = query.Get("q")
	} else {
		if query.Get("query") != "" {
			searchString = query.Get("query")
		}
	}

	if query.Get("c") != "" {
		count, err = strconv.Atoi(query.Get("c"))
		if err != nil {
			log.Fatal(fmt.Sprintf("Count query parameter 'c' is invalid '%s'!", err.Error()))
		}
	} else {
		if query.Get("count") != "" {
			count, err = strconv.Atoi(query.Get("count"))
			if err != nil {
				log.Fatal("Count query parameter 'count' is invalid!")
			}
		}
	}

	if query.Get("o") != "" {
		output = query.Get("o")
	} else {
		if query.Get("ouput") != "" {
			output = query.Get("output")
		}
	}

	return searchString, count, output
}

func (searchFn *SearchFn) outputContentType(output string) string {
	switch output {
	case "yaml":
		return "application/yaml"
	case "json":
		return "application/json"
	}
	return "text/html"
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
