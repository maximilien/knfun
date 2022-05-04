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

package common

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (commonFn *CommonFn) InitCommonQueryParams(request *http.Request) {
	commonFn.SearchString = commonFn.ExtractQueryStringParam(request, []string{"q", "query", "search-string", "s"}, commonFn.SearchString)
	commonFn.Count = commonFn.ExtractQueryIntParam(request, []string{"c", "count"}, commonFn.Count)
	commonFn.Output = commonFn.ExtractQueryStringParam(request, []string{"o", "output"}, commonFn.Output)
}

func (commonFn *CommonFn) ExtractQueryStringParams(request *http.Request, paramNames []string) map[string]string {
	query := request.URL.Query()
	valueMap := map[string]string{}

	for _, paramName := range paramNames {
		if query.Get(paramName) != "" {
			valueMap[paramName] = query.Get(paramName)
		}
	}
	return valueMap
}

func (commonFn *CommonFn) ExtractQueryIntParams(request *http.Request, paramNames []string) map[string]int {
	query := request.URL.Query()
	valueMap := map[string]int{}

	for _, paramName := range paramNames {
		if query.Get(paramName) != "" {
			intValue, err := strconv.Atoi(query.Get(paramName))
			if err != nil {
				log.Print(fmt.Sprintf("`%s` query parameter value is invalid '%s'!", paramName, err.Error()))
			}
			valueMap[paramName] = intValue
		}
	}
	return valueMap
}

func (commonFn *CommonFn) ExtractQueryStringParam(request *http.Request, paramNames []string, defaultValue string) string {
	query := request.URL.Query()
	stringValue := defaultValue

	for _, paramName := range paramNames {
		if query.Get(paramName) != "" {
			stringValue = query.Get(paramName)
			break
		}
	}
	return stringValue
}

func (commonFn *CommonFn) ExtractQueryIntParam(request *http.Request, paramNames []string, defaultValue int) int {
	query := request.URL.Query()
	intValue := defaultValue

	for _, paramName := range paramNames {
		if query.Get(paramName) != "" {
			iValue, err := strconv.Atoi(query.Get(paramName))
			if err != nil {
				log.Print(fmt.Sprintf("`%s` query parameter value is invalid '%s'!", paramName, err.Error()))
			}
			intValue = iValue
			break
		}
	}
	return intValue
}

func (commonFn *CommonFn) OutputContentType(output string) string {
	switch output {
	case "yaml":
		return "application/yaml"
	case "json":
		return "application/json"
	}
	return "text/html"
}
