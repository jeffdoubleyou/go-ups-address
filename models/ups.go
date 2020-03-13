package models

type Response struct {
	ResponseStatus *ResponseStatus `json:"ResponseStatus"`
}

type ResponseStatus struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}
