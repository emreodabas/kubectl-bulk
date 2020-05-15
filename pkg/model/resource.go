package model

type Resource struct {
	Name       string
	Namespaced bool
	Kind       string
	ShortName  []string
	Verbs      []string
}
