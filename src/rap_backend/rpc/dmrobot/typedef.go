package dmrobot

import "fmt"

type GetParams interface {
	ToParamsStr() string
}

type IntentionAnnotate struct {
	Intention  string `json:"intention" form:"intention"`
	ZhAnnotate string `json:"zh_annotate" form:"zh_annotate"`
	IdAnnotate string `json:"id_annotate" form:"id_annotate"`
	EnAnnotate string `json:"en_annotate" form:"en_annotate"`
}

type GetRobotIntentionAnnotateListReq struct {
	RobotName string
}

type GetRobotIntentionAnnotateListResp struct {
	List []*IntentionAnnotate
}

type GetRobotIntentionAnnotateListTransResp struct {
	List []*IntentionAnnotate `json:"intention_annotate_list"`
}

type GetRobotIntentionAnnotateListTransReq struct {
	RobotName string `json:"robot_name"`
}

func (r *GetRobotIntentionAnnotateListTransReq) ToParamsStr() string {
	return fmt.Sprintf("robot_name=%s", r.RobotName)
}
