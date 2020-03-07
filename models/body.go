package models

type Request struct {
	Skip   int    `json:"skip"`
	Limit  int    `json:"limit"`
	Filter Filter `json:"filter"`
}

type Filter struct {
	Type         int    `json:"type"`
	Date         string `json:"date"`
	User         string `json:"user"`
	TargetIp     string `json:"targetIp"`
	Behavior     string `json:"behavior"`     // app
	BehaviorType string `json:"behaviorType"` // serv
	Url          string `json:"url"`
	Mac          string `json:"mac"`
	Port         string `json:"port"`
}

type Response struct {
	List []Action `json:"list"`
	Total int64 `json:"total"`
}
