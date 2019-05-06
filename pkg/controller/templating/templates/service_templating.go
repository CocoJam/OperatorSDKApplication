package templates
import (
	corev1 "k8s.io/api/core/v1"
)

type Service struct{
	Meta DeploymentMetaTemplate
	S corev1.Service
}

func (s *Service) init(){
	s.S.TypeMeta = s.Meta.TypeMeta()
	s.S.ObjectMeta = s.Meta.ObjectMeta()
	s.Meta = s.Meta
}

func (s *Service) Selector(Selector map[string]string){
	s.S.Spec.Selector = Selector
}

func (s *Service) ServiceSpec(ClusterIP string){
	s.S.Spec.ClusterIP = ClusterIP
}

func (s *Service) ServicePort(Selector map[string]string){
	s.S.Spec.Selector = Selector
}

func (s *Service) ServiceType(ServiceType string){
	switch ServiceType {
	case string(corev1.ServiceTypeClusterIP):
		s.S.Spec.Type = corev1.ServiceTypeClusterIP
	case string(corev1.ServiceTypeNodePort):
		s.S.Spec.Type = corev1.ServiceTypeNodePort
	case string(corev1.ServiceTypeLoadBalancer):
		s.S.Spec.Type = corev1.ServiceTypeLoadBalancer
	case string(corev1.ServiceTypeExternalName):
		s.S.Spec.Type = corev1.ServiceTypeExternalName
	default:

	}
}

func(s *Service) BootStrap() corev1.Service{
	s.init()
	return s.S
} 

type ServicePort struct{}