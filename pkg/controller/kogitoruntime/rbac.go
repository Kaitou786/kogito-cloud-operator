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

package kogitoruntime

import (
	"github.com/kiegroup/kogito-cloud-operator/pkg/client"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	serviceAccountName = "kogito-service-viewer"
	roleName           = "kogito-service-viewer"
	roleBindingName    = "kogito-service-viewer"
	roleAPIGroup       = "rbac.authorization.k8s.io"
)

func createServiceAccountIfNotExists(client *client.Client, namespace string) (err error) {
	if err = kubernetes.ResourceC(client).CreateIfNotExists(&v1.ServiceAccount{
		ObjectMeta: v12.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
	}); err != nil {
		return
	}
	return nil
}

func createRoleIfNotExists(client *client.Client, namespace string) (err error) {
	if err = kubernetes.ResourceC(client).CreateIfNotExists(&rbac.Role{
		ObjectMeta: v12.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
		},
		Rules: []rbac.PolicyRule{
			{
				Verbs:     []string{"list", "get", "watch", "update", "patch"},
				APIGroups: []string{""},
				Resources: []string{"services", "configmaps"},
			},
		},
	}); err != nil {
		return
	}
	return nil
}

func createRoleBindingIfNotExists(client *client.Client, namespace string) (err error) {
	if err = kubernetes.ResourceC(client).CreateIfNotExists(&rbac.RoleBinding{
		ObjectMeta: v12.ObjectMeta{
			Name:      roleBindingName,
			Namespace: namespace,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "ServiceAccount",
				Name: serviceAccountName,
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: roleAPIGroup,
			Name:     roleName,
			Kind:     "Role",
		},
	}); err != nil {
		return
	}
	return nil
}
