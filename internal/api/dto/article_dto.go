package dto

type CreateArticleRequest struct {
	Title string  `json:"title" validate:"required,min=3"`    // The title of the article, required and must be at least 3 characters long
	Text  string  `json:"content" validate:"required,min=10"` // The content of the article, required and must be at least 10 characters long
	Img   *string `json:"img"`                                // The image URL, can be null or a valid string
}
