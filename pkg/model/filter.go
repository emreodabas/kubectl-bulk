package model

type Filters struct {
	Filters []Filter `json:"actions"`
}

type Filter struct {
	Name           string `json:"name"`
	Description    string `json:"desc"`
	Sample         string `json:"sample"`
	NeedFilterable bool   `json:"needFilter"`
}

var FilterList = []Filter{
	{
		Name:           "none",
		Description:    "no filter",
		Sample:         "--",
		NeedFilterable: false,
	},
	{
		Name:           "label",
		Description:    "filter with Label",
		Sample:         "kubectl get pods -l app=nginx",
		NeedFilterable: true,
	},
	{
		Name:           "field-selector",
		Description:    "filter with field selector",
		Sample:         " kubectl get pods --field-selector=\"status.phase!=Running\"",
		NeedFilterable: false,
	},
	{
		Name:           "multi-select",
		Description:    "Multi select from resource list",
		Sample:         "--",
		NeedFilterable: false,
	},
	{
		Name:           "grep",
		Description:    "using grep function from resource list",
		Sample:         "--",
		NeedFilterable: false,
	},
}
