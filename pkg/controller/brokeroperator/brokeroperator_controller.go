package brokeroperator

import (
	"context"
	"fmt"
	"reflect"
	// er "errors"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"github.com/go-logr/logr"
	kss "github.com/example-inc/app-operator/pkg/controller/templating/KafkaStatefulSet"
)

var log = logf.Log.WithName("controller_brokeroperator")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new BrokerOperator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBrokerOperator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("brokeroperator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource BrokerOperator
	err = c.Watch(&source.Kind{Type: &kafkav1alpha1.BrokerOperator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner BrokerOperator
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kafkav1alpha1.BrokerOperator{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileBrokerOperator{}

// ReconcileBrokerOperator reconciles a BrokerOperator object
type ReconcileBrokerOperator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

type ConditionSpecUpdate = func(r ResourceGetDeploy) error


type ResourceGetDeploy interface{
	getResourcePtr() interface{}
	getOperatorPtr() *kafkav1alpha1.BrokerOperator
	getReconcileOperator() interface{}
	getResourceNameSpace()(string,string)
	findResourceFromInstance()(reconcile.Result,error)
	deployment(recon reconcile.Result,err error) (reconcile.Result,error)
	// ConditionalReconcile(condition bool,)
	specConditionalUpdate(con ConditionSpecUpdate, condition bool,reqLogger logr.Logger)  (reconcile.Result,error)
}

type BrokerStatefulSet struct{
	resourcePtr *appsv1.StatefulSet
	operatorPtr *kafkav1alpha1.BrokerOperator 
	r *ReconcileBrokerOperator
	KafkaStatefulSet *kss.KafkaStatefulSet
}

func(bss BrokerStatefulSet) getResourcePtr() interface{}{
	return bss.resourcePtr
}

func (bss BrokerStatefulSet) getOperatorPtr() *kafkav1alpha1.BrokerOperator{
	return bss.operatorPtr
}

func (bss BrokerStatefulSet) getReconcileOperator() interface{}{
	return bss.r
}
func (bss BrokerStatefulSet) getResourceNameSpace() (string,string){
	return bss.resourcePtr.Name, bss.resourcePtr.Namespace
}

func (bss BrokerStatefulSet) findResourceFromInstance()( reconcile.Result,error){
	err := bss.r.client.Get(context.TODO(), types.NamespacedName{Name: bss.operatorPtr.Name, Namespace: bss.operatorPtr.Namespace}, bss.resourcePtr)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, err
}


func (bss BrokerStatefulSet) deployment(recon reconcile.Result,err error) (reconcile.Result,error){
	if err != nil && errors.IsNotFound(err) {
		kssTemplate:= kss.KafkaStatefulSet{}
		bss.KafkaStatefulSet = &kssTemplate
		dep := kssTemplate.BootStrap(bss.operatorPtr)
		bss.resourcePtr = &dep
		err = bss.r.client.Create(context.TODO(), bss.resourcePtr)
		return recon, err
	}
	return recon,err
}


func (r BrokerStatefulSet) specConditionalUpdate(con ConditionSpecUpdate,condition bool,reqLogger logr.Logger)  (reconcile.Result,error){
	if condition {
		ptrValueOf:= reflect.ValueOf(r.getResourcePtr())
		ptrType:= ptrValueOf.Type()
		Name, Namespace:= r.getResourceNameSpace()
		err := con(r)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", Namespace, "Deployment.Name", Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

var _ ResourceGetDeploy = (*BrokerStatefulSet)(nil)

type BrokerService struct{
	resourcePtr *corev1.Service
	operatorPtr *kafkav1alpha1.BrokerOperator 
	r *ReconcileBrokerOperator
	KafkaService *kss.KafkaService
	Headless bool
}

func(bss BrokerService) getResourcePtr()interface{}{
	return bss.resourcePtr
}

func (bss BrokerService) getOperatorPtr() *kafkav1alpha1.BrokerOperator{
	return bss.operatorPtr
}

func (bss BrokerService) getReconcileOperator()interface{}{
	return bss.r
}

func (bss BrokerService) getResourceNameSpace() (string,string){
	return bss.resourcePtr.Name, bss.resourcePtr.Namespace
}

func (bss BrokerService) findResourceFromInstance()(reconcile.Result,error){
	var err error
	if bss.Headless{
		err := bss.r.client.Get(context.TODO(), types.NamespacedName{Name: bss.operatorPtr.Name+"-headless", Namespace: bss.operatorPtr.Namespace}, bss.resourcePtr)
	}else{
		err := bss.r.client.Get(context.TODO(), types.NamespacedName{Name: bss.operatorPtr.Name, Namespace: bss.operatorPtr.Namespace}, bss.resourcePtr)
	}
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, err
}


func (bs BrokerService) deployment(recon reconcile.Result,err error) (reconcile.Result,error){
	if err != nil && errors.IsNotFound(err) {
		ksTemplate:= kss.KafkaService{}
		bs.KafkaService = &ksTemplate
		dep := ksTemplate.BootStrap(bs.operatorPtr, bs.Headless)
		bs.resourcePtr = &dep
		err = bs.r.client.Create(context.TODO(), bs.resourcePtr)
		return recon, err
	}
	return recon,err
}

func (r BrokerService) specConditionalUpdate(con ConditionSpecUpdate,condition bool,reqLogger logr.Logger)  (reconcile.Result,error){
	if condition {
		ptrValueOf:= reflect.ValueOf(r.getResourcePtr())
		ptrType:= ptrValueOf.Type()
		Name, Namespace:= r.getResourceNameSpace()
		err := con(r)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", Namespace, "Deployment.Name", Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

var _ ResourceGetDeploy = (*BrokerService)(nil)

func GetResourceInstanceDeploy(r ResourceGetDeploy, reqLogger logr.Logger) (reconcile.Result,error){
	recon, err := r.findResourceFromInstance()
	ptrValueOf:= reflect.ValueOf(r.getResourcePtr())
	ptrType:= ptrValueOf.Type()
	Name, Namespace:= r.getResourceNameSpace()
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment",fmt.Sprintf("%T.Namespace",ptrType), Namespace, fmt.Sprintf("%T.Name",ptrType), Name)
		recon, err = r.deployment(recon,err)
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}
	if parseReconcile(recon,err){
		reqLogger.Error(err, "Failed to create new Deployment", fmt.Sprintf("%T.Namespace",ptrType), Namespace, fmt.Sprintf("%T.Name",ptrType), Name)
		return recon, err
	}
	return reconcile.Result{}, err
}

// Reconcile reads that state of the cluster for a BrokerOperator object and makes changes based on the state read
// and what is in the BrokerOperator.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBrokerOperator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling BrokerOperator")

	// Fetch the BrokerOperator instance
	instance := &kafkav1alpha1.BrokerOperator{}
	recon , err := r.getBrokerInstance(instance, request)
	if(parseReconcile(recon,err)){
		return recon, err
	}

	var bss = BrokerStatefulSet{resourcePtr: &appsv1.StatefulSet{}, operatorPtr: instance, r : r}
	recon, err = GetResourceInstanceDeploy(bss, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bss.resourcePtr, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	size := instance.Spec.Replicas
	// ConditionSpecUpdate := func(r ResourceGetDeploy) error{
		*bss.resourcePtr.Spec.Replicas = size
		// err = r.client.Update(context.TODO(), found)
	// }
	// if *bss.resourcePtr.Spec.Replicas != size {
	// 	*bss.resourcePtr.Spec.Replicas = &size
	// 	err = r.client.Update(context.TODO(), found)
	// 	if err != nil {
	// 		reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	// 		return reconcile.Result{}, err
	// 	}
	// 	return reconcile.Result{Requeue: true}, nil
	// }

	var bs = BrokerService{resourcePtr: &corev1.Service{}, operatorPtr: instance, r : r, Headless: false}
	recon, err = GetResourceInstanceDeploy(bs, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bs.resourcePtr, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	bs.Headless = true
	recon, err = GetResourceInstanceDeploy(bs, reqLogger)
	if err := controllerutil.SetControllerReference(instance, bs.resourcePtr, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForBroker(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set BrokerOperator instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func parseReconcile(recon reconcile.Result, err error) bool{
	if !recon.Requeue || err != nil {
		return true
	} else {
		return false
	}
}




// func (bss *BrokerStatefulSet) findDeploymentResourceFromInstance(instance *kafkav1alpha1.BrokerOperator ,resource *appsv1.Deployment) ( reconcile.Result,error){
// 	err := bss.r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, resource)
// 	if err != nil && errors.IsNotFound(err) {
// 		return reconcile.Result{}, err
// 	}
// 	return reconcile.Result{Requeue: true}, err
// }

func (r *ReconcileBrokerOperator) getBrokerInstance (instance *kafkav1alpha1.BrokerOperator, request reconcile.Request) (reconcile.Result,error){
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, err
}
 
// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *kafkav1alpha1.BrokerOperator) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
func labelsForBroker(name string) map[string]string {
	return 	map[string]string{"app": "Broker", "Kafka_Broker_cr": name}
}
