package dto

import "magnolia-test-backend/internal/model"

type SalesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewSalesResponse(s model.Sales) SalesResponse {
	return SalesResponse{
		ID:   s.SalesID,
		Name: s.Name,
	}
}
