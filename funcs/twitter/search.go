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
	searchFn.InitCommonQueryParams(request)
	log.Printf("TwitterFn.Search: q=\"%s\", c=\"%d\", o=\"%s\"", searchFn.SearchString, searchFn.Count, searchFn.Output)

	tweetsData, err := searchFn.Search()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	writer.Header().Add("Content-Type", searchFn.OutputContentType(searchFn.Output))
	fmt.Fprintf(writer, "%s\n", common.Flatten(&tweetsData, searchFn.Output, tweetsData.ToText))
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

func (tweet TweetData) ToText(in interface{}) string {
	sb := bytes.NewBufferString("")
	sb.WriteString(fmt.Sprintf("\n🐦 %s\n", tweet.Text))
	if len(tweet.ImageURLs) > 0 {
		for _, imageUrl := range tweet.ImageURLs {
			sb.WriteString(fmt.Sprintf("- 📸 %s\n", imageUrl))
		}
	}
	return sb.String()
}

func (tweets TweetsData) ToText(in interface{}) string {
	sb := bytes.NewBufferString("")
	for _, tweetData := range tweets {
		sb.WriteString(tweetData.ToText(tweetData))
		sb.WriteString("------\n")
	}
	return sb.String()
}
