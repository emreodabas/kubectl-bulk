package model

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Command struct {
	Namespace     string
	Label         *v1.LabelSelector
	FieldSelector string
	Action        Action
	Resource      Resource
	List          []unstructured.Unstructured
	Filter        Filter
}
