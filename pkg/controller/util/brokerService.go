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
	kss "github.com/example-inc/app-operator/pkg/controller/templating/KafkaStatefulSet"
)


type BrokerService struct{
	ResourcePtr *corev1.Service
	OperatorPtr *kafkav1alpha1.BrokerOperator 
	R *ReconcileBrokerOperator
	KafkaService *kss.KafkaService
	Headless bool
}

func(bss BrokerService) getResourcePtr()interface{}{
	return bss.OperatorPtr
}

func (bss BrokerService) getOperatorPtr() *kafkav1alpha1.BrokerOperator{
	return bss.OperatorPtr
}

func (bss BrokerService) getReconcileOperator()interface{}{
	return bss.R
}

func (bss BrokerService) getResourceNameSpace() (string,string){
	return bss.ResourcePtr.Name, bss.ResourcePtr.Namespace
}

func (bss BrokerService) findResourceFromInstance()(recon reconcile.Result,err error){
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

func (bs BrokerService) deployment(recon reconcile.Result,err error) (reconcile.Result,error){
	if err != nil && errors.IsNotFound(err) {
		ksTemplate:= kss.KafkaService{}
		bs.KafkaService = &ksTemplate
		dep := ksTemplate.BootStrap(bs.OperatorPtr, bs.Headless)
		bs.ResourcePtr = &dep
		err = bs.R.Client.Create(context.TODO(), bs.ResourcePtr)
		return recon, err
	}
	return recon,err
}

func (r BrokerService) SpecConditionalUpdate(con ConditionSpecUpdate,condition bool,reqLogger logr.Logger)  (reconcile.Result,error){
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

func (bss BrokerService)GetPodList(podList *corev1.PodList) (error){
	labelSelector := labels.SelectorFromSet(labelsForBroker(bss.ResourcePtr.Name))
	listOps := &client.ListOptions{Namespace: bss.ResourcePtr.Namespace, LabelSelector: labelSelector}
	return bss.R.Client.List(context.TODO(), listOps, podList)
}

var _ ResourceGetDeploy = (*BrokerService)(nil)

func labelsForBroker(name string) map[string]string {
	return 	map[string]string{"app": "Broker", "Kafka_Broker_cr": name}
}
