package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
)

func DoHttpPutRequest(url string, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	if len(contentType) == 0 {
		r.Header["Content-Type"] = []string{"application/json"}
	} else {
		r.Header["Content-Type"] = []string{contentType}
	}
	rs, err1 := http.DefaultClient.Do(r)
	if err1 != nil {
		return nil, err1
	}
	return rs, nil
}

func DoHttpPutJson(url string, request interface{}, result interface{}) error {
	var err error

	var b []byte
	if b, err = json.Marshal(request); err != nil {
		return err
	}

	var r *http.Response
	if r, err = DoHttpPutRequest(url, "", bytes.NewBuffer(b)); err != nil {
		return err
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	return nil
}

func DoHttpPostRequest(url string, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// fix EOF error
	// it prevents the connection from being re-used
	// see https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi/23963271
	r.Close = true

	if len(contentType) == 0 {
		r.Header["Content-Type"] = []string{"application/json"}
	} else {
		r.Header["Content-Type"] = []string{contentType}
	}
	// 设置超时时间
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	r.WithContext(ctx)
	rs, err1 := http.DefaultClient.Do(r)
	if err1 != nil {
		return nil, err1
	}
	return rs, nil
}

func DoHttpPostJson(url string, request interface{}, result interface{}) error {
	var err error

	var b []byte
	if b, err = json.Marshal(request); err != nil {
		return err
	}

	var r *http.Response
	if r, err = DoHttpPostRequest(url, "", bytes.NewBuffer(b)); err != nil {
		seelog.Errorf("http request, url:%s err:%s request:%s \n", url, err.Error(), string(b))
		return err
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		seelog.Errorf("http request, url:%s err:%s request:%s response:%s \n", url, err.Error(), string(b), string(body))
		return err
	}

	if r.StatusCode != 200 {
		return errors.New("http状态码错误: " + r.Status)
	}

	if err = json.Unmarshal(body, &result); err != nil {
		seelog.Errorf("http request, url:%s err:%s request:%s response:%s \n", url, err.Error(), string(b), string(body))
		return err
	}

	return nil
}

func DoHttpDeleteRequest(url string) (*http.Response, error) {
	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	rs, err1 := http.DefaultClient.Do(r)
	if err1 != nil {
		return nil, err1
	}
	return rs, nil
}

func DoHttpDeleteJson(url string, data interface{}) error {
	var err error
	var r *http.Response
	if r, err = DoHttpDeleteRequest(url); err != nil {
		return err
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}

	return nil
}

func DoHttpGetRequest(url string, head http.Header) (*http.Response, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if head != nil && len(head) > 0 {
		for k, val := range head {
			for _, v := range val {
				r.Header.Add(k, v)
			}
		}
	}
	// 设置超时时间
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	r.WithContext(ctx)
	rs, err1 := http.DefaultClient.Do(r)
	if err1 != nil {
		return nil, err1
	}
	return rs, nil
}

func DoHttpGetJson(url string, data interface{}, head http.Header) error {
	var err error
	var r *http.Response
	if r, err = DoHttpGetRequest(url, head); err != nil {
		return err
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(body, &data); err != nil {
		reason := fmt.Sprintf("url:%s, response:%s,reason:%s", url, string(body), err.Error())
		return errors.New(reason)
	}

	return nil
}

func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func DoHttpPostFormURLencode(endpoint string, data url.Values, result interface{}) error {
	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		seelog.Errorf("DoHttpPostFormURLencode url:%s, resp:%v, err:%s", endpoint, res, err.Error())
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}
	return nil
}
