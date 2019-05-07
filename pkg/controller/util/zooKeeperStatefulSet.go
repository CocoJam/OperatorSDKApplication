package operatorInterface
import (
	"context"
	"reflect"
	"fmt"
	// er "errors"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/go-logr/logr"
	kss "github.com/example-inc/app-operator/pkg/controller/templating/ZooKeeperStatefulSet"
)
type ZooKeeperStatefulSet struct{
	ResourcePtr *appsv1.StatefulSet
	OperatorPtr *kafkav1alpha1.ZooKeeperOperator 
	R *ReconcileZooKeeperOperator
	KafkaStatefulSet *kss.ZooKeeperStatefulSet
}

func(bss ZooKeeperStatefulSet) getResourcePtr() interface{}{
	return bss.ResourcePtr
}

func (bss ZooKeeperStatefulSet) getOperatorPtr() interface{}{
	return bss.OperatorPtr
}

func (bss ZooKeeperStatefulSet) getReconcileOperator() interface{}{
	return bss.R
}
func (bss ZooKeeperStatefulSet) getResourceNameSpace() (string,string){
	return bss.ResourcePtr.Name, bss.ResourcePtr.Namespace
}

func (bss ZooKeeperStatefulSet) findResourceFromInstance()( reconcile.Result,error){
	err := bss.R.Client.Get(context.TODO(), types.NamespacedName{Name: bss.OperatorPtr.Name, Namespace: bss.OperatorPtr.Namespace}, bss.ResourcePtr)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, err
}


func (bss ZooKeeperStatefulSet) deployment(recon reconcile.Result,err error) (reconcile.Result,error){
	if err != nil && errors.IsNotFound(err) {
		kssTemplate:= kss.ZooKeeperStatefulSet{}
		bss.KafkaStatefulSet = &kssTemplate
		dep := kssTemplate.BootStrap(bss.OperatorPtr)
		bss.ResourcePtr = &dep
		err = bss.R.Client.Create(context.TODO(), bss.ResourcePtr)
		return recon, err
	}
	return recon,err
}


func (r ZooKeeperStatefulSet) SpecConditionalUpdate(con ConditionSpecUpdate,condition bool,reqLogger logr.Logger)  (reconcile.Result,error){
	if condition {
		ptrValueOf:= reflect.ValueOf(r.getResourcePtr())
		ptrType:= ptrValueOf.Type()
		Name, Namespace:= r.getResourceNameSpace()
		err := con(r)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", fmt.Sprintf("%T.Namespace",ptrType), Namespace, fmt.Sprintf("%T.Name",ptrType), Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

func (bss ZooKeeperStatefulSet)GetPodList(podList *corev1.PodList) (error){
	labelSelector := labels.SelectorFromSet(labelsForZooKeeper(bss.ResourcePtr.Name))
	listOps := &client.ListOptions{Namespace: bss.ResourcePtr.Namespace, LabelSelector: labelSelector}
	return bss.R.Client.List(context.TODO(), listOps, podList)
}

func (bss ZooKeeperStatefulSet)GetPodListByLabel(podList *corev1.PodList, ls map[string]string) (error){
	labelSelector := labels.SelectorFromSet(ls)
	listOps := &client.ListOptions{Namespace: bss.ResourcePtr.Namespace, LabelSelector: labelSelector}
	return bss.R.Client.List(context.TODO(), listOps, podList)
}

var _ ResourceGetDeploy = (*ZooKeeperStatefulSet)(nil)
