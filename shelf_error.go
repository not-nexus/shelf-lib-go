package shelflib

type ShelfError struct {
	Code    string
    HasError bool
	Message string
    Parent error
}

func CreateShelfErrorFromError(parent error) *ShelfError {
    return &ShelfError{HasError: true, Message: parent.Error(), Parent: parent}
}

func CreateShelfError(message string, code string) *ShelfError {
    return &ShelfError{Code: code, HasError: true, Message: message}
}

func (this *ShelfError) Error() string {
	return this.Message
}
