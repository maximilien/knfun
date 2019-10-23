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

	yaml "gopkg.in/yaml.v2"
)

type Data struct{}

func (data Data) ToText() string {
	sb := bytes.NewBufferString("")
	sb.WriteString(fmt.Sprintf("%#v\n", data))
	return sb.String()
}

func (data Data) ToYAML() string {
	yData, err := yaml.Marshal(&data)
	if err != nil {
		panic("Error YAML marshalling Data")
	}
	return string(yData)
}

func (data Data) ToJSON() string {
	jData, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		panic("Error JSON marshalling Data")
	}
	return string(jData)
}

func (data Data) Flatten(output string) string {
	outputData := ""
	switch output {
	case "yaml":
		outputData = data.ToYAML()
		break
	case "json":
		outputData = data.ToJSON()
		break
	default:
		outputData = data.ToText()
	}
	return outputData
}
