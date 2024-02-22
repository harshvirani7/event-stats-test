package model

type Data struct {
	Unique   string   `json:"unique"`
	Info     Info     `json:"info"`
	Workflow Workflow `json:"workflow"`
}

type Info struct {
}

type Workflow struct {
}
