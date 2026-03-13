package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cihub/seelog"
)

const (
	TASK_STATUS_NOTDISTRIBUTE = "not distribute"
	TASK_STATUS_NOTREADY      = "not ready"
	TASK_STATUS_NOTSTARTED    = "not started"
	TASK_STATUS_LABELING      = "labeling"
	TASK_STATUS_LABELFINISHED = "label finished"
	TASK_STATUS_APPROVING     = "approving"
	TASK_STATUS_APPROVED      = "approved"

	//TASK_STATUS_INIT       = "initial"    //初始化

	//task status
	TASK_STATUS_CREATED   = "created"   //已创建、待分配
	TASK_STATUS_ALLOCATED = "allocated" //已分配、待标注							callId status - 1

	TASK_STATUS_ANALYSTING = "analysting" //结果分析中					task status
	TASK_STATUS_AUDITING   = "auditing"   // 审核中
	//callid status

	TASK_STATUS_ANNOTATING = "annotating" //标注中
	TASK_STATUS_ANNOTATED  = "annotated"  //已标注						callId status - 2       realtion status - 标注完成

	CALLID_STATUS_PRE_AUDIT = "pre_audit" //待审核						callId status - 3
	TASK_STATUS_AUDITED     = "audited"   //已审核							callId status - 5

	TASK_STATUS_PRE_ANALYST = "pre_analyst" //待结果分析					callId status - 6

	TASK_STATUS_ANALYSTED = "analysted" //结果分析完成					callId status - 8
	TASK_STATUS_COMPLETED = "completed" //已完成						callId status - 9
	TASK_STATUS_DELETED   = "deleted"   //已删除

	TASK_PAGE_TAB_MANAGE  = "taskmanage" //任务管理
	TASK_PAGE_TAB_ALLOCAT = "allocator"  //任务分配
	TASK_PAGE_TAB_ANNOTAT = "annotator"  //通话标注
	TASK_PAGE_TAB_AUDIT   = "auditor"    //通话审核
	TASK_PAGE_TAB_ANALYS  = "analysts"   //结果分析

	USER_TYPE_ALLOCATOR = "user_allocator"
	USER_TYPE_ANNOTATOR = "user_annotator"
	USER_TYPE_AUDITOR   = "user_auditor"
	USER_TYPE_ANALYSTS  = "user_analysts"

	EXCEL_DOWNLOAD_PATH = "https://rap.airudder.com/report/"

	LABELWORK_STATUS_NOTSTART = "not labeled"
	LABELWORK_STATUS_FINISHED = "finished"
	LABELWORK_STATUS_APPROVED = "approved"

	RECORDING_PATH = "https://rap.airudder.com/recording/"

	TASK_SEARCH_TYPE_ID       = "id"
	TASK_SEARCH_TYPE_NAME     = "name"
	TASK_SEARCH_TYPE_CREATOR  = "creator"
	TASK_SEARCH_TYPE_STATUS   = "status"
	TASK_SEARCH_TYPE_OPERATOR = "operator"
)

var (

	//task 状态流程顺序
	TaskStatusFlow = []string{
		TASK_STATUS_CREATED,    //创建
		TASK_STATUS_ALLOCATED,  //已分配
		TASK_STATUS_ANNOTATING, //标注中
		TASK_STATUS_ANNOTATED,  //标注完待审核
		TASK_STATUS_AUDITING,   //审核中
		TASK_STATUS_AUDITED,    //已审核待结果分析
		TASK_STATUS_ANALYSTING, //结果分析中
		TASK_STATUS_COMPLETED,  //已完成
	}

	//callid 状态流转顺序
	CallIdStatusFlow = []string{
		TASK_STATUS_CREATED,     //创建, 待分配
		TASK_STATUS_ALLOCATED,   //已分配，待标注
		TASK_STATUS_ANNOTATED,   //已标注
		CALLID_STATUS_PRE_AUDIT, //待审核
		TASK_STATUS_AUDITED,     //已审核
		TASK_STATUS_PRE_ANALYST, //待结果分析
		TASK_STATUS_ANALYSTED,   //结果分析完成
		TASK_STATUS_COMPLETED,   //已完成
	}
)

