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

set -o pipefail

source_dirs="funcs"

username=${USERNAME:-}
cr_url=${CR_URL:-}
cr_pat=${CR_PAT:-}

knfun_config=$HOME/.knfun.yaml

# Store for later
if [ -z "$1" ]; then
    ARGS=("")
else
    ARGS=("$@")
fi

set -eu

# Temporary fix for iTerm issue https://gitlab.com/gnachman/iterm2/issues/7901
S=""
if [ -n "${ITERM_PROFILE:-}" ]; then
  S=" "
fi
# Run build
run() {
  # Switch on modules unconditionally
  export GO111MODULE=on

  # Jump into project directory
  pushd $(basedir) >/dev/null 2>&1

  # Print help if requested
  if $(has_flag --help -h); then
    display_help
    exit 0
  fi

  if $(has_flag --watch -w); then
    # Build and test first
    go_build

    if $(has_flag --test -t); then
       go_test
    fi

    # Go in endless loop, to be stopped with CTRL-C
    watch
  fi

  # Fast mode: Only compile and maybe run test
  if $(has_flag --fast -f); then
    go_build

    if $(has_flag --test -t); then
       go_test
    fi
    exit 0
  fi

  # Run only tests
  if $(has_flag --test -t); then
    go_test
    exit 0
  fi

  # Run only codegen
  if $(has_flag --codegen -c); then
    codegen
    exit 0
  fi

  # Run build images and push
  if $(has_flag --docker -d); then
    login_cr
    build_images
    push_images
    exit 0
  fi

  # Run only build images
  if $(has_flag --images -i); then
    login_cr
    build_images
    exit 0
  fi

  # Run only push images
  if $(has_flag --push -p); then
    login_cr
    push_images
    exit 0
  fi

  # Run only scan images
  if $(has_flag --scan -s); then
    login_cr
    scan_images
    exit 0
  fi

  # Default flow
  codegen
  go_build

  echo "────────────────────────────────────────────"
  echo "success"
}

codegen() {
  # Update dependencies
  update_deps

  # Format source code and cleanup imports
  source_format

  # Check for license headers
  check_license

  # Auto generate cli docs
  generate_docs
}

checks() {
  check_username  
  check_cr_url
  check_cr_pat
}

check_username() {
  if [[ -z "${username}" ]]; then
    echo "Please set environment variable USERNAME with your container registry username, e.g., Docker or GH"
    exit 1
  fi
}

check_cr_url() {
  if [[ -z "${cr_url}" ]]; then
    echo "Please set environment variable CR_URL with your container registry URL of choice: docker.io or ghcr.io"
    exit 1
  fi
}

check_cr_pat() {
  if [[ -z "${cr_pat}" ]]; then
    echo "Please set environment variable CR_PAT with your container registry of choice personal access token"
    exit 1
  fi
}

login_cr() {
  checks
  echo ${cr_pat} | docker login ${cr_url} -u $username --password-stdin
}

build_images() {
  echo "🚧 🐳 build images"

  echo "   🚧 🐳 twitter-fn"
  docker build --platform linux/amd64 -f ./funcs/twitter/Dockerfile -t ${cr_url}/${username}/twitter-fn .

  echo "   🚧 🐳 watson-fn"
  docker build --platform linux/amd64 -f ./funcs/watson/Dockerfile -t ${cr_url}/${username}/watson-fn .

  echo "   🚧 🐳 gvision-fn"
  docker build --platform linux/amd64 -f ./funcs/gvision/Dockerfile -t ${cr_url}/${username}/gvision-fn .

  echo "   🚧 🐳 summary-fn"
  docker build --platform linux/amd64 -f ./funcs/summary/Dockerfile -t ${cr_url}/${username}/summary-fn .
}

push_images() {
  echo "📤 🐳 push images"

  echo "   📤 🐳 twitter-fn"
  docker push ${cr_url}/${username}/twitter-fn

  echo "   📤 🐳 watson-fn"
  docker push ${cr_url}/${username}/watson-fn

  echo "   📤 🐳 gvision-fn"
  docker push ${cr_url}/${username}/gvision-fn

  echo "   📤 🐳 summary-fn"
  docker push ${cr_url}/${username}/summary-fn
}

