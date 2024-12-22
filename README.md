# ArchiveLib REST API

An API for library services: add book, get book(s), borrow and return book with login authentication

## Description

1. As a librarian, I would like to see a complete list of books in the library.

- Each book should include details such as ID, title, description, quantity, and cover.

2. As a librarian, I would like to search a book by its title.

3. As a librarian, I would like to add a new book to the library collection.

- The added book must have a title, description, and quantity (other fields can remain empty).
- Duplicate titles are not allowed.
- The quantity must be 0 or higher.
- Titles cannot exceed 35 characters.

4. As a librarian, I would like to see the author’s detail when viewing the list of books.

- Author details should include their ID and name.
- Each book should reference one author using the author_id field.
- The response should be a list of books joined with their respective author data.
- When a book is successfully added, the response should include the book record (without author relation).

5. As a user, I would like to borrow a book.

- A borrowing record should be saved in a new table.
- If the book does not exist or out of stock, an error should be returned.
- Each request allows borrowing only one book.
- The book’s stock should decrease by 1 for every successful borrowing.

6. As a librarian, I would like the users to be able to return a book.

7. As a user, I would like to login so that I can borrow a book.

## Tech Stack

Go (Golang)

## How to Run

1. Clone this repository.
2. Make sure Go has been installed.
   Go version that is used in this app: `go1.19.13`
3. Setup the PostgreSQL database (see DDL queries in `schema.sql`).
4. Create `.env` file and adjust the variables accordingly (see `.env.example`).
5. Run the app: `go run .`
6. Run the unit tests: `go test ./...`

## Author

Feibs (2024)
