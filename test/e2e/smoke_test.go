// Copyright 2019 The Knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build e2e

package e2e

import (
	"testing"

	"gotest.tools/assert"
	"knative.dev/client/pkg/util"
)

func TestSmoke(t *testing.T) {
	t.Parallel()
	test := NewE2eTest(t, false)
	test.Setup(t)
	defer test.Teardown(t)

	t.Run("verifies twitter-fn search", func(t *testing.T) {
		test.twitterFn_Search_Local(t)
	})

	t.Run("verifies watson-fn vr classify", func(t *testing.T) {
		test.watsonFn_VR_Classify_Local(t)
	})
}

func (test *e2eTest) twitterFn_Search_Local(t *testing.T) {
	out, err := test.funcs.Run("twitter-fn", []string{"search", "NBA", "-c", "10", "-o", "json"})
	assert.NilError(t, err)

	assert.Check(t, util.ContainsAll(out, "Using config file:"))
}

func (test *e2eTest) watsonFn_VR_Classify_Local(t *testing.T) {
	out, err := test.funcs.Run("watson-fn", []string{"vr", "classify", "http://pbs.twimg.com/media/EHb34-KXYAESI46.jpg", "-o", "json"})
	assert.NilError(t, err)

	assert.Check(t, util.ContainsAll(out, "Using config file:"))
}
