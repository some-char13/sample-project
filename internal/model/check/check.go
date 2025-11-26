package check

import (
	"time"
)

type Result struct {
	ID           int       `json:"id" binding:"required"`
	ServiceID    int       `json:"serviceId" binding:"required"`
	ResponseCode int       `json:"respCode" binding:"required"`
	TimeChecked  time.Time `json:"timeChecked"`
	RespDuration int       `json:"respDuration" binding:"required"`
}

type ResultRequest struct {
	ServiceID    int       `json:"serviceId" binding:"required"`
	ResponseCode int       `json:"respCode" binding:"required"`
	RespDuration int       `json:"respDuration" binding:"required"`
	TimeChecked  time.Time `json:"lastCheck"`
}

func NewResult(serviceID, responseCode, respDuration int) *Result {
	return &Result{
		ServiceID:    serviceID,
		ResponseCode: responseCode,
		RespDuration: respDuration,
		TimeChecked:  time.Now().UTC(),
	}
}
