package operatorInterface
import (
	"context"
	"reflect"
	"fmt"
	// er "errors"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/go-logr/logr"
	kss "github.com/example-inc/app-operator/pkg/controller/templating/ZooKeeperStatefulSet"
)


type ZooKeeperService struct{
	ResourcePtr *corev1.Service
	OperatorPtr *kafkav1alpha1.ZooKeeperOperator 
	R *ReconcileZooKeeperOperator
	KafkaService *kss.ZooKeeperService
	Headless bool
}

func(bss ZooKeeperService) getResourcePtr()interface{}{
	return bss.OperatorPtr
}

func (bss ZooKeeperService) getOperatorPtr() interface{}{
	return bss.OperatorPtr
}

func (bss ZooKeeperService) getReconcileOperator()interface{}{
	return bss.R
}

func (bss ZooKeeperService) getResourceNameSpace() (string,string){
	return bss.ResourcePtr.Name, bss.ResourcePtr.Namespace
}

func (bss ZooKeeperService) findResourceFromInstance()(recon reconcile.Result,err error){
	if bss.Headless{
		err = bss.R.Client.Get(context.TODO(), types.NamespacedName{Name: bss.OperatorPtr.Name+"-headless", Namespace: bss.OperatorPtr.Namespace}, bss.ResourcePtr)
	}else{
		err = bss.R.Client.Get(context.TODO(), types.NamespacedName{Name: bss.OperatorPtr.Name, Namespace: bss.OperatorPtr.Namespace}, bss.ResourcePtr)
	}
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, err
}

func (bs ZooKeeperService) deployment(recon reconcile.Result,err error) (reconcile.Result,error){
	if err != nil && errors.IsNotFound(err) {
		ksTemplate:= kss.ZooKeeperService{}
		bs.KafkaService = &ksTemplate
		dep := ksTemplate.BootStrap(bs.OperatorPtr, bs.Headless)
		bs.ResourcePtr = &dep
		err = bs.R.Client.Create(context.TODO(), bs.ResourcePtr)
		return recon, err
	}
	return recon,err
}

func (r ZooKeeperService) SpecConditionalUpdate(con ConditionSpecUpdate,condition bool,reqLogger logr.Logger)  (reconcile.Result,error){
	if condition {
		err := con(r)
		if err != nil {
			ptrValueOf:= reflect.ValueOf(r.getResourcePtr())
			ptrType:= ptrValueOf.Type()
			Name, Namespace:= r.getResourceNameSpace()
			reqLogger.Error(err, "Failed to create new Deployment", fmt.Sprintf("%T.Namespace",ptrType), Namespace, fmt.Sprintf("%T.Name",ptrType), Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

func (bss ZooKeeperService)GetPodList(podList *corev1.PodList) (error){
	labelSelector := labels.SelectorFromSet(labelsForZooKeeper(bss.ResourcePtr.Name))
	listOps := &client.ListOptions{Namespace: bss.ResourcePtr.Namespace, LabelSelector: labelSelector}
	return bss.R.Client.List(context.TODO(), listOps, podList)
}

func (bss ZooKeeperService)GetPodListByLabel(podList *corev1.PodList, ls map[string]string) (error){
	labelSelector := labels.SelectorFromSet(ls)
	listOps := &client.ListOptions{Namespace: bss.ResourcePtr.Namespace, LabelSelector: labelSelector}
	return bss.R.Client.List(context.TODO(), listOps, podList)
}

var _ ResourceGetDeploy = (*ZooKeeperService)(nil)

func labelsForZooKeeper(name string) map[string]string {
	return 	map[string]string{"app": "ZooKeeper", "ZooKeeper_cr": name}
}
