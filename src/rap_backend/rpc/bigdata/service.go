package bigdata

import (
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"rap_backend/config"
	"rap_backend/utils"
	"sync"
)

type Service interface {
	GetCallInfoList(req GetCallInfoReq) ([]*CallInfo, error)
	GetCallInfoMap(req GetCallInfoReq) (map[string]*CallInfo, error)
	GetCallInfo(req GetCallInfoReq) (*CallInfo, error)
}

type service struct {
}

func GetSvc() Service {
	return &service{}
}

func (s *service) GetCallInfo(req GetCallInfoReq) (*CallInfo, error) {
	list, err := s.GetCallInfoList(req)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("call info not found")
	}
	return list[0], nil
}

func (s *service) GetCallInfoMap(req GetCallInfoReq) (m map[string]*CallInfo, err error) {
	callInfoList, err := s.GetCallInfoList(req)
	if err != nil {
		return nil, err
	}

	m = make(map[string]*CallInfo)

	for _, callInfo := range callInfoList {
		m[callInfo.CallId] = callInfo
	}

	return m, nil
}

func (s *service) GetCallInfoList(req GetCallInfoReq) ([]*CallInfo, error) {
	if len(req.CallIds) == 0 {
		return nil, errors.New("empty callid")
	}

	var wg sync.WaitGroup

	ret := make([]*CallInfo, 0)
	var requestErr error

	var lock sync.Mutex

	//大数据接口最大只能支持99 callId查询
	for i := 0; i < len(req.CallIds); i += 200 {
		for j := i; j < i+200; j++ {
			if j == len(req.CallIds) {
				break
			}
			f := func(index int) {
				defer wg.Done()

				callIdStr := req.CallIds[index]
				//				callIdStr := strings.Join(req.CallIds[startI:endI], ",")
				var transResp []*CallInfo
				err := doHttpPost(config.BIGDATA_HOST+"/api/rap/callidinfo", GetCallInfoTransReq{
					Call_id_list: callIdStr,
				}, &transResp)
				if err != nil {
					seelog.Errorf("failed to request bigdata interface, callIds: %s, err: %+v", callIdStr, err)
					requestErr = err
					return
				}

				lock.Lock()
				ret = append(ret, transResp...)
				lock.Unlock()
			}
			go f(j)
			wg.Add(1)
		}

		wg.Wait()

		if requestErr != nil {
			return nil, requestErr
		}
	}

	return ret, nil
}

func doHttpPost(url string, request interface{}, response interface{}) error {
	commonResp := CommonResp{
		Body: response,
	}
	err := utils.DoHttpPostJson(url, request, &commonResp)
	if err != nil {
		seelog.Errorf("请求bigdata失败, err: %+v", err)
		return err
	}

	if commonResp.Code != 200 {
		return fmt.Errorf("request bigdata code is not 200, code: %d, message: %s", commonResp.Code, commonResp.Message)
	}

	return nil
}

type CommonResp struct {
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}
