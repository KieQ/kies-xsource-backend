package dto

type AfterSaleVoyageCheckProgressResult int64

const (
	AfterSaleVoyageCheckProgressResultNeverStarted AfterSaleVoyageCheckProgressResult = iota
	AfterSaleVoyageCheckProgressResultInTrip
	AfterSaleVoyageCheckProgressResultCleared
)

type AfterSaleVoyageCheckProgressResponse struct {
	Progress AfterSaleVoyageCheckProgressResult `json:"progress"`
	Level    int8                         `json:"level"`
}

type AfterSaleVoyageStartOrContinueTripRequest struct {
	Level int8 `json:"level"`
}

type AfterSaleVoyageStartOrContinueTripResponse struct {
	Level int8 `json:"level"`
	Cleared bool `json:"cleared"`
	Seed int32 `json:"seed"`
}

type AfterSaleVoyageStartOverRequest struct {
	Level int8 `json:"level"`
}

type AfterSaleVoyageStartOverResponse struct {
	Level int8 `json:"level"`
	Seed int32 `json:"seed"`
}

type AfterSaleVoyageCheckResultRequest struct {
	Level int8 `json:"level"`
	IV *string `json:"iv"`
	Result *string `json:"result"`
}

type AfterSaleVoyageCheckResultResponse struct {
	Pass bool `json:"pass"`
	FailReason string `json:"fail_reason"`
}

type AfterSaleVoyageNextStepRequest struct {
	FinishedLevel int8 `json:"finished_level"`
}

type AfterSaleVoyageNextStepResponse struct {
	Level int8 `json:"level"`
	Cleared bool `json:"cleared"`
	Seed int32 `json:"seed"`
}

type AfterSaleVoyageFinalRewardRequest struct {
}

type AfterSaleVoyageFinalRewardResponse struct {
}
