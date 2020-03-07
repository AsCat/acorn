package models

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/AsCat/acorn/kubernetes"
)

type ServiceRoleBindings []ServiceRoleBinding
type ServiceRoleBinding struct {
	meta_v1.TypeMeta
	Metadata meta_v1.ObjectMeta `json:"metadata"`
	Spec     struct {
		Subjects interface{} `json:"subjects"`
		RoleRef  interface{} `json:"roleRef"`
	} `json:"spec"`
}

func (srbs *ServiceRoleBindings) Parse(serviceRoleBindings []kubernetes.IstioObject) {
	for _, srb := range serviceRoleBindings {
		serviceRoleBinding := ServiceRoleBinding{}
		serviceRoleBinding.Parse(srb)
		*srbs = append(*srbs, serviceRoleBinding)
	}
}

func (srb *ServiceRoleBinding) Parse(serviceRoleBinding kubernetes.IstioObject) {
	srb.TypeMeta = serviceRoleBinding.GetTypeMeta()
	srb.Metadata = serviceRoleBinding.GetObjectMeta()
	srb.Spec.Subjects = serviceRoleBinding.GetSpec()["subjects"]
	srb.Spec.RoleRef = serviceRoleBinding.GetSpec()["roleRef"]
}
