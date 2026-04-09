package handler

type createFlagRequest struct {
	ProjectID   string `json:"project_id"   example:"77c00606-0099-4642-83e4-0d03c6f78c36"`
	Key         string `json:"key"          example:"checkout-v2"`
	Name        string `json:"name"         example:"Checkout V2"`
	Description string `json:"description"  example:"New checkout flow"`
	Type        string `json:"type"         example:"boolean"`
}

type updateFlagRequest struct {
	Key         string `json:"key"          example:"checkout-v2"`
	Name        string `json:"name"         example:"Checkout V2"`
	Description string `json:"description"  example:"Updated description"`
	Type        string `json:"type"         example:"boolean"`
}
