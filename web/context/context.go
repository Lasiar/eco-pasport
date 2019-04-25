package web

import (
	"context"
	"net/http"
)

type contextKey string

// ключи для ResponseContext
const (
	ResponseDataKey contextKey = "response.data"
)

// ResponseContext контекст для ответа
type ResponseContext struct {
	context.Context
	err error
	key interface{}
	val interface{}
}

// Err возврашает ошибку если не nil иначе ошибку родителя
func (rc *ResponseContext) Err() error {
	if rc.err != nil {
		return rc.err
	}
	return rc.Context.Err()
}

// Value возврашает значение контекста по ключу
// если в данном контексте нет то ищет у родителя
func (rc *ResponseContext) Value(key interface{}) interface{} {
	if key == rc.key {
		return rc.val
	}
	return rc.Context.Value(key)
}

// WithResponseContext возвращает копию родителя,
// если ошибка не nil то вернет ошибку родителя,
// возврашает поле data по ключу key
func WithResponseContext(parent context.Context, key, val interface{}, err error) context.Context {
	return &ResponseContext{
		Context: parent,
		err:     err,
		key:     key,
		val:     val,
	}
}

// SetError устанавливает в контекст ошибку
func SetError(r *http.Request, err error) {
	*r = *r.WithContext(
		WithResponseContext(
			r.Context(),
			ResponseDataKey,
			nil,
			err,
		),
	)
}

// SetResponse устанавливает в контекст json ответ
func SetResponse(r *http.Request, data interface{}) {
	*r = *r.WithContext(
		WithResponseContext(
			r.Context(),
			ResponseDataKey,
			data,
			nil,
		),
	)
}
