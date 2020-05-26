package model

type Actions struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Name         string `json:"name"`
	Descriptions string `json:"desc"`
	Sample       string `json:"sample"`
}

var Actionlist = []Action{
	Action{Name: "GET", Descriptions: "Getting bulk information of resources", Sample: "kubectl-bulk get po images "},
	Action{Name: "UPDATE", Descriptions: "Bulk updating of resources", Sample: "kubectl-bulk update po images "},
	Action{Name: "DELETE", Descriptions: "Deleting of resources", Sample: "kubectl-bulk get delete images "},
}
