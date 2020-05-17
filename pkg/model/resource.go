package model

import "k8s.io/apimachinery/pkg/runtime/schema"

type Resource struct {
	Name             string
	Namespaced       bool
	Kind             string
	ShortName        []string
	Verbs            []string
	GroupVersionKind schema.GroupVersionKind
	GroupVersion     []schema.GroupVersion
}
