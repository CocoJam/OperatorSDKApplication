package operatorInterface
import(
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"

)
var _ reconcile.Reconciler = &ReconcileBrokerOperator{}

type ReconcileBrokerOperator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
	ReconcileFunc reconcileFunction
}
type reconcileFunction func(r *ReconcileBrokerOperator,request reconcile.Request) (reconcile.Result, error)

func (r *ReconcileBrokerOperator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	return r.ReconcileFunc(r,request)
}

func (r *ReconcileBrokerOperator) GetBrokerInstance (instance *kafkav1alpha1.BrokerOperator, request reconcile.Request) (reconcile.Result,error){
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
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
 