package dto

type SendEmailRequest struct {
	Email string `json:"email" validate:"required,custom_email"`
	Name  string `json:"name" validate:"required,min=2"`
	Phone string `json:"phone" validate:"required,min=5"`
	Text  string `json:"text" validate:"required,min=5,max=500"`
}
