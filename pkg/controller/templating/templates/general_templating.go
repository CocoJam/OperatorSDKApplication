package templates
import (
	// "reflect"
	// "context"
	// er "errors"
	// kafkav1alpha1 "github.com/example-inc/app-operator/pkg/apis/kafka/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	// "sigs.k8s.io/controller-runtime/pkg/handler"
	// "sigs.k8s.io/controller-runtime/pkg/manager"
	// "sigs.k8s.io/controller-runtime/pkg/reconcile"
	// logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	// "sigs.k8s.io/controller-runtime/pkg/source"
)

// Meta labels and Tags
type MetaDeployment interface{
	TypeMeta() metav1.TypeMeta
	ObjectMeta() metav1.ObjectMeta
}

type DeploymentMetaTemplate struct{
	Kind string
	APIVersion string
	ObjectName string
	ObjectNamespace string
}

func (dmt *DeploymentMetaTemplate) TypeMeta() metav1.TypeMeta{
	return metav1.TypeMeta{
		Kind: dmt.Kind,
		APIVersion: dmt.APIVersion,
	}
}

func (dmt *DeploymentMetaTemplate) ObjectMeta() metav1.ObjectMeta{
	return metav1.ObjectMeta{
		Name: dmt.ObjectNamespace,
		Namespace: dmt.ObjectNamespace,
	}
}

//Selector
type Selector interface{
	Selector() *metav1.LabelSelector
}

//
// type SpecDeployment interface{
// 	Selector() *metav1.LabelSelector
// 	PersistentVolumeClaim() corev1.PersistentVolumeClaim
// }

//PodTemplateSpec
func (PTS *PodTemplateSpec) MetaPodTemplateSpecs(ls map[string]string ){
	PTS.PTS.ObjectMeta = metav1.ObjectMeta{
		Labels: ls,
	}
}

// func (PTS *PodTemplateSpec) MetaPodTemplateSpecs(meta DeploymentMetaTemplate ){
// 	PTS.PTS.ObjectMeta = meta.ObjectMeta()
// }

type ContainerASSemble struct{
	Container corev1.Container
}

func (ASS *ContainerASSemble) ImageFactory(name string, image string ){
	ASS.Container.Name = name
	ASS.Container.Image = image
}


func (ASS *ContainerASSemble) ContainerWorkDir(workingDir string){
	ASS.Container.WorkingDir = workingDir
}

func ContainerPortGenerator(name string, ContainerPort int) corev1.ContainerPort{
	container := int32(ContainerPort)
	return corev1.ContainerPort{Name: name, ContainerPort: container}
}

func (ASS *ContainerASSemble) ContainerPort(Ports[]corev1.ContainerPort){
	ASS.Container.Ports = Ports
}

func (ASS *ContainerASSemble) CommandWithArgs(command []string, args []string){
	ASS.Container.Command = command
	ASS.Container.Args = args
}

func (ASS *ContainerASSemble) EnvVar(envMap map[string]string){
	envVarSlice := make([]corev1.EnvVar, len(envMap))
	if(len(envVarSlice)<1){
		return
	}
	counter :=0
	for k,v := range envMap{
		envVarSlice[counter] = corev1.EnvVar{Name: k, Value: v}
		counter++
	}
	ASS.Container.Env = append(envVarSlice,  ASS.Container.Env...)
}

func (ASS *ContainerASSemble) EnvVarSourceFieldRef(envMap map[string]string){
	envVarSlice := make([]corev1.EnvVar, len(envMap))
	if(len(envVarSlice)<1){
		return
	}
	counter :=0
	for k,v := range envMap{
		envVarSlice[counter] = corev1.EnvVar{
			Name: k,
			ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				// APIVersion: k,
				FieldPath: v},
			}}
		counter++
	}
	ASS.Container.Env = append(envVarSlice,  ASS.Container.Env...)
}

func (ASS *ContainerASSemble) VolumeMounts(envMap map[string]string){
	VolumeSlice := make([]corev1.VolumeMount, len(envMap))
	if(len(VolumeSlice)<1){
		return
	}
	counter :=0
	for k,v := range envMap{
		VolumeSlice[counter] = corev1.VolumeMount{Name: k, MountPath: v}
		counter++
	}
	ASS.Container.VolumeMounts = append(VolumeSlice,  ASS.Container.VolumeMounts...)
}

func (PTS *PodTemplateSpec) metaPodTemplateSpecsContainter(Containers []corev1.Container){
	PTS.PTS.Spec.Containers = Containers
}


// PVC templating
func (PVC *PersistentVolumeClaimASSemble) Meta(metaName string){
	PVC.PVC.ObjectMeta = metav1.ObjectMeta{
		Name: metaName,
	}
}

