package model

type NewOrUpdateMenuFunctionParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Active      bool   `json:"active"`
}