func IsAfterCurrentTaskStatus(currentStatus string, com string) bool {
	var currentStatusIndex, comIndex int
	for i, status := range TaskStatusFlow {
		if status == currentStatus {
			currentStatusIndex = i
		}
		if status == com {
			comIndex = i
		}
	}
	return currentStatusIndex <= comIndex
}

// GetCallIdRemainStatus 获取status以及之后的状态
func GetCallIdRemainStatus(status string) []string {
	for i, s := range CallIdStatusFlow {
		if s == status {
			return CallIdStatusFlow[i:]
		}
	}
	return []string{}
}

type RingTypeTaskKey struct {
	Key       string
	Creator   string
	Allocator []string
}

var RingTypeTaskKeys = map[string]RingTypeTaskKey{
	"aba2cba8854447798727622ca16353cr": {
		Key:       "aba2cba8854447798727622ca16353cr",
		Creator:   "hua.liu@airudder.com", //login_name
		Allocator: []string{"hua.liu@airudder.com", "kyrie.ji@airudder.com", "junliang.chen@airudder.com"},
	},
	"885c6508ce504b2da7cf9a18a67bd5w4": {
		Key:       "885c6508ce504b2da7cf9a18a67bd5w4",
		Creator:   "rebecca.li@airudder.com",
		Allocator: []string{"PTE02"},
	},
}

func Post(url string, data interface{}, contentType string) string {

	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		seelog.Errorf("post %s failed %s", url, err.Error())
		return err.Error()
	}
	seelog.Infof("alrarm:%s", jsonStr)
	defer resp.Body.Close()

	result, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		seelog.Errorf("post %s read response failed %s", url, err2.Error())
		return err.Error()
	}

	return string(result)
}

func AsUInt8(bytearray []byte) []uint8 {
	var uintArray = make([]uint8, 0)
	for _, b := range bytearray {
		uintArray = append(uintArray, b)
	}
	return uintArray
}

//	func AsInt(int16array []int16) []int{
//		var byteArray = make([]byte, 0)
//		//bytesBuffer := bytes.NewBuffer(int16array)
//		for i:=range(int16array){
//			tmp := []byte{
//				byte(i),
//				byte(i >> 8),
//			}
//			byteArray = append(byteArray, tmp...)
//		}
//		bytebuff := bytes.NewBuffer(byteArray)
//		var data int64
//		binary.Read(bytebuff, binary.BigEndian, &data)
//
//		return int(data)
//	}
type Req struct {
	Name       string `json:"username"`         // 报警接收人 必填
	Phone      string `json:"phone"`            // 报警人电话 非必填
	Supervisor string `json:"supervisor"`       // 上级接收人
	SvPhone    string `json:"supervisor_phone"` // 上级接收人电话 非必填
	Product    string `json:"product"`          // 产品名 非必填
	Title      string `json:"title"`            // 报警标题 必填
	Body       string `json:"body"`             // 报警内容 必填
	Dept       string `json:"dept"`             // 部门 用来确定发到哪个飞书群里 待定 必填
	Level      int    `json:"level"`            // 报警等级 说明见上 必填
	Cc         string `json:"cc"`               // 同时通知的人 非必填 '|' 分割用户名
}

