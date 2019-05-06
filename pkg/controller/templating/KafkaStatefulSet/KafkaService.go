package broker
import (
	corev1 "k8s.io/api/core/v1"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	"github.com/example-inc/app-operator/pkg/controller/templating/templates"
)

type KafkaService struct{
	templateSS templates.Service
}

type DefaultkafkaService struct{
	templateSS templates.Service
}

func(kafkaService *KafkaService) BootStrap(broker *kafkav1alpha1.BrokerOperator, headless bool) corev1.Service {
	ls:= map[string]string{"app": "Broker", "Kafka_Broker_cr": broker.Name}
	Service := templates.Service{}
	meta:= templates.DeploymentMetaTemplate{
		Kind: "Service",
		APIVersion: "v1",
		ObjectName: broker.Name,
		ObjectNamespace: broker.Namespace,
		Labels: ls,
	}
	Service.Meta = meta
	Service.Selector(ls)
	// Service.ServiceType("")
	Service.ServicePort(broker.Spec.ContainerPorts)
	if(headless){
		meta.ObjectName +="-headless"
		Service.ServiceSpec("None")
	}
	return Service.BootStrap()
}

type DefaultkafkaHeadlessService struct{
	templateSS templates.Service
}