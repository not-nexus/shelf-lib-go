package shelflib

type ShelfError struct {
	code    string
	message string
}

func NewShelfError(message string, code string) *ShelfError {
	return &ShelfError{code, message}
}
