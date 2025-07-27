package repository

import (
	"fmt"
	serviceID "sample_project/internal/model"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
)

var new_srv []*service.Service
var new_res []*check.Result

func AddItem(item serviceID.SrvID) {
	switch v := item.(type) {
	case *service.Service:
		new_srv = append(new_srv, v)
		fmt.Printf(
			"Added service: %d %s %s %d %s\n",
			v.GetServiceID(),
			v.GetServiceName(),
			v.GetServiceURL(),
			v.GetServiceInterval(),
			v.ServiceCreatedAt(),
		)
	case *check.Result:
		new_res = append(new_res, v)
		fmt.Printf(
			"New check result: %d %d %d %s %d\n",
			v.GetResultID(),
			v.GetServiceID(),
			v.GetResponseCode(),
			v.GetTimeChecked(),
			v.GetResponseDuration(),
		)
	default:
		fmt.Printf("Unknown type: %T", item)
	}
}