// Alarm ...
//func Alarm(title string, body string, level int) {
//	req := Req{
//		Name:      "king.qin",
//		Supervisor: "king.qin",
//		Title:      "[postsuperman]" + title,
//		Dept:       "VOIP",
//		Body:       body,
//		Level:      level,
//		Cc:         "xiao.li",
//	}
//	data, _ := json.Marshal(req)
//	//request, err := http.NewRequest("POST", config.AlarmUrl, bytes.NewBuffer(data))
//	request, err := http.NewRequest("POST", config.AlarmUrl, bytes.NewBuffer(data))
//	request.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	resp, err := client.Do(request)
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode == 200 {
//
//	} else {
//		// 自行处理
//	}
//}

func OutPutDataWithXLSX(list interface{}, headers []string, filepath string, filename string) (err error) {
	filename = filepath + "/" + filename
	seelog.Infof("report file path:%s", filename)
	_, fileerr := os.Lstat(filename)
	if fileerr == nil {
		fileerr2 := os.Remove(filename)
		if fileerr2 != nil {
			seelog.Errorf("remove file %s failed %s", filename, fileerr2.Error())
			return fileerr2
		}
	}

	sheet1 := "Sheet1"
	f := excelize.NewFile()

	character := string(65 + len(headers) - 1)

	//headers
	styleHeader, _ := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"},"font":{"bold":false,"italic":false,"family":"Calibri","size":10,"color":"#000000"}}`)
	f.SetCellStyle(sheet1, "A2", character+"2", styleHeader)

	for k, v := range headers {
		f.SetCellValue(sheet1, string(65+k)+"2", v)
	}

	///* -------------------- list content -------------------- */

	getValue := reflect.ValueOf(list)
	if getValue.Kind() != reflect.Slice {
		return errors.New("list must be slice")
	}
	length := getValue.Len()

	if length > 0 {
		line := 3
		for i := 0; i < length; i++ {
			value := getValue.Index(i)
			typel := value.Type()
			if typel.Kind() != reflect.Struct {
				return errors.New("list must be slice of struct")
			}

			lineChr := strconv.Itoa(line)
			f.SetCellStyle(sheet1, "A"+lineChr, character+lineChr, styleHeader)

			n := value.NumField()
			for i := 0; i < n; i++ {
				val := value.Field(i)
				f.SetCellValue(sheet1, string(65+i)+lineChr, val.Interface())
			}
			line++
		}
	}

	err = f.SaveAs(filename)
	if err != nil {
		seelog.Errorf("save excel %s error:%s", filename, err.Error())
		return err
	}

	return nil
}

//for [][]string

func OutPutDataWithXLSX2(list interface{}, headers []string, filepath string, filename string) (err error) {
	filename = filepath + "/" + filename
	seelog.Infof("report file path:%s", filename)
	_, fileerr := os.Lstat(filename)
	if fileerr == nil {
		fileerr2 := os.Remove(filename)
		if fileerr2 != nil {
			seelog.Errorf("remove file %s failed %s", filename, fileerr2.Error())
			return fileerr2
		}
	}

	sheet1 := "Sheet1"
	f := excelize.NewFile()

	character := string(65 + len(headers) - 1)

	//headers
	styleHeader, _ := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"},"font":{"bold":false,"italic":false,"family":"Calibri","size":10,"color":"#000000"}}`)
	f.SetCellStyle(sheet1, "A2", character+"2", styleHeader)

	for k, v := range headers {
		f.SetCellValue(sheet1, string(65+k)+"2", v)
	}

	///* -------------------- list content -------------------- */

	getValue := reflect.ValueOf(list)
	//if getValue.Kind() != reflect.Slice {
	//	return errors.New("list must be slice")
	//}
	length := getValue.Len()

	if length > 0 {
		line := 3
		for i := 0; i < length; i++ {
			value := getValue.Index(i)
			//typel := value.Type()
			//if typel.Kind() != reflect.Struct {
			//	return errors.New("list must be slice of struct")
			//}

			lineChr := strconv.Itoa(line)
			f.SetCellStyle(sheet1, "A"+lineChr, character+lineChr, styleHeader)

			n := value.Len()
			for i := 0; i < n; i++ {
				val := value.Index(i)
				f.SetCellValue(sheet1, string(65+i)+lineChr, val.String())
			}
			line++
		}
	}

	err = f.SaveAs(filename)
	if err != nil {
		seelog.Errorf("save excel %s error:%s", filename, err.Error())
		return err
	}

	return nil
}
