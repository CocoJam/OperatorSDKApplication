package operatorInterface
import (
	"fmt"
	"reflect"
	// kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"github.com/go-logr/logr"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net"
	"os"
)

type ConditionSpecUpdate func(r ResourceGetDeploy) error

type ResourceGetDeploy interface{
	getResourcePtr() interface{}
	getOperatorPtr() interface{}
	getReconcileOperator() interface{}
	getResourceNameSpace()(string,string)
	findResourceFromInstance()(reconcile.Result,error)
	deployment(recon reconcile.Result,err error) (reconcile.Result,error)
	SpecConditionalUpdate(con ConditionSpecUpdate, condition bool,reqLogger logr.Logger)  (reconcile.Result,error)
	GetPodList(podList *corev1.PodList) (error)
	GetPodListByLabel(podList *corev1.PodList, ls map[string]string) (error)
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


func execCommand (namespace string, podName string, stdinReader io.Reader, container *v1.Container, command ...string) (string, error) {
	
	execReq := kubeClient.CoreV1().RESTClient().Post()
	execReq = execReq.Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")

	execReq.VersionedParams(&v1.PodExecOptions{
		Container: container.Name,
		Command:   command,
		Stdout:    true,
		Stderr:    true,
		Stdin:     stdinReader != nil,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(inClusterConfig, "POST", execReq.URL())

	if err != nil {
		reqLogger.Error("Creating remote command executor failed: %v", err)
		return "", err
	}

	stdOut := bytes.Buffer{}
	stdErr := bytes.Buffer{}

	reqLogger.Debugf("Executing command '%v' in namespace='%s', pod='%s', container='%s'", command, namespace, podName, container.Name)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: bufio.NewWriter(&stdOut),
		Stderr: bufio.NewWriter(&stdErr),
		Stdin:  stdinReader,
		Tty:    false,
	})

	reqLogger.Debugf("Command stderr: %s", stdErr.String())
	reqLogger.Debugf("Command stdout: %s", stdOut.String())

	if err != nil {
		reqLogger.Infof("Executing command failed with: %v", err)

		return "", err
	}

	reqLogger.Debug("Command succeeded.")
	if stdErr.Len() > 0 {
		return "", fmt.Errorf("stderr: %v", stdErr.String())
	}

	return stdOut.String(), nil

}
