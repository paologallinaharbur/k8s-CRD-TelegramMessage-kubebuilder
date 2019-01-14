/*
Copyright 2018 Paolo.Gallina.

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

package telegrammessage

import (
	"context"
	"log"
	"net/http"
	"net/url"

	harburv1beta1 "message/pkg/apis/harbur/v1beta1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new TelegramMessage Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTelegramMessage{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("telegrammessage-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to TelegramMessage
	err = c.Watch(&source.Kind{Type: &harburv1beta1.TelegramMessage{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Uncomment watch a Deployment created by TelegramMessage - change this for objects you create
	err = c.Watch(&source.Kind{Type: &harburv1beta1.TelegramMessage{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &harburv1beta1.TelegramMessage{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileTelegramMessage{}

// ReconcileTelegramMessage reconciles a TelegramMessage object
type ReconcileTelegramMessage struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TelegramMessage object and makes changes based on the state read
// and what is in the TelegramMessage.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=harbur.harbur.io,resources=telegrammessages,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileTelegramMessage) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the TelegramMessage instance
	instance := &harburv1beta1.TelegramMessage{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	log.Printf("Checking status of resource %s", instance.Name)

	if instance.Status.Delivered != "Yes" || instance.Status.MessageDelivered != instance.Spec.MessageToDeliver {
		log.Printf("Sending message %s\n", instance.Spec.MessageToDeliver)
		resp, err := http.Post("https://api.telegram.org/bot"+instance.Spec.Token+"/sendMessage?chat_id="+instance.Spec.ChatID+"&text="+url.QueryEscape(instance.Spec.MessageToDeliver), "", nil)
		if err != nil {
			log.Printf("Please doublecheck token and chatId, %s", err)
			instance.Status.Delivered = "No"
			instance.Status.MessageDelivered = "ERROR NO MESSAGE SENT"
		} else {
			instance.Status.Delivered = "Yes"
			instance.Status.MessageDelivered = instance.Spec.MessageToDeliver
		}
		defer resp.Body.Close()
		err = r.Update(context.TODO(), instance)
		if err != nil {
			log.Printf("%s\n", err)
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
