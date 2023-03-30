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
	VoyageStatus AfterSaleStartVoyageResult `json:"voyage_status"`
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
