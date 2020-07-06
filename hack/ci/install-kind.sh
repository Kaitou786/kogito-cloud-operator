#!/bin/bash
# Copyright 2020 Red Hat, Inc. and/or its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


default_kind_version=v0.8.1

if [[ -z ${KIND_VERSION} ]]; then
    KIND_VERSION=$default_kind_version
fi

echo "The script will install kind ${KIND_VERSION} at /usr/bin/  Please make sure the user have sudo privileges before running the script."

( which kind > /dev/null && echo "kind is already installed") || (sudo curl -L https://kind.sigs.k8s.io/dl/$KIND_VERSION/kind-$(uname)-amd64 -o /usr/bin/kind && sudo chmod +x /usr/bin/kind)

#for verification
kind version