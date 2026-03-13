package dmrobot

import (
	"fmt"
	"testing"
)

func TestGetCallInfo(t *testing.T) {
	resp, err := GetSvc().GetRobotIntentionAnnotateList(GetRobotIntentionAnnotateListReq{
		RobotName: "caolulu_online_robot",
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
