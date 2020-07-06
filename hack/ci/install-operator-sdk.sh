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

default_operator_sdk_version=v0.18.2

if [[ -z ${OPERATOR_SDK_VERSION} ]]; then
    OPERATOR_SDK_VERSION=$default_operator_sdk_version
fi

echo "The script will install operator-sdk ${OPERATOR_SDK_VERSION} at /usr/bin/  Please make sure the user have sudo privileges before running the script."

(which operator-sdk > /dev/null && echo "operator-sdk is already installed") || (sudo curl -L https://github.com/operator-framework/operator-sdk/releases/download/$OPERATOR_SDK_VERSION/operator-sdk-$OPERATOR_SDK_VERSION-x86_64-linux-gnu -o /usr/bin/operator-sdk && sudo chmod +x /usr/bin/operator-sdk)
#For verification
operator-sdk version