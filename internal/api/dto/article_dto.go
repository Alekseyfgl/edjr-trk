package dto

type CreateArticleRequest struct {
	Title string  `json:"title" validate:"required,min=3"`             // The title of the article, required and must be at least 3 characters long
	Text  string  `json:"text" validate:"required,min=10"`             // The content of the article, required and must be at least 10 characters long
	Img   *string `json:"img" validate:"omitempty,img_base64_or_null"` // The image URL, optional, can be null or a valid Base64 string
}
