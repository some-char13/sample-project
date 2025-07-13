package check

import "time"

type Result struct {
	id           int
	serviceId    int
	responseCode int
	timeChecked  time.Time
	respDuration int
}

func NewResult(serviceId, responseCode, respDuration int) *Result {
	return &Result{
		serviceId:    serviceId,
		responseCode: responseCode,
		timeChecked:  time.Now().UTC(),
		respDuration: respDuration,
	}
}

func (r *Result) GetID() int {
	return r.id
}
func (r *Result) GetServiceId() int {
	return r.serviceId
}
func (r *Result) GetResponseCode() int {
	return r.responseCode
}
func (r *Result) GetResponseDuration() int {
	return r.respDuration
}
func (r *Result) GetTimeChecked() time.Time {
	return r.timeChecked
}

func (r *Result) SetID(id int) {
	r.id = id
}
