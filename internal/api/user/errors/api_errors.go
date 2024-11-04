package errors

import "fmt"

var (
	ErrUserModelIsNil      = fmt.Errorf("user model is nil")       // ErrUserModelIsNil модель юзера nil
	ErrDescUserIsNil       = fmt.Errorf("desc user is nil")        // ErrDescUserIsNil grpc запрос юзера nil
	ErrDescUserUpdateIsNil = fmt.Errorf("desc user update is nil") // ErrDescUserUpdateIsNil grpc запрос юзера для обновления nil
)
