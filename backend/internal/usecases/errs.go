package usecases

import "errors"

var (
	ErrInvalidSortBy = errors.New("Invalid sort_by type")
	ErrInvalidYear	 = errors.New("Invalid year")
	ErrInvalidMonth	 = errors.New("Invalid month")
	ErrInvalidDate	 = errors.New("Invalid date")
	ErrRepository	 = errors.New("Repository error")
)