scan_images() {
  echo "🔒 🐳 scan images"

  echo "   🔒 🐳 twitter-fn"
  docker scan ${cr_url}/${username}/twitter-fn

  echo "   🔒 🐳 watson-fn"
  docker scan ${cr_url}/${username}/watson-fn

  echo "   🔒 🐳 gvision-fn"
  docker scan ${cr_url}/${username}/gvision-fn

  echo "   🔒 🐳 summary-fn"
  docker scan ${cr_url}/${username}/summary-fn
}

go_fmt() {
  echo "🧹 ${S}Format"
  find $(echo $source_dirs) -name "*.go" -print0 | xargs -0 gofmt -s -w
}

source_format() {
  set +e
  which goimports >/dev/null 2>&1
  if [ $? -ne 0 ]; then
     echo "✋ No 'goimports' found. Please use"
     echo "✋   go install golang.org/x/tools/cmd/goimports"
     echo "✋ to enable import cleanup. Import cleanup skipped."

     # Run go fmt instead
     go_fmt
  else
     echo "🧽 ${S}Format"
     goimports -w $(echo $source_dirs)
     find $(echo $source_dirs) -name "*.go" -print0 | xargs -0 gofmt -s -w
  fi
  set -e
}

go_build() {
  echo "🚧 Compile"
  go build -mod=vendor -ldflags "$(build_flags $(basedir))" -o twitter-fn ./funcs/twitter/...
  go build -mod=vendor -ldflags "$(build_flags $(basedir))" -o watson-fn ./funcs/watson/...
  go build -mod=vendor -ldflags "$(build_flags $(basedir))" -o gvision-fn ./funcs/gvision/...
  go build -mod=vendor -ldflags "$(build_flags $(basedir))" -o summary-fn ./funcs/summary/...
}

go_test() {
  export PATH=$PWD:$PATH
  local basedir=$(basedir)
  if [ ! -f $knfun_config ]; then
    echo "Please create a ~/.knfun.yaml file with all funcs credentials"
    echo "  see: https://github.com/maximilien/knfun#credentials"
    exit 1
  fi

  local test_output=$(mktemp /tmp/knfun-test-output.XXXXXX)

  local red=""
  local reset=""
  # Use color only when a terminal is set
  if [ -t 1 ]; then
    red="[31m"
    reset="[39m"
  fi

  echo "🧪 ${S}Tests"
  set +e
  echo "  🧪 e2e"  
  go test ${basedir}/test/e2e/ -run TestSmoke -test.v --tags 'e2e' "$@" >$test_output 2>&1
  local err=$?
  if [ $err -ne 0 ]; then
    echo "🔥 ${red}Failure${reset}"
    cat $test_output | sed -e "s/^.*\(FAIL.*\)$/$red\1$reset/"
    rm $test_output
    exit $err
  fi
  rm $test_output
}

check_license() {
  echo "⚖️ ${S}License"
  local required_keywords=("Authors" "Apache License" "LICENSE-2.0")
  local extensions_to_check=("sh" "go" "yaml" "yml" "json")

  local check_output=$(mktemp /tmp/kn-client-licence-check.XXXXXX)
  for ext in "${extensions_to_check[@]}"; do
    find . -name "*.$ext" -a \! -path "./vendor/*" -a \! -path "./.*" -print0 |
      while IFS= read -r -d '' path; do
        for rword in "${required_keywords[@]}"; do
          if ! grep -q "$rword" "$path"; then
            echo "   $path" >> $check_output
          fi
        done
      done
  done
  if [ -s $check_output ]; then
    echo "🔥 No license header found in:"
    cat $check_output | sort | uniq
    echo "🔥 Please fix and retry."
    rm $check_output
    exit 1
  fi
  rm $check_output
}


update_deps() {
  echo "🕸️ ${S}Update"
  go mod vendor
}

generate_docs() {
  echo "📖 Docs"
  # rm -rf "./docs/cmd"
  # mkdir -p "./docs/cmd"
  # go run "./hack/generate-docs.go" "."
}

