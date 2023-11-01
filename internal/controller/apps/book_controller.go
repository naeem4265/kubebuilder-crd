/*
Copyright 2023.

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

package apps

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdappsv1 "github.com/naeem4265/kubebuilder-crd/api/apps/v1"
)

var (
	depOwnerKey = ".metadata.controller"
	svcOwnerKey = ".metadata.controller"
	apiGVStr    = crdappsv1.GroupVersion.String()
)

// BookReconciler reconciles a Book object
type BookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.naeem4265.com,resources=books,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.naeem4265.com,resources=books/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.naeem4265.com,resources=books/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Book object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *BookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Load the book by name
	// TODO(user): your logic here
	var book crdappsv1.Book
	if err := r.Get(ctx, req.NamespacedName, &book); err != nil {
		log.Error(err, "Unable to fetch book")
		// We'll ignore not-found errors, since they can't be fixed by an immediate requeue
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var deploy appsv1.Deployment
	depName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      book.Spec.DeploymentName,
	}
	// get the deployments, and update the status
	if err := r.Get(ctx, depName, &deploy); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		// if no Deployment found , or found deployments are not owned by book, create one on the cluster
		deploy := newDeployment(&book)
		if err := r.Create(ctx, &deploy); err != nil {
			log.Error(err, "Unable to create deployment")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		fmt.Println("Created deployment successfully")
	} else {
		// update deployment status if exist
		if deploy.Spec.Replicas != book.Spec.Replicas {
			deploy.Spec.Replicas = book.Spec.Replicas
			if err := r.Update(ctx, &deploy); err != nil {
				log.Error(err, "Unable to update deployment")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		}
		// update book status
		book.Status.AvailableReplicas = deploy.Status.AvailableReplicas
		_ = r.Status().Update(ctx, &book)
	}

	// Same for service.
	var srv corev1.Service
	if err := r.Get(ctx, depName, &srv); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		// if no Service found , or found deployments are not owned by book, create one on the cluster
		srv := customService(&book)
		if err := r.Create(ctx, &srv); err != nil {
			log.Error(err, "Unable to create Service")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		fmt.Println("Created Service successfully")
	}

	fmt.Println("Reconcile done")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	/*
		// Did not use.
		if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, depOwnerKey, func(rawObj client.Object) []string {
			// grab the deploy object, extract the owner...
			deploy := rawObj.(*appsv1.Deployment)
			owner := metav1.GetControllerOf(deploy)
			if owner == nil {
				return nil
			}
			// ...make sure it's a Nicedeploy...
			if owner.APIVersion != apiGVStr || owner.Kind != "Neymar" {
				return nil
			}

			// ...and if so, return it
			return []string{owner.Name}
		}); err != nil {
			return err
		}

		if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Service{}, svcOwnerKey, func(rawObj client.Object) []string {
			// grab the service object, extract the owner...
			svc := rawObj.(*corev1.Service)
			owner := metav1.GetControllerOf(svc)
			if owner == nil {
				return nil
			}
			if owner.APIVersion != apiGVStr || owner.Kind != "Neymar" {
				return nil
			}
			// ...and if so, return it
			return []string{owner.Name}
		}); err != nil {
			return err
		}
	*/
	fmt.Println("SetupWithManager successful. ")
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdappsv1.Book{}).
		// Watch deployment and if owner of this deployment is book, then call Reconcile.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		// Watches()
		Complete(r)
}
