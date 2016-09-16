package shelflib

type ShelfError struct {
	Code    string
	Message string
}

func NewShelfError(message string, code string) *ShelfError {
	return &ShelfError{Code: code, Message: message}
}

func (this *ShelfError) Error() string {
	return this.Message
}
