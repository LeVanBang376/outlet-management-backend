package outlet

import "errors"

var (
	ErrNotFound = errors.New("outlet not found")

	// Validation errors
	ErrNameRequired    = errors.New("name is required")
	ErrAddressRequired = errors.New("address is required")
	ErrInvalidChannel  = errors.New("invalid channel")
	ErrInvalidTier     = errors.New("invalid tier")
	ErrInvalidStage    = errors.New("invalid stage")

	// Business errors
	ErrSalesNotFound = errors.New("sales not found")
)
