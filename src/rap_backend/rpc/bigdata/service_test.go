package bigdata

import (
	"fmt"
	"testing"
)

func TestGetCallInfo(t *testing.T) {
	resp, err := GetSvc().GetCallInfoList(GetCallInfoReq{CallIds: []string{"f7ec1e3db15eb544be4d7c319f004f39"}})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
