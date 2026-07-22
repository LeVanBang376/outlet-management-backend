package customerrors

import "errors"

var (
	OutletErrNotFound = errors.New("outlet not found")

	// Validation errors
	OutletErrNameRequired        = errors.New("name is required")
	OutletErrAddressRequired     = errors.New("address is required")
	OutletErrInvalidChannel      = errors.New("invalid channel")
	OutletErrInvalidTier         = errors.New("invalid tier")
	OutletErrInvalidStage        = errors.New("invalid stage")
	OutletErrInvalidScheduleDate = errors.New("invalid schedule date")

	// Business errors
	OutletErrSalesNotFound = errors.New("sales not found")
)
