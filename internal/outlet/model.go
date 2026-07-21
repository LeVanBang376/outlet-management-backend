package outlet

import "time"

type Outlet struct {
	OutletID  uint      `json:"outlet_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Channel   string    `json:"channel"`
	Tier      string    `json:"tier"`
	SalesID   uint      `json:"sales_id"`
	Stage     string    `json:"stage"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
