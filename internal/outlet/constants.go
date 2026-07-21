package outlet

const (
	// Outlet stage
	StageRawLead          = "raw_lead"
	StageSQL              = "sql"
	StageCustomerSampling = "customer_sampling"
	StageProposalSent     = "proposal_sent"
	StageWon              = "won"
	StageLost             = "lost"

	// Outlet channels
	ChannelCafe       = "cafe"
	ChannelRestaurant = "restaurant"
	ChannelHotel      = "hotel"
	ChannelBar        = "bar"
	ChannelBakery     = "bakery"

	// Outlet tiers
	TierPremium = "premium"
	TierMid     = "mid"
	TierMass    = "mass"
)

var ValidChannels = map[string]bool{
	ChannelCafe:       true,
	ChannelRestaurant: true,
	ChannelHotel:      true,
	ChannelBar:        true,
	ChannelBakery:     true,
}

var ValidTiers = map[string]bool{
	TierPremium: true,
	TierMid:     true,
	TierMass:    true,
}

var ValidStages = map[string]bool{
	StageRawLead:          true,
	StageSQL:              true,
	StageCustomerSampling: true,
	StageProposalSent:     true,
	StageWon:              true,
	StageLost:             true,
}
