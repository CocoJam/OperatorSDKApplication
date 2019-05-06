package zookeeper
import (
	"fmt"
	"strings"
	kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"github.com/example-inc/app-operator/pkg/controller/templating/templates"
)

type ZooKeeperStatefulSet struct{
	templateSS templates.StatefulSet
}

type DefaultZooKeeperStatefulSet struct{
	templateSS templates.StatefulSet
}

func (defaultBroker *DefaultZooKeeperStatefulSet) bootStrap(zookeeper *kafkav1alpha1.ZooKeeperOperator)  appsv1.StatefulSet{
	ls:= map[string]string{"app": "Zookeeper", "Zookeeper_cr": zookeeper.Name}
	meta:= templates.DeploymentMetaTemplate{
		Kind: "StatefulSet",
		APIVersion: "apps/v1beta1",
		ObjectName: zookeeper.Name,
		ObjectNamespace: zookeeper.Namespace,
	}
	defaultBroker.templateSS.Meta = meta

	defaultBroker.templateSS.Replicas(zookeeper.Spec.Replicas)
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
	containerCompose(zookeeper,&ContainerAssemble)
	defaultBroker.templateSS.SS.Spec.Template.Spec.Containers = []corev1.Container{ContainerAssemble.Container}

	PersistentVolumeClaim:= templates.PersistentVolumeClaimASSemble{}
	PersistentVolumeClaimSlice := make([]corev1.PersistentVolumeClaim, zookeeper.Spec.MountNum)
	for i:=0;i < zookeeper.Spec.MountNum; i++{
		pvcCompose(zookeeper, &PersistentVolumeClaim, i)
		PersistentVolumeClaimSlice[i] = PersistentVolumeClaim.PVC
	}
	defaultBroker.templateSS.SS.Spec.VolumeClaimTemplates = PersistentVolumeClaimSlice
	return defaultBroker.templateSS.BootStrap()
}

func containerCompose(zookeeper *kafkav1alpha1.ZooKeeperOperator,ContainerAss *templates.ContainerASSemble){
	ContainerAssemble := *ContainerAss
	ContainerAssemble.ImageFactory(zookeeper.Spec.ContainerName,zookeeper.Spec.Image)
	// ContainerAssemble.ContainerWorkDir(zookeeper.Spec.WorkDir)
	conatinerSlice := make([]corev1.ContainerPort, len(zookeeper.Spec.ContainerPorts))
	counter := 0
	for k,v := range zookeeper.Spec.ContainerPorts{
		conatinerSlice[counter] = templates.ContainerPortGenerator(k,v)
		counter++
	}
	ContainerAssemble.ContainerPort(conatinerSlice)
	ContainerAssemble.CommandWithArgs(zookeeper.Spec.Commands,zookeeper.Spec.Args)
	EnvMap:=map[string]string{"KAFKA_HEAP_OPTS": zookeeper.Spec.Heap,
		"ZOOKEEPER_CLIENT_PORT": serverlist(zookeeper, zookeeper.Spec.Replicas),
		"ZOOKEEPER_SERVERS": zookeeper.Spec.LogDir,
	}
	ContainerAssemble.EnvVar(EnvMap)
	EnvSourceFieldMap:=map[string]string{
		"ZOOKEEPER_SERVER_ID": "metadata.name",
	}
	ContainerAssemble.EnvVarSourceFieldRef(EnvSourceFieldMap)
	VolumeMountsMap := make(map[string]string, zookeeper.Spec.MountNum)
	for i:=0;i < zookeeper.Spec.MountNum; i++{
		VolumeMountsMap["datadir-"+string(i)]= "/var/lib/zookeeper/data-"+string(i)
		VolumeMountsMap["datalogdir-"+string(i)]= "/var/lib/zookeeper/log"+string(i)
	}
	ContainerAssemble.VolumeMounts(VolumeMountsMap)
}
func serverlist(zookeeper *kafkav1alpha1.ZooKeeperOperator, Replicas int) string{
	serversSlice := make([]string, Replicas)
	for i:=0;i<Replicas;i++{
		s := fmt.Sprintf("%s-%d.%s-headless.%s:%d:%d",
		 zookeeper.Name,
		 i, zookeeper.Name, zookeeper.Namespace,
		 zookeeper.Spec.ServerPort,
		 zookeeper.Spec.LeaderElectionPort,
		)
		serversSlice[i] = s
	}
	return strings.Join(serversSlice,";")
}


func pvcCompose(zookeeper *kafkav1alpha1.ZooKeeperOperator, PersistentVolumeClaimPTR *templates.PersistentVolumeClaimASSemble, i int){
	PersistentVolumeClaim:= *PersistentVolumeClaimPTR		
	PersistentVolumeClaim.Meta("datadir-"+string(i))
	PersistentVolumeClaim.AccessModes([]string{"ReadWriteOnce"})
	PersistentVolumeClaim.Resource(map[string]string{"storage":"1Gi"})
}