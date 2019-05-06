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
	// "k8s.io/apimachinery/pkg/api/resource"
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

type StatefulSet struct{
	meta DeploymentMetaTemplate
	SS appsv1.StatefulSet
	SpecTemplate DeploymentSpecTemplate
}

func (ks *StatefulSet) init(){
	ks.SS.TypeMeta = ks.meta.TypeMeta()
	ks.SS.ObjectMeta = ks.meta.ObjectMeta()
	ks.SpecTemplate.Meta = ks.meta
}

func (ks *StatefulSet) Replicas(num int){
	replicas  := int32(num)
	ks.SS.Spec.Replicas = &replicas
	ks.SpecTemplate.Replicas = num
}

func (ks *StatefulSet) SpecSelector(ls map[string]string){
	ks.SS.Spec.Selector =  &metav1.LabelSelector{
		MatchLabels: ls,
	}
	ks.SpecTemplate.Selector = &ls
}

func (ks *StatefulSet) PodTemplateSpecObjectMeta(){
	ks.SS.Spec.Template.ObjectMeta =  metav1.ObjectMeta{
		Labels: *ks.SpecTemplate.Selector,
	}
}

func (ks *StatefulSet) PodTemplateSpecSpec(containerASS ContainerASSemble){
	ks.SS.Spec.Template.Spec.Containers = append(ks.SS.Spec.Template.Spec.Containers, containerASS.Container)
}

func(ks *StatefulSet) VolumeClaimTemplates(pvc corev1.PersistentVolumeClaim){
	ks.SS.Spec.VolumeClaimTemplates = append( ks.SS.Spec.VolumeClaimTemplates, pvc)
}

func(ks *StatefulSet) bootStrap() appsv1.StatefulSet{
	ks.init()
	return ks.SS
} 
