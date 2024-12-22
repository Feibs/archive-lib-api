package dto

type BookResponse struct {
	Id          int             `json:"id"`
	Author      *AuthorResponse `json:"author,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Quantity    int             `json:"quantity"`
	Cover       string          `json:"cover,omitempty"`
}

type BookRequest struct {
	Title       string  `json:"title" binding:"required,max=35"`
	AuthorId    *int    `json:"author_id" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	Quantity    *int    `json:"quantity" binding:"required,gte=0"`
	Cover       *string `json:"cover"`
}
