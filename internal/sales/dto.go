package sales

type SalesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewSalesResponse(s Sales) SalesResponse {
	return SalesResponse{
		ID:   s.SalesID,
		Name: s.Name,
	}
}
