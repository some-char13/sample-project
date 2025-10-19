package check

import (
	"fmt"
	"time"
)

// type Result struct {
// 	Id           int       `json:"id"`
// 	ServiceId    int       `json:"service_id"`
// 	ResponseCode int       `json:"response_code"`
// 	TimeChecked  time.Time `json:"time_checked"`
// 	RespDuration int       `json:"resp_duration"`
// }

type Result struct {
	Id           int       `json:"id" binding:"required"`
	ServiceId    int       `json:"service_id" binding:"required"`
	ResponseCode int       `json:"resp_code" binding:"required"`
	TimeChecked  time.Time `json:"time_cheched"`
	RespDuration int       `json:"resp_duration" binding:"required"`
}

func NewResult(id, serviceId, responseCode, respDuration int) *Result {
	return &Result{
		Id:           id,
		ServiceId:    serviceId,
		ResponseCode: responseCode,
		TimeChecked:  time.Now().UTC(),
		RespDuration: respDuration,
	}
}

func (c *Result) String() string {
	return fmt.Sprintf("id: %d, serviceId: %d, responseCode: %d, timeChecked: %s, respDuration: %v", c.Id, c.ServiceId, c.ResponseCode, c.TimeChecked, c.RespDuration)
}

func (c *Result) UnformString() string {
	if c == nil {
		return "nil"
	}
	// return fmt.Sprintf("%d %s %s %d %s",
	// 	s.Id, s.Name, s.Url, s.Interval, s.Created)
	// return fmt.Sprint(
	// 	s.Id, s.Name, s.Url, s.Interval, s.Created)
	return fmt.Sprintf("%d,%d,%d,%s,%d",
		c.Id, c.ServiceId, c.ResponseCode, c.TimeChecked, c.RespDuration)
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
