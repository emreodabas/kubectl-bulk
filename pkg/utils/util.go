package utils

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Unique(unstructureds *unstructured.UnstructuredList) []unstructured.Unstructured {
	keys := make(map[string]bool)
	list := []unstructured.Unstructured{}
	for _, entry := range unstructureds.Items {
		if _, value := keys[entry.GetName()]; !value {
			keys[entry.GetName()] = true
			list = append(list, entry)
		}
	}
	return list
}
