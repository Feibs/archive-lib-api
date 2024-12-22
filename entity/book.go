package entity

type Book struct {
	Id          int
	Author      *Author
	Title       string
	Description string
	Quantity    int
	Cover       *string
}

type BookPost struct {
	Id          int
	Title       string
	Description string
	Quantity    int
	Cover       *string
	AuthorId    int
}

func (bp BookPost) ConvertToBook() *Book {
	book := Book{
		Id:          bp.Id,
		Title:       bp.Title,
		Description: bp.Description,
		Quantity:    bp.Quantity,
	}

	if bp.Cover != nil {
		book.Cover = bp.Cover
	}

	return &book
}
