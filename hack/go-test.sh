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


if [[ -z ${ENVTEST_ASSETS_DIR} ]]; then
  ENVTEST_ASSETS_DIR=testbin
fi

mkdir -p $(pwd)/${ENVTEST_ASSETS_DIR}
test -f $(pwd)/${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo testbin/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.6.3/hack/setup-envtest.sh
sed -i "s,#\!.*,#\!\/bin\/bash,g" $(pwd)/${ENVTEST_ASSETS_DIR}/setup-envtest.sh
source $(pwd)/${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(pwd)/${ENVTEST_ASSETS_DIR}; setup_envtest_env $(pwd)/${ENVTEST_ASSETS_DIR}; \
go test ./cmd/... -p=1 -count=1 -coverprofile cmd-cover.out; \
go test ./pkg/... -p=1 -count=1 -coverprofile pkg-cover.out; \
go test ./controllers/... -p=1 -count=1 -coverprofile controllers-cover.out