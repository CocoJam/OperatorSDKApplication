package broker
import (
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"github.com/example-inc/app-operator/pkg/controller/templating/templates"
)

type KafkaStatefulSet struct{
	templateSS templates.StatefulSet
}

type DefaultKafkaStatefulSet struct{
	templateSS templates.StatefulSet
}

func (defaultBroker *KafkaStatefulSet) BootStrap(broker *kafkav1alpha1.BrokerOperator) appsv1.StatefulSet{
	ls:= map[string]string{"app": "Broker", "Kafka_Broker_cr": broker.Name}
	meta:= templates.DeploymentMetaTemplate{
		Kind: "StatefulSet",
		APIVersion: "apps/v1beta1",
		ObjectName: broker.Name,
		ObjectNamespace: broker.Namespace,
	}
	defaultBroker.templateSS.Meta = meta

	defaultBroker.templateSS.Replicas(broker.Spec.Replicas)
	defaultBroker.templateSS.SS.Spec.Selector= &metav1.LabelSelector{MatchLabels: ls,}

	PodManagementPolicy := templates.PodManagementPolicy{}
	PodManagementPolicy.PodManagementPolicy("OrderedReady")
	defaultBroker.templateSS.SS.Spec.PodManagementPolicy = PodManagementPolicy.PMP
	
	UpdateStrategy:= templates.UpdateStrategy{}
	UpdateStrategy.UpdateStrategy("RollingUpdate",0)
	defaultBroker.templateSS.SS.Spec.UpdateStrategy = UpdateStrategy.US

	defaultBroker.templateSS.SS.Spec.ServiceName = ""
	
	PodTemplateSpec := templates.PodTemplateSpec{}
	PodTemplateSpec.MetaPodTemplateSpecs(ls)

	
	ContainerAssemble := templates.ContainerASSemble{}
	containerCompose(broker,&ContainerAssemble)
	defaultBroker.templateSS.SS.Spec.Template.Spec.Containers = []corev1.Container{ContainerAssemble.Container}

	PersistentVolumeClaim:= templates.PersistentVolumeClaimASSemble{}
	PersistentVolumeClaimSlice := make([]corev1.PersistentVolumeClaim, broker.Spec.MountNum)
	for i:=0;i < broker.Spec.MountNum; i++{
		pvcCompose(broker, &PersistentVolumeClaim, i)
		PersistentVolumeClaimSlice[i] = PersistentVolumeClaim.PVC
	}
	defaultBroker.templateSS.SS.Spec.VolumeClaimTemplates = PersistentVolumeClaimSlice
	return defaultBroker.templateSS.BootStrap()
}

func containerCompose(broker *kafkav1alpha1.BrokerOperator,ContainerAss *templates.ContainerASSemble){
	ContainerAssemble := *ContainerAss
	ContainerAssemble.ImageFactory(broker.Spec.ContainerName,broker.Spec.Image)
	ContainerAssemble.ContainerWorkDir(broker.Spec.WorkDir)
	conatinerSlice := make([]corev1.ContainerPort, len(broker.Spec.ContainerPorts))
	counter := 0
	for k,v := range broker.Spec.ContainerPorts{
		conatinerSlice[counter] = templates.ContainerPortGenerator(k,v)
		counter++
	}
	ContainerAssemble.ContainerPort(conatinerSlice)
	ContainerAssemble.CommandWithArgs(broker.Spec.Commands,broker.Spec.Args)
	EnvMap:=map[string]string{"KAFKA_HEAP_OPTS": broker.Spec.Heap,
		"KAFKA_ZOOKEEPER_CONNECT": broker.Spec.ZooKeeperConnect,
		"KAFKA_LOG_DIRS": broker.Spec.LogDir,"KAFKA_METRIC_REPORTERS":broker.Spec.MetricReporters,
		"CONFLUENT_METRICS_REPORTER_BOOTSTRAP_SERVERS": broker.Spec.ReposterBootStrapServer,
	}
	ContainerAssemble.EnvVar(EnvMap)
	EnvSourceFieldMap:=map[string]string{"POD_IP": "status.podIP",
		"HOST_IP": "status.hostIP",
		"POD_NAME": "metadata.name",
		"POD_NAMESPACE": "metadata.namespace",
	}
	ContainerAssemble.EnvVarSourceFieldRef(EnvSourceFieldMap)
	VolumeMountsMap := make(map[string]string, broker.Spec.MountNum)
	for i:=0;i < broker.Spec.MountNum; i++{
		VolumeMountsMap["datadir-"+string(i)]= "/opt/kafka/data-"+string(i)
	}
	ContainerAssemble.VolumeMounts(VolumeMountsMap)
}

func pvcCompose(broker *kafkav1alpha1.BrokerOperator, PersistentVolumeClaimPTR *templates.PersistentVolumeClaimASSemble, i int){
	PersistentVolumeClaim:= *PersistentVolumeClaimPTR		
	PersistentVolumeClaim.Meta("datadir-"+string(i))
	PersistentVolumeClaim.AccessModes([]string{"ReadWriteOnce"})
	PersistentVolumeClaim.Resource(map[string]string{"storage":"1Gi"})
}