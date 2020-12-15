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

/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"github.com/kiegroup/kogito-cloud-operator/pkg/client"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	"github.com/kiegroup/kogito-cloud-operator/pkg/framework"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure/services"
	"github.com/kiegroup/kogito-cloud-operator/pkg/logger"
	imagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kiegroup/kogito-cloud-operator/api/v1beta1"
)

// KogitoRuntimeReconciler reconciles a KogitoRuntime object
type KogitoRuntimeReconciler struct {
	*client.Client
	Log    logger.Logger
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KogitoRuntime object and makes changes based on the state read
// and what is in the KogitoRuntime.Spec
// +kubebuilder:rbac:groups=app.kiegroup.org,resources=kogitoruntimes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.kiegroup.org,resources=kogitoruntimes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets;deployments;replicasets,verbs=*
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create;list;delete
// +kubebuilder:rbac:groups=infinispan.org,resources=infinispans,verbs=get;create;list;delete;watch
// +kubebuilder:rbac:groups=kafka.strimzi.io,resources=kafkas;kafkatopics,verbs=get;create;list;delete;watch
// +kubebuilder:rbac:groups=keycloak.org,resources=keycloaks,verbs=get;create;list;delete;watch
// +kubebuilder:rbac:groups=apps,resourceNames=kogito-operator,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=eventing.knative.dev,resources=brokers,verbs=get;list;watch
// +kubebuilder:rbac:groups=eventing.knative.dev,resources=triggers,verbs=get;list;watch;create;delete;update
// +kubebuilder:rbac:groups=sources.knative.dev,resources=sinkbindings,verbs=get;list;watch;create;delete;update
// +kubebuilder:rbac:groups=integreatly.org,resources=grafanadashboards,verbs=get;create;list;watch;create;delete;update
func (r *KogitoRuntimeReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	r.Log.Info("Reconciling for", "KogitoRuntime", req.Name, "Namespace", req.Namespace)

	instance, err := infrastructure.FetchKogitoRuntimeService(r.Client, req.Name, req.Namespace)
	if err != nil {
		return
	}
	if instance == nil {
		r.Log.Debug("Instance not found", "KogitoRuntime", req.Name, "Namespace", req.Namespace)
		return
	}

	if err = r.setupRBAC(req.Namespace); err != nil {
		return
	}

	if err = infrastructure.MountProtoBufConfigMapOnDataIndex(r.Client, instance); err != nil {
		r.Log.Error(err, "Fail to mount Proto Buf config map of Kogito runtime on DataIndex", "Instance", instance.Name)
		return
	}

	definition := services.ServiceDefinition{
		Request:            req,
		DefaultImageTag:    infrastructure.LatestTag,
		SingleReplica:      false,
		OnDeploymentCreate: onDeploymentCreate,
		OnObjectsCreate:    r.onObjectsCreate,
		OnGetComparators:   onGetComparators,
		CustomService:      true,
	}
	requeueAfter, err := services.NewServiceDeployer(definition, instance, r.Client, r.Scheme).Deploy()
	if err != nil {
		return
	}
	if requeueAfter > 0 {
		r.Log.Info("Waiting for all resources to be created, scheduling for 30 seconds from now")
		result.RequeueAfter = requeueAfter
		result.Requeue = true
	}
	return
}

// SetupWithManager registers the controller with manager
func (r *KogitoRuntimeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log.Debug("Adding watched objects for KogitoRuntime controller")
	// Create a new controller
	c, err := controller.New("kogitoruntime-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	pred := predicate.Funcs{
		// Don't watch delete events as the resource removals will be handled by its finalizer
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.MetaNew.GetDeletionTimestamp().IsZero()
		},
	}
	err = c.Watch(&source.Kind{Type: &v1beta1.KogitoRuntime{}}, &handler.EnqueueRequestForObject{}, pred)
	if err != nil {
		return err
	}

	watchedObjects := []framework.WatchedObjects{
		{
			GroupVersion: routev1.GroupVersion,
			AddToScheme:  routev1.Install,
			Objects:      []runtime.Object{&routev1.Route{}},
		},
		{
			GroupVersion: imagev1.GroupVersion,
			AddToScheme:  imagev1.Install,
			Objects:      []runtime.Object{&imagev1.ImageStream{}},
		},
		{
			Objects: []runtime.Object{&corev1.Service{}, &appsv1.Deployment{}, &corev1.ConfigMap{}},
		},
	}
	controllerWatcher := framework.NewControllerWatcher(r.Client, mgr, c, &v1beta1.KogitoRuntime{})
	if err = controllerWatcher.Watch(watchedObjects...); err != nil {
		return err
	}
	return nil
}

func (r *KogitoRuntimeReconciler) setupRBAC(namespace string) (err error) {
	// create service viewer role
	if err = kubernetes.ResourceC(r.Client).CreateIfNotExists(getServiceViewerRole(namespace)); err != nil {
		r.Log.Error(err, "Fail to create role for service viewer")
		return
	}

	// create service viewer service account
	if err = kubernetes.ResourceC(r.Client).CreateIfNotExists(getServiceViewerServiceAccount(namespace)); err != nil {
		r.Log.Error(err, "Fail to create service account for service viewer")
		return
	}

	// create service viewer rolebinding
	if err = kubernetes.ResourceC(r.Client).CreateIfNotExists(getServiceViewerRoleBinding(namespace)); err != nil {
		r.Log.Error(err, "Fail to create role binding for service viewer")
		return
	}
	return
}
