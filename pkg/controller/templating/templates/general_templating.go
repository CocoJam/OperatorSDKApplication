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
func (pts *PodTemplateSpec) metaPodTemplateSpecs(meta DeploymentMetaTemplate ){
	pts.pts.ObjectMeta = meta.ObjectMeta()
}

type ContainerAssemble struct{
	Container corev1.Container
}

func (ass *ContainerAssemble) ImageFactory(name string, image string ){
	ass.Container.Name = name
	ass.Container.Image = image
}


func (ass *ContainerAssemble) ContainerWorkDir(workingDir string){
	ass.Container.WorkingDir = workingDir
}

func (ass *ContainerAssemble) ContainerPort(Ports[]corev1.ContainerPort){
	ass.Container.Ports = Ports
}

func (ass *ContainerAssemble) CommandWithArgs(command []string, args []string){
	ass.Container.Command = command
	ass.Container.Args = args
}

func (ass *ContainerAssemble) EnvVar(envMap map[string]string){
	envVarSlice := make([]corev1.EnvVar, len(envMap))
	counter :=0
	for k,v := range envMap{
		envVarSlice[counter] = corev1.EnvVar{Name: k, Value: v}
		counter++
	}
	ass.Container.Env = append(envVarSlice,  ass.Container.Env...)
}

func (ass *ContainerAssemble) EnvVarSourceFieldRef(envMap map[string]string){
	envVarSlice := make([]corev1.EnvVar, len(envMap))
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
	ass.Container.Env = append(envVarSlice,  ass.Container.Env...)
}

func (ass *ContainerAssemble) VolumeMounts(envMap map[string]string){
	VolumeSlice := make([]corev1.VolumeMount, len(envMap))
	counter :=0
	for k,v := range envMap{
		VolumeSlice[counter] = corev1.VolumeMount{Name: k, MountPath: v}
		counter++
	}
	ass.Container.VolumeMounts = append(VolumeSlice,  ass.Container.VolumeMounts...)
}

func (pts *PodTemplateSpec) metaPodTemplateSpecsContainter(Containers []corev1.Container){
	pts.pts.Spec.Containers = Containers
}


// PVC templating
func (pvc *PersistentVolumeClaim) metaPVC(meta DeploymentMetaTemplate){
	pvc.pvc.TypeMeta = meta.TypeMeta()
	pvc.pvc.ObjectMeta = meta.ObjectMeta()
}

func (pvc *PersistentVolumeClaim) AccessModes(AccessModes[]string){
	sliceAccess := make([]corev1.PersistentVolumeAccessMode, len(AccessModes))
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
	pvc.pvc.Spec.AccessModes = sliceAccess
}

func (pvc *PersistentVolumeClaim) Selector(selets map[string]string){
	pvc.pvc.Spec.Selector = &metav1.LabelSelector{MatchLabels: selets}
}

func (pvc *PersistentVolumeClaim) Resource(Resources map[string]string ){
	resourceMap := make(map[corev1.ResourceName]resource.Quantity, len(Resources))
	
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
	pvc.pvc.Spec.Resources.Requests = resourceMap
}

func (pvc *PersistentVolumeClaim) VolumeName(vol string){
	pvc.pvc.Spec.VolumeName = vol
}
func (pvc *PersistentVolumeClaim) StorageClassName(storage *string){
	pvc.pvc.Spec.StorageClassName = storage
}

// func (pvc *PersistentVolumeClaim) VolumeMode(mode string){
// 	switch mode {
// 	case string(corev1.PersistentVolumeBlock):
// 		pvc.pvc.Spec.VolumeMode = &corev1.PersistentVolumeFilesystem
// 	case string(corev1.PersistentVolumeFilesystem):
// 		pvc.pvc.Spec.VolumeMode = corev1.PersistentVolumeFilesystem
// 	default:

// 	}
// }

type PodTemplateSpec struct{
	pts corev1.PodTemplateSpec
}


type PersistentVolumeClaim struct{
	pvc corev1.PersistentVolumeClaim
}


type DeploymentSpecTemplate struct{
	Meta DeploymentMetaTemplate
	Replicas int32
	Selector *map[string]string
	Template corev1.PodTemplateSpec
	VolumeClaimTemplates []corev1.PersistentVolumeClaim
	ServiceName string
	PodManagementPolicy appsv1.PodManagementPolicyType
	UpdateStrategy 	appsv1.StatefulSetUpdateStrategy
}
