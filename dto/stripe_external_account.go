package dto

type StripeExternalAccountInput struct {
	Country           string `binding:"required"`
	Currency          string `binding:"required"`
	RoutingNumber     string `binding:"required"`
	AccountNumber     string `binding:"required"`
	AccountHolderName string `binding:"required"`
	AccountHolderType string `binding:"required"`
}
