package fileprocess

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"rap_backend/dao"
	"time"

	"github.com/cihub/seelog"
)

func PrepareTaskRecordingFiles(taskId string, callidList []string) {
	subPath := path.Join(CFG_LOCALRECORDINGFILEPATH, taskId)
	_, checkErr := os.Stat(subPath)
	if checkErr != nil {
		res := os.Mkdir(subPath, os.ModePerm)
		if res != nil {
			seelog.Errorf("Create recording sub path failed, %s", subPath)
		}
	}
	//查询是否已存在
	recExist := make(map[string]string, 0)
	recs, err := dao.GetCallRecordingInfoByCallID(callidList)
	if err == nil && recs != nil {
		for _, rec := range *recs {
			recExist[rec.CallId] = rec.Address
		}
	}

	for _, callid := range callidList {
		if _, ok := recExist[callid]; ok {
			//recording表已存在，直接下载
			DownloadExistRecordingFile(subPath, callid)
			continue
		}
		fileBackup, err := DownloadRecordingFile(subPath, callid)
		if err != nil {
			seelog.Errorf("download recroding files failed:%s", err.Error())
			continue
		}

		time.Sleep(500 * time.Millisecond)

		err2 := CreateSubChannelRecordingFile(subPath, callid)
		if err2 != nil {
			seelog.Errorf("Create Sub Channel Recording File failed:%s", err2.Error())
			continue
		}
		//上传到oss
		ossRappath := "rap/" + time.Now().Format("200601")
		ossOri := fileBackup + "/" + CFG_OSSFILEPATH + "/" + callid + ".wav"
		ossRap_a := fileBackup + "/" + ossRappath + "/" + callid + "_a.wav"
		ossRap_b := fileBackup + "/" + ossRappath + "/" + callid + "_b.wav"
		err3 := UploadSubChannelRecordingFile(subPath, ossRappath, callid, fileBackup)
		if err3 != nil {
			seelog.Errorf("UploadSubChannelRecordingFile Recording File failed:%s", err3.Error())
			continue
		}
		//入数据库
		tcs := make([]dao.CallidRecording, 0)
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid,
			Address: ossOri,
		})
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid + "_a",
			Address: ossRap_a,
		})
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid + "_b",
			Address: ossRap_b,
		})
		err = dao.CreateCallRecordingInBatches(tcs)
		if err != nil {
			seelog.Errorf("CreateCallRecordingInBatches failed:%s", err.Error())
			return
		}
		//清除本地文件
		// ori_localfilepath := subPath + "/" + callid + ".wav"
		// ch1_localfilepath := subPath + "/" + callid + "_a.wav"
		// ch2_localfilepath := subPath + "/" + callid + "_b.wav"
		// os.Remove(ori_localfilepath)
		// os.Remove(ch1_localfilepath)
		// os.Remove(ch2_localfilepath)
	}
	seelog.Infof("prepare task recording files finished.")
	return
}

func DownloadRecordingFile(localpath string, callid string) (string, error) {
	fileBackup, err := DownloadOssFile(CFG_OSSFILEPATH, localpath, callid+".wav")
	if err != nil {
		seelog.Errorf("download file failed:%s", err.Error())
		return "", err
	}
	return fileBackup, nil
}

func CreateSubChannelRecordingFile(localpath, callid string) error {
	ori_localfilepath := localpath + "/" + callid + ".wav"
	ch1_localfilepath := localpath + "/" + callid + "_a.wav"
	ch2_localfilepath := localpath + "/" + callid + "_b.wav"
	strCmd := "ffmpeg -i " + ori_localfilepath + " -ar 8000 -acodec pcm_s16le -map_channel 0.0.0 " + ch1_localfilepath + " -map_channel 0.0.1 " + ch2_localfilepath
	//strCmd := "ls"
	cmd := exec.Command("/bin/bash", "-c", strCmd)

	stdout, err0 := cmd.StdoutPipe()
	if err0 != nil {
		seelog.Errorf("Error:can not obtain stdout pipe for command:%s\n", err0.Error())
		return err0
	}

	if err := cmd.Start(); err != nil {
		seelog.Errorf("separate channel failed %s, %s,", callid, err)
		return err
	}

	_, err2 := ioutil.ReadAll(stdout)
	if err2 != nil {
		seelog.Errorf("ReadAll Stdout:", err2.Error())
		return err2
	}
	if err3 := cmd.Wait(); err3 != nil {
		seelog.Errorf("wait:", err3.Error())
		return err3
	}

	// seelog.Infof("stdout:\n %s", bytes)
	return nil
}

