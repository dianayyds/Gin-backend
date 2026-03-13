package fileprocess

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cihub/seelog"
)

var GBucket *oss.Bucket
var GBucket2 *oss.Bucket
var GBucketFra *oss.Bucket

var CFG_ENDPINT string
var CFG_ENDPINT2 string
var CFG_ENDPINT_FRA string
var CFG_ACCESSKEYID string
var CFG_ACCESSKEYID2 string
var CFG_ACCESSKEYID_FRA string
var CFG_ACCESSKEYSECRET string
var CFG_ACCESSKEYSECRET2 string
var CFG_ACCESSKEYSECRET_FRA string
var CFG_BUCKETNAME string
var CFG_BUCKETNAME2 string
var CFG_BUCKETNAME_FRA string
var CFG_OSSFILEPATH string
var CFG_LOCALRECORDINGFILEPATH string
var CFG_LOCALREPORTFILEPATH string
var CFG_RECORDINGURL string
var CFG_DOWNLOADURL string

type OSSCfg struct {
	Endpoint               string `json:"endpoint"`
	Endpoint2              string `json:"endpoint2"`
	EndpointFra            string `json:"endpointfra"`
	AccessKeyId            string `json:"accesskeyid"`
	AccessKeyId2           string `json:"accesskeyid2"`
	AccessKeyIdFra         string `json:"accesskeyidfra"`
	AccessKeySecret        string `json:"accesskeysecret"`
	AccessKeySecret2       string `json:"accesskeysecret2"`
	AccessKeySecretFra     string `json:"accesskeysecretfra"`
	BucketName             string `json:"bucketname"`
	BucketName2            string `json:"bucketname2"`
	BucketNameFra          string `json:"bucketnamefra"`
	OssFilePath            string `json:"ossfilepath"`
	LocalRecordingFilePath string `json:"localrecordingfilepath"`
	LocalReportFilePath    string `json:"localreportfilepath"`
	RecordingURL           string `json:"recordingUrl"`
	DownloadURL            string `json:"downloadUrl"`
}

func LoadConfigData(data string) int {
	configData := OSSCfg{}
	if err := json.Unmarshal([]byte(data), &configData); err != nil {
		seelog.Error("parse db cfg error ", err)
		return -1
	}

	CFG_ENDPINT = configData.Endpoint
	if CFG_ENDPINT == "" {
		seelog.Error("Get OSS endpoint error")
		return -1
	}

	CFG_ENDPINT2 = configData.Endpoint2
	if CFG_ENDPINT2 == "" {
		seelog.Error("Get OSS endpoint2 error")
		return -1
	}

	CFG_ENDPINT_FRA = configData.EndpointFra
	if CFG_ENDPINT_FRA == "" {
		seelog.Error("Get OSS CFG_ENDPINT_FRA error")
		return -1
	}

	CFG_ACCESSKEYID = configData.AccessKeyId
	if CFG_ACCESSKEYID == "" {
		seelog.Error("Get access key id error")
		return -1
	}

	CFG_ACCESSKEYID2 = configData.AccessKeyId2
	if CFG_ACCESSKEYID2 == "" {
		seelog.Error("Get access key2 id error")
		return -1
	}

	CFG_ACCESSKEYID_FRA = configData.AccessKeyIdFra
	if CFG_ACCESSKEYID_FRA == "" {
		seelog.Error("Get access CFG_ACCESSKEYID_FRA id error")
		return -1
	}

	CFG_ACCESSKEYSECRET = configData.AccessKeySecret
	if CFG_ACCESSKEYSECRET == "" {
		seelog.Error("Get access key secret error")
		return -1
	}

	CFG_ACCESSKEYSECRET2 = configData.AccessKeySecret2
	if CFG_ACCESSKEYSECRET2 == "" {
		seelog.Error("Get access key secret2 error")
		return -1
	}

	CFG_ACCESSKEYSECRET_FRA = configData.AccessKeySecretFra
	if CFG_ACCESSKEYSECRET_FRA == "" {
		seelog.Error("Get access CFG_ACCESSKEYSECRET_FRA secret error")
		return -1
	}

	CFG_BUCKETNAME = configData.BucketName
	if CFG_BUCKETNAME == "" {
		seelog.Error("Get bucket name error")
		return -1
	}

	CFG_BUCKETNAME2 = configData.BucketName2
	if CFG_BUCKETNAME2 == "" {
		seelog.Error("Get bucket2 name error")
		return -1
	}

	CFG_BUCKETNAME_FRA = configData.BucketNameFra
	if CFG_BUCKETNAME_FRA == "" {
		seelog.Error("Get bucket CFG_BUCKETNAME_FRA error")
		return -1
	}

	CFG_OSSFILEPATH = configData.OssFilePath
	if CFG_OSSFILEPATH == "" {
		seelog.Error("Get oss file path error")
		return -1
	}

	CFG_LOCALRECORDINGFILEPATH = configData.LocalRecordingFilePath
	if CFG_LOCALRECORDINGFILEPATH == "" {
		seelog.Error("Get local recording file path error")
		return -1
	}

	CFG_LOCALREPORTFILEPATH = configData.LocalReportFilePath
	if CFG_LOCALREPORTFILEPATH == "" {
		seelog.Error("Get local report file path error")
		return -1
	}

	CFG_RECORDINGURL = configData.RecordingURL
	if CFG_RECORDINGURL == "" {
		seelog.Error("Get local recordingurl file path error")
		return -1
	}

	CFG_DOWNLOADURL = configData.DownloadURL
	if CFG_DOWNLOADURL == "" {
		seelog.Error("Get local downloadurl file path error")
		return -1
	}
	return 0
}

