package model

type AddNewContactMessageParams struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type UpdateContactMessageParams struct {
	Status int16 `json:"status"`
}

type ContactMessageQuery struct {
	Limit  int   `json:"limit"`
	Page   int   `json:"page"`
	Status int16 `json:"status"`
}
