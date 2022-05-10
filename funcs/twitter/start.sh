#!/bin/bash

# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

/usr/local/twitter-fn search NBA -o json -c 10 -p 8080 -S \
                --twitter-api-key $TWITTER_API_KEY \
                --twitter-api-secret-key $TWITTER_API_SECRET_KEY \
                --twitter-access-token $TWITTER_ACCESS_TOKEN \
                --twitter-access-token-secret $TWITTER_ACCESS_TOKEN_SECRET