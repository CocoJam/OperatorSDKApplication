package brokeroperator

import (
	"context"
	// er "errors"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type MetaDeployment interface{
	TypeMeta() metav1.TypeMeta
	ObjectMeta() metav1.ObjectMeta
}

type DeploymentMetaTemplate struct{
	Kind string
	APIVersion string
	ObjectName string
	ObjectNamespace string
}

func (dmt *DeploymentMetaTemplate) TypeMeta() metav1.TypeMeta{
	return metav1.TypeMeta{
		Kind: dmt.Kind,
		APIVersion: dmt.APIVersion,
	}
}

func (dmt *DeploymentMetaTemplate) ObjectMeta() metav1.ObjectMeta{
	return metav1.ObjectMeta{
		Name: dmt.ObjectNamespace,
		Namespace: dmt.ObjectNamespace,
	}
}

type Selector interface{
	Selector() *metav1.LabelSelector
}

type SpecDeployment interface{
	Selector() *metav1.LabelSelector
	PersistentVolumeClaim() corev1.PersistentVolumeClaim
}

type DeploymentSpecTemplate struct{
	Replicas int32
	Selector map[string]string
	Template corev1.PodTemplateSpec
	VolumeClaimTemplates []corev1.PersistentVolumeClaim
	ServiceName string
	PodManagementPolicy appsv1.PodManagementPolicyType
	UpdateStrategy 	appsv1.StatefulSetUpdateStrategy
}

func (sd *DeploymentSpecTemplate) PersistentVolumeClaim() corev1.PersistentVolumeClaim{
	return corev1.PersistentVolumeClaim{}
}

type PersistentVolumeClaimSpec interface{
	AccessModes(AccessModes[]string)
	Selector(selets map[string]string) *metav1.LabelSelector
	Resources() corev1.ResourceRequirements
	VolumeName() string 
	StorageClassName() *string
	VolumeMode() *corev1.PersistentVolumeMode
}

func (pvc *PersistentVolumeClaim) AccessModes(AccessModes[]string){
	sliceAccess := make([]corev1.PersistentVolumeAccessMode, len(AccessModes))
	for i,access := range AccessModes{
		switch access {
		case string(corev1.ReadWriteOnce):
			sliceAccess[i] = corev1.ReadWriteOnce
		case string(corev1.ReadOnlyMany):
			sliceAccess[i] = corev1.ReadOnlyMany
		case string(corev1.ReadWriteMany):
			sliceAccess[i] = corev1.ReadWriteMany
		default:
		}
	}
	pvc.Spec.AccessModes = sliceAccess
}

func (pvc *PersistentVolumeClaim) Selector(selets map[string]string){
	pvc.Spec.Selector = &metav1.LabelSelector{MatchLabels: selets}
}

func (pvc *PersistentVolumeClaim) Resource(Resources map[string]string ){
	resourceMap := make(map[corev1.ResourceName]resource.Quantity, len(Resources))
	
	for k,v := range Resources{
		switch k {
		case string(corev1.ResourceCPU):
			resourceMap[corev1.ResourceCPU] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceMemory):
			resourceMap[corev1.ResourceMemory] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceStorage):
			resourceMap[corev1.ResourceStorage] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceEphemeralStorage):
			resourceMap[corev1.ResourceEphemeralStorage] = resource.Quantity{Format: resource.DecimalSI}
		default:

		}
	}
	pvc.Spec.Resources.Requests = resourceMap
}

func (pvc *PersistentVolumeClaim) VolumeName(vol string){
	pvc.Spec.VolumeName = vol
}
func (pvc *PersistentVolumeClaim) StorageClassName(storage *string){
	pvc.Spec.StorageClassName = storage
}
func (pvc *PersistentVolumeClaim) VolumeMode(mode string){
	switch mode {
	case string(corev1.PersistentVolumeBlock):
		pvc.Spec.VolumeMode = &corev1.PersistentVolumeBlock
	case string(corev1.PersistentVolumeFilesystem):
		pvc.Spec.VolumeMode = &corev1.PersistentVolumeFilesystem
	default:

	}
}


type PersistentVolumeClaim struct{
	dmt *DeploymentMetaTemplate
	Spec corev1.PersistentVolumeClaimSpec
	Status corev1.PersistentVolumeClaimStatus
}


type ResourceFunction interface{
	findResourceFromInstance() (reconcile.Result,error)
	deployment()(reconcile.Result, error)
}

type BrokerStatefulSet struct{
	resourcePtr *appsv1.StatefulSet
	operatorPtr *kafkav1alpha1.BrokerOperator 
	r *ReconcileBrokerOperator
}



func (bss *BrokerStatefulSet) findResourceFromInstance()( reconcile.Result,error){
	err := bss.r.client.Get(context.TODO(), types.NamespacedName{Name: bss.operatorPtr.Name, Namespace: bss.operatorPtr.Namespace}, bss.resourcePtr)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, err
}


func (bss *BrokerStatefulSet) deployment(reconcile.Result, error){
	ss := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{

		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: appsv1.StatefulSetSpec{},
		Status: appsv1.StatefulSetStatus{},
	}
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
	recon, err = bss.findResourceFromInstance()
	
	if(parseReconcile(recon,err)){
		return recon, err
	}else{
		
	}

	
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMemcached(memcached)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}



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