package dto

type MotoristRequestSituationCreationInput struct {
	Code        string `json:"code" binding:""`
	Title       string `json:"title" binding:""`
	SubTitlte   string `json:"subtitle" binding:""`
	Note        string `json:"note" binding:""`
	Description string `json:"description" binding:""`
}

type MotoristRequestSituationUpdateInput struct {
	Code        *string `json:"code" binding:""`
	Title       *string `json:"title" binding:""`
	SubTitlte   *string `json:"subtitle" binding:""`
	Note        *string `json:"note" binding:""`
	Description *string `json:"description" binding:""`
}
