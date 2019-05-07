package zookeeperoperator

import (
	"context"

	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	util "github.com/example-inc/app-operator/pkg/controller/util"
)

var log = logf.Log.WithName("controller_zookeeperoperator")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ZooKeeperOperator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &util.ReconcileZooKeeperOperator{Client: mgr.GetClient(), Scheme: mgr.GetScheme(),ReconcileFunc: Reconcile }
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("zookeeperoperator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ZooKeeperOperator
	err = c.Watch(&source.Kind{Type: &kafkav1alpha1.ZooKeeperOperator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ZooKeeperOperator
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kafkav1alpha1.ZooKeeperOperator{},
	})
	if err != nil {
		return err
	}

	return nil
}

// Reconcile reads that state of the cluster for a ZooKeeperOperator object and makes changes based on the state read
// and what is in the ZooKeeperOperator.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func Reconcile(r *util.ReconcileZooKeeperOperator,request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ZooKeeperOperator")

	// Fetch the ZooKeeperOperator instance
	instance := &kafkav1alpha1.ZooKeeperOperator{}
	recon , err := r.GetZooKeeperInstance(instance, request)
	if(util.ParseReconcile(recon,err)){
		return recon, err
	}
	var bss = util.ZooKeeperStatefulSet{ResourcePtr: &appsv1.StatefulSet{}, OperatorPtr: instance, R: r}
	recon, err = util.GetResourceInstanceDeploy(bss, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bss.ResourcePtr, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	size := instance.Spec.Replicas

	ConditionSpecUpdate := func(resource util.ResourceGetDeploy) error{
		*bss.ResourcePtr.Spec.Replicas = size
		return r.Client.Update(context.TODO(), bss.ResourcePtr)
	}

	recon, err = bss.SpecConditionalUpdate(ConditionSpecUpdate,
		 *bss.ResourcePtr.Spec.Replicas != size, reqLogger)
	if(util.ParseReconcile(recon,err)){
			return recon, err
	}

	var bs = util.ZooKeeperService{ResourcePtr: &corev1.Service{}, OperatorPtr: instance, R : r, Headless: false}
	recon, err = util.GetResourceInstanceDeploy(bs, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bs.ResourcePtr, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	bs.Headless = true
	recon, err = util.GetResourceInstanceDeploy(bs, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bs.ResourcePtr, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	podList := &corev1.PodList{}
	err = bss.GetPodList(podList)

	return reconcile.Result{}, nil
}