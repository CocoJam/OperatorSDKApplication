package operatorInterface
import (
	"fmt"
	"reflect"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/go-logr/logr"
)

type ConditionSpecUpdate func(r ResourceGetDeploy) error

type ResourceGetDeploy interface{
	getResourcePtr() interface{}
	getOperatorPtr() *kafkav1alpha1.BrokerOperator
	getReconcileOperator() interface{}
	getResourceNameSpace()(string,string)
	findResourceFromInstance()(reconcile.Result,error)
	deployment(recon reconcile.Result,err error) (reconcile.Result,error)
	SpecConditionalUpdate(con ConditionSpecUpdate, condition bool,reqLogger logr.Logger)  (reconcile.Result,error)
	GetPodList(podList *corev1.PodList) (error)
}


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
	if ParseReconcile(recon,err){
		reqLogger.Error(err, "Failed to create new Deployment", fmt.Sprintf("%T.Namespace",ptrType), Namespace, fmt.Sprintf("%T.Name",ptrType), Name)
		return recon, err
	}
	return reconcile.Result{}, err
}

func ParseReconcile(recon reconcile.Result, err error) bool{
	if !recon.Requeue || err != nil {
		return true
	} else {
		return false
	}
}
