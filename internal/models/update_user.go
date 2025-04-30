package models

type UpdateUser struct {
	Firstname *string `json:"firstname,omitempty"`
	Lastname  *string `json:"lastname,omitempty"`
	Email     *string `json:"email,omitempty"`
	Age       *uint   `json:"age,omitempty"`
}
