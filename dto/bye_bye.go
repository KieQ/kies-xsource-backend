package dto

type ByeByeCreateResult int64

const (
	ByeByeCreateResultSuccess ByeByeCreateResult = iota
	ByeByeCreateResultExist
	ByeByeCreateResultFail
)

type ByeByeCreateRequest struct {
}

type ByeByeCreateResponse struct {
	State    ByeByeCreateResult `json:"state"`
	ByeByeID string             `json:"bye_bye_id"`
}

type ByeByeFetchResultRequest struct {
	ByeByeID string `json:"bye_bye_id"`
}

type ByeByeFetchResultResponse struct {
	//TODO
}