func UploadSubChannelRecordingFile(localpath, osspath, callid, fileBackup string) error {
	err := UploadLocalFile(localpath, osspath, callid+"_a.wav", fileBackup)
	if err != nil {
		seelog.Errorf("upload file %s failed:%s", callid, err.Error())
		return err
	}
	err2 := UploadLocalFile(localpath, osspath, callid+"_b.wav", fileBackup)
	if err2 != nil {
		seelog.Errorf("upload file %s failed:%s", callid, err2.Error())
		return err2
	}
	return nil
}

func DownloadExistRecordingFile(localpath string, callid string) (string, error) {
	rcid := []string{callid, callid + "_a", callid + "_b"}
	recs, err := dao.GetCallRecordingInfoByCallID(rcid)
	if err != nil {
		return "", err
	}

	for _, rec := range *recs {
		DownloadExistOssFile(rec.Address, localpath, rec.CallId+".wav")
	}
	return "", nil
}

func PrepareTaskRecordingFilesLost(taskId string) {
	var callidList = []string{}
	subPath := path.Join(CFG_LOCALRECORDINGFILEPATH, taskId)
	_, checkErr := os.Stat(subPath)
	if checkErr != nil {
		res := os.Mkdir(subPath, os.ModePerm)
		if res != nil {
			seelog.Errorf("Create recording sub path failed, %s", subPath)
		}
	}
	taskcalls, _, err := dao.GetTaskCallsByTaskId(taskId, "", 0, nil, 0, 5000)
	if err != nil {
		seelog.Errorf("GetTaskCallsByTaskId failed, %s", err.Error())
		return
	}
	for _, tc := range *taskcalls {
		callidList = append(callidList, tc.CallId)
	}
	//查询是否已存在
	recExist := make(map[string]string, 0)
	recs, err := dao.GetCallRecordingInfoByCallID(callidList)
	if err == nil && recs != nil {
		for _, rec := range *recs {
			recExist[rec.CallId] = rec.Address
		}
	}

	for _, callid := range callidList {
		if _, ok := recExist[callid]; ok {
			//recording表已存在，直接下载
			DownloadExistRecordingFile(subPath, callid)
			continue
		}
		fileBackup, err := DownloadRecordingFile(subPath, callid)
		if err != nil {
			seelog.Errorf("download recroding files failed:%s", err.Error())
			continue
		}

		time.Sleep(500 * time.Millisecond)

		err2 := CreateSubChannelRecordingFile(subPath, callid)
		if err2 != nil {
			seelog.Errorf("Create Sub Channel Recording File failed:%s", err2.Error())
			continue
		}
		//上传到oss
		ossRappath := "rap/" + time.Now().Format("200601")
		ossOri := fileBackup + "/" + CFG_OSSFILEPATH + "/" + callid + ".wav"
		ossRap_a := fileBackup + "/" + ossRappath + "/" + callid + "_a.wav"
		ossRap_b := fileBackup + "/" + ossRappath + "/" + callid + "_b.wav"
		err3 := UploadSubChannelRecordingFile(subPath, ossRappath, callid, fileBackup)
		if err3 != nil {
			seelog.Errorf("UploadSubChannelRecordingFile Recording File failed:%s", err3.Error())
			continue
		}
		//入数据库
		tcs := make([]dao.CallidRecording, 0)
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid,
			Address: ossOri,
		})
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid + "_a",
			Address: ossRap_a,
		})
		tcs = append(tcs, dao.CallidRecording{
			CallId:  callid + "_b",
			Address: ossRap_b,
		})
		err = dao.CreateCallRecordingInBatches(tcs)
		if err != nil {
			seelog.Errorf("CreateCallRecordingInBatches failed:%s", err.Error())
			return
		}
	}
	return
}
