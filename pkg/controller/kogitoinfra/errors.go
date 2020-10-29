// Copyright 2020 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kogitoinfra

import (
	"fmt"
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
)

// reconciliationError type for KogitoInfra reconciliation cycle cases.
// MUST be public because it's used by the sub packages of kogitoinfra.
type reconciliationError struct {
	Reason     v1alpha1.KogitoInfraConditionReason
	innerError error
}

// String stringer implementation
func (e reconciliationError) String() string {
	return e.innerError.Error()
}

// Error error implementation
func (e reconciliationError) Error() string {
	return e.innerError.Error()
}

// newResourceNotFoundError ...
func newResourceNotFoundError(kind, instance, namespace string) reconciliationError {
	return reconciliationError{
		Reason:     v1alpha1.ResourceNotFound,
		innerError: fmt.Errorf("%s instance(%s) not found in namespace %s", kind, instance, namespace),
	}
}

// newResourceAPINotFoundError ...
func newResourceAPINotFoundError(resource *v1alpha1.Resource) reconciliationError {
	return reconciliationError{
		Reason:     v1alpha1.ResourceAPINotFound,
		innerError: fmt.Errorf("%s CRD is not available in the cluster, this feature is not available. Please install the required Operator first. ", resource.APIVersion),
	}
}

// newUnsupportedAPIError ...
func newUnsupportedAPIError(instance *v1alpha1.KogitoInfra) reconciliationError {
	return reconciliationError{
		Reason: v1alpha1.UnsupportedAPIKind,
		innerError: fmt.Errorf("API %s is not supported for kind %s. Supported APIs are: %v",
			instance.Spec.Resource.APIVersion,
			instance.Spec.Resource.Kind,
			getSupportedResources()),
	}
}

func newResourceNotReadyError(instance *v1alpha1.KogitoInfra, err error) reconciliationError {
	return reconciliationError{
		Reason:     v1alpha1.ResourceNotReady,
		innerError: err,
	}
}

func getSupportedResources() []string {
	res := getSupportedInfraResources()
	keys := make([]string, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	return keys
}

func reasonForError(err error) v1alpha1.KogitoInfraConditionReason {
	switch t := err.(type) {
	case reconciliationError:
		return t.Reason
	}
	return v1alpha1.ReconciliationFailure
}