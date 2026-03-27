package apperrors

import "errors"

var (
	ErrSlugNotFound = errors.New("Slug not found")
	ErrUrlNotFound  = errors.New("Url not found")
)
