package bigdata

type CallInfo struct {
	BillSec      int      `json:"BillSec"`
	CalleeNumber string   `json:"CalleeNumber"`
	Company      string   `json:"Company"`
	Intention    string   `json:"Intention"`
	RobotName    string   `json:"RobotName"`
	RobotType    string   `json:"RobotType"`
	SIPLine      string   `json:"SIPLine"`
	StartTime    string   `json:"StartTime"`
	TalkRound    int      `json:"TalkRound"`
	SenTence     []string `json:"SenTence"`
	CallId       string   `json:"callid"`
}

type GetCallInfoReq struct {
	CallIds []string
}

type GetCallInfoResp struct {
	Calls []*CallInfo
}

type GetCallInfoTransResp []*CallInfo

type GetCallInfoTransReq struct {
	Call_id_list string `json:"Call_id_list"`
}
