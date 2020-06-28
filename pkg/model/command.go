package model

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Command struct {
	Namespace       string
	LabelFilter     string
	FieldSelector   string
	GrepFilter      []string
	Action          Action
	Resource        Resource
	List            []unstructured.Unstructured
	Filter          Filter
	SelectedFilters []Filter
	UpdatedList     []unstructured.Unstructured
	SelectedSpec    []string
}
