package internal

import (
	"fmt"
	"runtime"
)

// ErrorHelper - функция добавления имени файла и строки к ошибке
func ErrorHelper(err error) error {
	if err == nil {
		return err
	}
	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("[%s:%d] %w", filename, line, err)
}
