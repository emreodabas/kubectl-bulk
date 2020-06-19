package model

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Resource struct {
	Name             string                  `json:"name"`
	Namespaced       bool                    `json:"namespaced"`
	Kind             string                  `json:"kind"`
	ShortName        []string                `json:"shortName"`
	Verbs            []string                `json:"verbs"`
	GroupVersionKind schema.GroupVersionKind `json:"gvk"`
	GroupVersion     []schema.GroupVersion   `json:"gv"`
}

type Resources struct {
	Resources []Resource `json:"resources"`
}
