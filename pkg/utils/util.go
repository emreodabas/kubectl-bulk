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

func Unique(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	keys := make(map[string]bool)
	list := []unstructured.Unstructured{}
	for _, entry := range unstructureds {
		if _, value := keys[entry.GetName()]; !value {
			keys[entry.GetName()] = true
			list = append(list, entry)
		}
	}
	return list
}

func Keys(list map[string]string) []string {
	var keys []string
	for k, _ := range list {
		keys = append(keys, k)
	}
	return keys
}

func RemoveStructure(s []unstructured.Unstructured, i int) []unstructured.Unstructured {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveItem(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
