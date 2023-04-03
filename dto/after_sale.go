package dto

type AfterSaleVoyageCheckProgressResult int64

const (
	AfterSaleVoyageCheckProgressResultNeverStarted AfterSaleVoyageCheckProgressResult = iota
	AfterSaleVoyageCheckProgressResultInTrip
	AfterSaleVoyageCheckProgressResultPass
)

type AfterSaleVoyageCheckProgressResponse struct {
	UserID   int32                        `json:"user_id"`
	Progress AfterSaleVoyageCheckProgressResult `json:"progress"`
	Level    int8                         `json:"level"`
}

type AfterSaleVoyageStartOrContinueTripRequest struct {
	Level int8 `json:"level"`
}

type AfterSaleVoyageStartOrContinueTripResponse struct {
	Level int8 `json:"level"`
	Passed bool `json:"passed"`
}

type AfterSaleVoyageStartOverRequest struct {
	Level int8 `json:"level"`
}

type AfterSaleVoyageStartOverResponse struct {
	Level int8 `json:"level"`
}

type AfterSaleVoyageCheckResultRequest struct {
}

type AfterSaleVoyageCheckResultResponse struct {
	Pass bool `json:"pass"`
}

type AfterSaleVoyageNextStepRequest struct {
}

type AfterSaleVoyageNextStepResponse struct {
}

type AfterSaleVoyageFinalRewardRequest struct {
}

type AfterSaleVoyageFinalRewardResponse struct {
}
