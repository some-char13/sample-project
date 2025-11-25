package check

import (
	"time"
)

type Result struct {
	Id           int       `json:"id" binding:"required"`
	ServiceId    int       `json:"service_id" binding:"required"`
	ResponseCode int       `json:"resp_code" binding:"required"`
	TimeChecked  time.Time `json:"time_checked"`
	RespDuration int       `json:"resp_duration" binding:"required"`
}

type ResultRequest struct {
	ServiceId    int       `json:"service_id" binding:"required"`
	ResponseCode int       `json:"resp_code" binding:"required"`
	RespDuration int       `json:"resp_duration" binding:"required"`
	TimeChecked  time.Time `json:"last_check"`
}

func NewResult(serviceId, responseCode, respDuration int) *Result {
	return &Result{
		ServiceId:    serviceId,
		ResponseCode: responseCode,
		RespDuration: respDuration,
		TimeChecked:  time.Now().UTC(),
	}
}