watch() {
    local command="./hack/build.sh --fast"
    local fswatch_opts='-e "^\..*$" -o pkg cmd'
    if $(has_flag --test -t); then
      command="$command --test"
    fi
    if $(has_flag --verbose); then
      fswatch_opts="$fswatch_opts -v"
    fi
    set +e
    which fswatch >/dev/null 2>&1
    if [ $? -ne 0 ]; then
      local green="[32m"
      local reset="[39m"

      echo "🤷 Watch: Cannot find ${green}fswatch${reset}"
      echo "🌏 Please see ${green}http://emcrisostomo.github.io/fswatch/${reset} for installation instructions"
      exit 1
    fi
    set -e

    echo "🔁 Watch"
    fswatch $fswatch_opts | xargs -n1 -I{} sh -c "$command && echo 👌 OK"
}

# Dir where this script is located
basedir() {
    # Default is current directory
    local script=${BASH_SOURCE[0]}

    # Resolve symbolic links
    if [ -L $script ]; then
        if readlink -f $script >/dev/null 2>&1; then
            script=$(readlink -f $script)
        elif readlink $script >/dev/null 2>&1; then
            script=$(readlink $script)
        elif realpath $script >/dev/null 2>&1; then
            script=$(realpath $script)
        else
            echo "ERROR: Cannot resolve symbolic link $script"
            exit 1
        fi
    fi

    local dir=$(dirname "$script")
    local full_dir=$(cd "${dir}/.." && pwd)
    echo ${full_dir}
}

# Checks if a flag is present in the arguments.
has_flag() {
    filters="$@"
    for var in "${ARGS[@]}"; do
        for filter in $filters; do
          if [ "$var" = "$filter" ]; then
              echo 'true'
              return
          fi
        done
    done
    echo 'false'
}

cross_build() {
  local basedir=$(basedir)
  local ld_flags="$(build_flags $basedir)"
  local pkg="github.com/knative/client/pkg/kn/commands"
  local failed=0

  echo "⚔️ ${S}Compile"

  export CGO_ENABLED=0
  echo "   🐧 kn-linux-amd64"
  GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "${ld_flags}" -o ./kn-linux-amd64 ./cmd/... || failed=1
  echo "   🍏 kn-darwin-amd64"
  GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags "${ld_flags}" -o ./kn-darwin-amd64 ./cmd/... || failed=1
  echo "   🎠 kn-windows-amd64.exe"
  GOOS=windows GOARCH=amd64 go build -mod=vendor -ldflags "${ld_flags}" -o ./kn-windows-amd64.exe ./cmd/... || failed=1

  return ${failed}
}

# Display a help message.
display_help() {
    local command="${1:-}"
    cat <<EOT
Knative client build script

Usage: $(basename $BASH_SOURCE) [... options ...]

with the following options:

-f  --fast                    Only compile (without dep update, formatting, testing, doc gen)
-t  --test                    Run tests when used with --fast or --watch
-c  --codegen                 Runs formatting, doc gen and update without compiling/testing
-d  --docker-images           Generates Docker images for each funcs (twitter-fn, watson-fn, summary-fn)
-w  --watch                   Watch for source changes and recompile in fast mode
-x  --all                     Build binaries for all platforms
-h  --help                    Display this help message
    --verbose                 More output
    --debug                   Debug information for this script (set -x)

You can add a symbolic link to this build script into your PATH so that it can be
called from everywhere. E.g.:

ln -s $(basedir)/hack/build.sh /usr/local/bin/kn_build.sh

Examples:

* Update deps, format, license check,
  gen docs, compile, test: ........... build.sh
* Compile only: ...................... build.sh --fast
* Run only tests: .................... build.sh --test
* Compile with tests: ................ build.sh -f -t
* Automatic recompilation: ........... build.sh --watch
* Build and push Docker images: ...... build.sh --docker -d
* Build Docker images: ............... build.sh --images -i
* Push Docker images: ................ build.sh --push -p
* Scan Docker images: ................ build.sh --scan -s
* Build cross platform binaries: ..... build.sh --all
EOT
}

if $(has_flag --debug); then
    export PS4='+($(basename ${BASH_SOURCE[0]}):${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
    set -x
fi

# Shared funcs with CI
source $(basedir)/hack/build-flags.sh

run $*
