package dto

type ImageDocumentInput struct {
	Reference *string `json:"reference" binding:"required"`

	// Reference *string `json:"reference" binding:"required"`
}
