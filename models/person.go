package models

type PersonRequest struct {
	Name    string `json:"name,omitempty" example:"John"`
	Surname string `json:"surname,omitempty" example:"Doe"`
} // @name Person

type PersonResponse struct {
	PersonRequest
	ID string `json:"id,omitempty"`
} // @name PersonResponse
