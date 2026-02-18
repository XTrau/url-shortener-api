package errors

import "errors"

var (
	SlugNotFound = errors.New("Slug not found")
	UrlNotFound  = errors.New("Url not found")
)
