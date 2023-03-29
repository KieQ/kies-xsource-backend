package dto

type AfterSaleStartVoyageResult int64

const (
	AfterSaleStartVoyageResultSuccess AfterSaleStartVoyageResult = iota
	AfterSaleStartVoyageResultExist
	AfterSaleStartVoyageResultFail
)

type AfterSaleStartVoyageRequest struct {
}

type AfterSaleCreateResponse struct {
	VoyageState AfterSaleStartVoyageResult `json:"voyage_state"`
}

type AfterSaleStartOverRequest struct {
}

type AfterSaleStartOverResponse struct {
}

type AfterSaleCheckResultRequest struct {
}

type AfterSaleCheckResultResponse struct {
	Pass bool `json:"pass"`
}

type AfterSaleNextStepRequest struct {
}

type AfterSaleNextStepResponse struct {
}

type AfterSaleFinalRewardRequest struct {
}

type AfterSaleFinalRewardResponse struct {
}
