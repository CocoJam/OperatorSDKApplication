package operatorInterface
import(
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"

)
var _ reconcile.Reconciler = &ReconcileZooKeeperOperator{}

type ReconcileZooKeeperOperator struct {
	Client client.Client
	Scheme *runtime.Scheme
	ReconcileFunc reconcileZooKeeperFunction
}
type reconcileZooKeeperFunction func(r *ReconcileZooKeeperOperator,request reconcile.Request) (reconcile.Result, error)

func (r *ReconcileZooKeeperOperator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	return r.ReconcileFunc(r,request)
}

func (r *ReconcileZooKeeperOperator) GetZooKeeperInstance (instance *kafkav1alpha1.ZooKeeperOperator, request reconcile.Request) (reconcile.Result,error){
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, err
}
 