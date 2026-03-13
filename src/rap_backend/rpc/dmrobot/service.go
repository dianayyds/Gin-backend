package dmrobot

import (
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"net/http"
	"rap_backend/config"
	"rap_backend/utils"
)

type Service interface {
	GetRobotIntentionAnnotateList(req GetRobotIntentionAnnotateListReq) ([]*IntentionAnnotate, error)
}

type service struct {
}

func GetSvc() Service {
	return &service{}
}

func (s *service) GetRobotIntentionAnnotateList(req GetRobotIntentionAnnotateListReq) ([]*IntentionAnnotate, error) {
	if req.RobotName == "" {
		return nil, errors.New("empty robot name")
	}

	var transResp GetRobotIntentionAnnotateListTransResp
	err := doHttpGet(config.DMROBOT_HOST+"/api/robot/intention_annotate/list", &GetRobotIntentionAnnotateListTransReq{
		RobotName: req.RobotName,
	}, &transResp)
	if err != nil {
		seelog.Errorf("failed to request dmrobot interface, err: %+v", err)
		return nil, err
	}

	return transResp.List, nil
}

func doHttpPost(url string, request interface{}, response interface{}) error {
	commonResp := CommonResp{
		Body: response,
	}
	err := utils.DoHttpPostJson(url, request, &commonResp)
	if err != nil {
		seelog.Errorf("请求dmrobot失败, err: %+v", err)
		return err
	}

	if commonResp.Code != 200 {
		return fmt.Errorf("request dmrobot code is not 200, code: %d, message: %s", commonResp.Code, commonResp.Message)
	}

	return nil
}

func doHttpGet(url string, request GetParams, response interface{}) error {
	commonResp := CommonResp{
		Body: response,
	}

	url = fmt.Sprintf("%s?%s", url, request.ToParamsStr())
	err := utils.DoHttpGetJson(url, &commonResp, http.Header{
		"Authorization": []string{config.DMROBOT_TOKEN},
	})
	if err != nil {
		seelog.Errorf("请求dmrobot失败, err: %+v", err)
		return err
	}

	if commonResp.Code != 0 {
		return fmt.Errorf("request dmrobot code is not 0, code: %d, message: %s", commonResp.Code, commonResp.Message)
	}

	return nil
}

type CommonResp struct {
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}
