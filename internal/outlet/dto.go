package outlet

type CreateOutletRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Channel string `json:"channel"`
	Tier    string `json:"tier"`
	SalesID uint   `json:"sales_id"`
	Stage   string `json:"stage"`
	Note    string `json:"note"`
}

func (r CreateOutletRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}

	if r.Address == "" {
		return ErrAddressRequired
	}

	if !ValidChannels[r.Channel] {
		return ErrInvalidChannel
	}

	if !ValidTiers[r.Tier] {
		return ErrInvalidTier
	}

	if !ValidStages[r.Stage] {
		return ErrInvalidStage
	}

	return nil
}

type UpdateOutletRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Channel string `json:"channel"`
	Tier    string `json:"tier"`
	SalesID uint   `json:"sales_id"`
	Stage   string `json:"stage"`
	Note    string `json:"note"`
}

func (r UpdateOutletRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}

	if r.Address == "" {
		return ErrAddressRequired
	}

	if !ValidChannels[r.Channel] {
		return ErrInvalidChannel
	}

	if !ValidTiers[r.Tier] {
		return ErrInvalidTier
	}

	if !ValidStages[r.Stage] {
		return ErrInvalidStage
	}

	return nil
}

type OutletResponse struct {
	OutletID uint   `json:"outlet_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Channel  string `json:"channel"`
	Tier     string `json:"tier"`
	SalesID  uint   `json:"sales_id"`
	Stage    string `json:"stage"`
	Note     string `json:"note"`
}