func (PVC *PersistentVolumeClaimASSemble) metaPVC(meta DeploymentMetaTemplate){
	PVC.PVC.TypeMeta = meta.TypeMeta()
	PVC.PVC.ObjectMeta = meta.ObjectMeta()
}

func (PVC *PersistentVolumeClaimASSemble) AccessModes(AccessModes[]string){
	sliceAccess := make([]corev1.PersistentVolumeAccessMode, len(AccessModes))
	if(len(sliceAccess)<1){
		return
	}
	for i,access := range AccessModes{
		switch access {
		case string(corev1.ReadWriteOnce):
			sliceAccess[i] = corev1.ReadWriteOnce
		case string(corev1.ReadOnlyMany):
			sliceAccess[i] = corev1.ReadOnlyMany
		case string(corev1.ReadWriteMany):
			sliceAccess[i] = corev1.ReadWriteMany
		default:
		}
	}
	PVC.PVC.Spec.AccessModes = sliceAccess
}

func (PVC *PersistentVolumeClaimASSemble) Selector(selets map[string]string){
	PVC.PVC.Spec.Selector = &metav1.LabelSelector{MatchLabels: selets}
}

func (PVC *PersistentVolumeClaimASSemble) Resource(Resources map[string]string ){
	resourceMap := make(map[corev1.ResourceName]resource.Quantity, len(Resources))
	if(len(resourceMap)<1){
		return
	}
	for k,_ := range Resources{
		switch k {
		case string(corev1.ResourceCPU):
			resourceMap[corev1.ResourceCPU] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceMemory):
			resourceMap[corev1.ResourceMemory] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceStorage):
			resourceMap[corev1.ResourceStorage] = resource.Quantity{Format: resource.DecimalSI}
		case string(corev1.ResourceEphemeralStorage):
			resourceMap[corev1.ResourceEphemeralStorage] = resource.Quantity{Format: resource.DecimalSI}
		default:

		}
	}
	PVC.PVC.Spec.Resources.Requests = resourceMap
}

func (PVC *PersistentVolumeClaimASSemble) VolumeName(vol string){
	PVC.PVC.Spec.VolumeName = vol
}
func (PVC *PersistentVolumeClaimASSemble) StorageClASSName(storage *string){
	PVC.PVC.Spec.StorageClassName = storage
}

// func (PVC *PersistentVolumeClaimASSemble) VolumeMode(mode string){
// 	switch mode {
// 	case string(corev1.PersistentVolumeBlock):
// 		PVC.PVC.Spec.VolumeMode = &corev1.PersistentVolumeFilesystem
// 	case string(corev1.PersistentVolumeFilesystem):
// 		PVC.PVC.Spec.VolumeMode = corev1.PersistentVolumeFilesystem
// 	default:

// 	}
// }

//PodManagementPolicy

func(pmp PodManagementPolicy) PodManagementPolicy(pmpType string){
	switch pmpType{
	case string(appsv1.OrderedReadyPodManagement):
		pmp.PMP = appsv1.OrderedReadyPodManagement
	case string(appsv1.ParallelPodManagement):
		pmp.PMP = appsv1.ParallelPodManagement
	}
}

func(US UpdateStrategy) UpdateStrategy(Type string, RollingUpdate int){
	switch Type{
	case string(appsv1.RollingUpdateStatefulSetStrategyType):
		US.US.Type = appsv1.RollingUpdateStatefulSetStrategyType
	case string(appsv1.OnDeleteStatefulSetStrategyType):
		US.US.Type = appsv1.OnDeleteStatefulSetStrategyType
	}
	switch {
	case RollingUpdate>0:
		num  := int32(RollingUpdate)
		US.US.RollingUpdate.Partition = &num
	default:
	}
}

type PodTemplateSpec struct{
	PTS corev1.PodTemplateSpec
}

type PersistentVolumeClaimASSemble struct{
	PVC corev1.PersistentVolumeClaim
}

type PodManagementPolicy struct{
	PMP appsv1.PodManagementPolicyType
}

type UpdateStrategy struct{
	US appsv1.StatefulSetUpdateStrategy
}

type DeploymentSpecTemplate struct{
	Meta DeploymentMetaTemplate
	Replicas int
	Selector *map[string]string
	Template corev1.PodTemplateSpec
	VolumeClaimTemplates []corev1.PersistentVolumeClaim
	ServiceName string
	PodManagementPolicy appsv1.PodManagementPolicyType
	UpdateStrategy 	appsv1.StatefulSetUpdateStrategy
}
