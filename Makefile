# Holds short targets for the go tool commands.
#
# Copyright 2017 Google Inc.
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
#

PKG = juju.nz/x/pipoint

VERSION = $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags "-X $(PKG).Version=$(VERSION)"

build:
	go get $(LDFLAGS) $(PKG)/pipoint

# Watch for changes, build, and push.
watch:
	watchman-make -p '**/*.go' -t push

run: build
	$(GOPATH)/bin/pipoint

check:
	go get -t $(PKG)/...
	go test $(shell go list $(PKG)/... | grep -vF /vendor)

coverage:
	go get -t $(PKG)
	go test -race -coverprofile=coverage.txt -covermode=atomic $(PKG)

push:
	GOARCH=arm GOARM=7 go get $(LDFLAGS) $(PKG)/pipoint
	rsync -zt $(GOPATH)/bin/linux_arm/pipoint pi-ed7:~

# Shh.  This is to test gitea.
