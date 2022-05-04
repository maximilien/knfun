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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ToTextFunc = func(in interface{}) string

func ToYAML(in interface{}) string {
	yData, err := yaml.Marshal(in)
	if err != nil {
		panic("error YAML marshalling Data")
	}
	return string(yData)
}

func ToJSON(in interface{}) string {
	jData, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		panic("error JSON marshalling Data")
	}
	return string(jData)
}

func ToText(in interface{}) string {
	sb := bytes.NewBufferString("")
	sb.WriteString(fmt.Sprintf("%#v\n", in))
	return sb.String()
}

func Flatten(in interface{}, output string, toText ToTextFunc) string {
	outputData := ""
	switch output {
	case "yaml":
		outputData = ToYAML(in)
		break
	case "json":
		outputData = ToJSON(in)
		break
	default:
		outputData = toText(in)
	}
	return outputData
}

func DownloadTmpFile(url string) (string, error) {
	file, err := ioutil.TempFile("", "knfun")
	if err != nil {
		return "", err
	}

	err = DownloadFile(file.Name(), url)
	return file.Name(), err
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
