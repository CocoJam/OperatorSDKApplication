package zookeeper
import (
	corev1 "k8s.io/api/core/v1"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	"github.com/example-inc/app-operator/pkg/controller/templating/templates"
)

type ZooKeeperService struct{
	templateSS templates.Service
}

type DefaultZooKeeperService struct{
	templateSS templates.Service
}

func(defaultService *DefaultZooKeeperService) bootStrap(zookeeper *kafkav1alpha1.ZooKeeperOperator, headless bool) corev1.Service {
	ls:= map[string]string{"app": "Zookeeper", "Zookeeper_cr": zookeeper.Name}
	Service := templates.Service{}
	meta:= templates.DeploymentMetaTemplate{
		Kind: "Service",
		APIVersion: "v1",
		ObjectName: zookeeper.Name,
		ObjectNamespace: zookeeper.Namespace,
		Labels: ls,
	}
	Service.Meta = meta
	Service.Selector(ls)
	// Service.ServiceType("")
	Service.ServicePort(zookeeper.Spec.ContainerPorts)
	if(headless){
		meta.ObjectName +="-headless"
		Service.ServiceSpec("None")
	}
	return Service.BootStrap()
}

type DefaultkafkaHeadlessService struct{
	templateSS templates.Service
}