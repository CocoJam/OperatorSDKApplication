package template
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
	"github.com/example-inc/app-operator/pkg/controller/templating/templates"
)

type kafkaStatefulSet struct{
	ss appsv1.StatefulSet
	SpecTemplate templates.DeploymentSpecTemplate
}

func (ks *kafkaStatefulSet) init(meta templates.DeploymentMetaTemplate){
	ks.ss.TypeMeta = meta.TypeMeta()
	ks.ss.ObjectMeta = meta.ObjectMeta()
	ks.SpecTemplate.Meta = meta
}

func (ks *kafkaStatefulSet) Replicas(num int32){
	ks.ss.Spec.Replicas = &num
	ks.SpecTemplate.Replicas = num
}

func (ks *kafkaStatefulSet) SpecSelector(ls map[string]string){
	ks.ss.Spec.Selector =  &metav1.LabelSelector{
		MatchLabels: ls,
	}
	ks.SpecTemplate.Selector = &ls
}

func (ks *kafkaStatefulSet) PodTemplateSpecObjectMeta(){
	ks.ss.Spec.Template.ObjectMeta =  metav1.ObjectMeta{
		Labels: *ks.SpecTemplate.Selector,
	}
}

func (ks *kafkaStatefulSet) PodTemplateSpecSpec(containerAss templates.ContainerAssemble){
	ks.ss.Spec.Template.Spec.Containers = append(ks.ss.Spec.Template.Spec.Containers, containerAss.Container)
}

func(ks *kafkaStatefulSet) VolumeClaimTemplates(pvc corev1.PersistentVolumeClaim){
	ks.ss.Spec.VolumeClaimTemplates = append( ks.ss.Spec.VolumeClaimTemplates, pvc)
}