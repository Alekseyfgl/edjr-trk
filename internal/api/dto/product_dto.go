package dto

type CreateProductRequest struct {
	Title     string  `json:"title" validate:"required,min=3"`
	Text      string  `json:"text" validate:"required,min=10,max=5000"`
	ShortText string  `json:"shortText" validate:"required,min=10,max=1000"`
	Img       *string `json:"img" validate:"omitempty,img_base64_or_null"`
}

type PatchProductRequest struct {
	Title     *string `json:"title" validate:"omitempty,min=3"`
	Text      *string `json:"text" validate:"omitempty,min=10,max=5000"`
	ShortText *string `json:"shortText" validate:"required,min=10,max=1000"`
	Img       *string `json:"img" validate:"omitempty,img_base64_or_null"`
}
