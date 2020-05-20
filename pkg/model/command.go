package model

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Command struct {
	Namespace     string
	Label         string
	FieldSelector string
	Action        Action
	Resource      Resource
	List          []unstructured.Unstructured
	Filter        Filter
}
