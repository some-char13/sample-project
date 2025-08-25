package check

import (
	"fmt"
	"time"
)

type Result struct {
	Id           int       `json:"id"`
	ServiceId    int       `json:"service_id"`
	ResponseCode int       `json:"response_code"`
	TimeChecked  time.Time `json:"time_checked"`
	RespDuration int       `json:"resp_duration"`
}

func NewResult(serviceId, responseCode, respDuration int) *Result {
	return &Result{
		ServiceId:    serviceId,
		ResponseCode: responseCode,
		TimeChecked:  time.Now().UTC(),
		RespDuration: respDuration,
	}
}

func (c *Result) String() string {
	return fmt.Sprintf("id: %d, serviceId: %d, responseCode: %d, timeChecked: %s, respDuration: %v", c.Id, c.ServiceId, c.ResponseCode, c.TimeChecked, c.RespDuration)
}

// func (r *Result) GetResultID() int {
// 	return r.id
// }
// func (r *Result) GetServiceID() int {
// 	return r.serviceId
// }
// func (r *Result) GetResponseCode() int {
// 	return r.responseCode
// }
// func (r *Result) GetResponseDuration() int {
// 	return r.respDuration
// }
// func (r *Result) GetTimeChecked() time.Time {
// 	return r.timeChecked
// }

// func (r *Result) SetID(id int) {
// 	r.id = id
// }
