// Copyright 2019 The knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"testing"
)

type funcs struct {
	t *testing.T
	l Logger
}

// Run the func with name with args
func (f funcs) Run(name string, args []string) (string, error) {
	out, err := f.RunWithOpts(name, args, runOpts{})
	return out, err
}

// Run the the func with name with args
func (f funcs) RunWithOpts(name string, args []string, opts runOpts) (string, error) {
	return runCLIWithOpts(name, args, opts, f.l)
}