func LoadConfig(cfg string) int {
	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		seelog.Error(err)
		return -1
	}

	return LoadConfigData(string(file))
}

func InitOSS() error {
	client, err := oss.New(CFG_ENDPINT, CFG_ACCESSKEYID, CFG_ACCESSKEYSECRET)
	if err != nil {
		seelog.Errorf("new oss Error:", err)
		return err
	}

	var err2, err3, err4, err5 error
	GBucket, err2 = client.Bucket(CFG_BUCKETNAME)
	if err2 != nil {
		seelog.Errorf("create bucket Error:", err2)
		return err2
	}

	client2, err3 := oss.New(CFG_ENDPINT2, CFG_ACCESSKEYID2, CFG_ACCESSKEYSECRET2)
	if err3 != nil {
		seelog.Errorf("new oss2 Error:", err3)
		return err3
	}
	GBucket2, err4 = client2.Bucket(CFG_BUCKETNAME2)
	if err4 != nil {
		seelog.Errorf("create bucket2 Error:", err4)
		return err4
	}

	clientfra, err3 := oss.New(CFG_ENDPINT_FRA, CFG_ACCESSKEYID_FRA, CFG_ACCESSKEYSECRET_FRA)
	if err3 != nil {
		seelog.Errorf("new ossfra Error:", err3)
		return err3
	}
	GBucketFra, err5 = clientfra.Bucket(CFG_BUCKETNAME_FRA)
	if err5 != nil {
		seelog.Errorf("create bucketfra Error:", err5)
		return err5
	}
	return nil
}

func UploadLocalFile(localFilePath, ossFilePath, fileName, fileBackup string) error {
	ossFile := ossFilePath + "/" + fileName
	localFile := localFilePath + "/" + fileName
	if fileBackup == "jkt" {
		err := GBucket2.PutObjectFromFile(ossFile, localFile)
		if err != nil {
			seelog.Errorf("upload local file to %s oss failed:%s, %s; %s", fileBackup, err.Error(), ossFile, localFile)
			return err
		}
		return nil
	}
	if fileBackup == "fra" {
		err := GBucketFra.PutObjectFromFile(ossFile, localFile)
		if err != nil {
			seelog.Errorf("upload local file to %s oss failed:%s, %s; %s", fileBackup, err.Error(), ossFile, localFile)
			return err
		}
		return nil
	}
	err := GBucket.PutObjectFromFile(ossFile, localFile)
	if err != nil {
		seelog.Errorf("upload local file to %s oss failed:%s, %s; %s", fileBackup, err.Error(), ossFile, localFile)
		return err
	}
	return nil
}

func DownloadOssFile(ossFilePath string, localFilePath string, fileName string) (string, error) {
	ff := ossFilePath + "/" + fileName
	lf := localFilePath + "/" + fileName
	seelog.Info("oss file path:", ff)
	err := GBucket.GetObjectToFile(ff, lf)
	fileBackup := "hk"
	if err != nil {
		err2 := GBucket2.GetObjectToFile(ff, lf)
		fileBackup = "jkt"
		if err2 != nil {
			err3 := GBucketFra.GetObjectToFile(ff, lf)
			fileBackup = "fra"
			if err3 != nil {
				// seelog.Errorf("download oss file to local failed 1:%s, file:%s", err.Error(), ff)
				// seelog.Errorf("download oss file to local failed 2:%s, file:%s", err2.Error(), ff)
				seelog.Errorf("download oss file to local failed: filename:%s, err1:%s, err2:%s, err3:%s", ff, err.Error(), err2.Error(), err3.Error())
				return "", err3
			}
		}

	}
	return fileBackup, nil
}

func DownloadExistOssFile(ossFilePath string, localFilePath string, fileName string) (string, error) {
	ff := ossFilePath
	lf := localFilePath + "/" + fileName
	seelog.Info("oss exist file path:", ff)
	if strings.HasPrefix(ossFilePath, "hk") {
		ff = strings.Replace(ff, "hk/", "", -1)
		err := GBucket.GetObjectToFile(ff, lf)
		if err != nil {
			seelog.Errorf("download exist oss file to local failed hk:%s, file:%s", err.Error(), ff)
			return "", err
		}
		return "", nil
	}
	if strings.HasPrefix(ossFilePath, "jkt") {
		ff = strings.Replace(ff, "jkt/", "", -1)
		err := GBucket2.GetObjectToFile(ff, lf)
		if err != nil {
			seelog.Errorf("download exist oss file to local failed jkt:%s, file:%s", err.Error(), ff)
			return "", err
		}
		return "", nil
	}
	if strings.HasPrefix(ossFilePath, "fra") {
		ff = strings.Replace(ff, "fra/", "", -1)
		err := GBucketFra.GetObjectToFile(ff, lf)
		if err != nil {
			seelog.Errorf("download exist oss file to local failed fra:%s, file:%s", err.Error(), ff)
			return "", err
		}
		return "", nil
	}

	return "", nil
}
