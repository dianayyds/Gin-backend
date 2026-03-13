package service

import (
	"errors"
	"net/http"
	"net/url"
	"rap_backend/config"
	"rap_backend/dao"
	"rap_backend/utils"

	"github.com/cihub/seelog"
)

var (
	// sso_domain    = "https://login.partner.microsoftonline.cn/b1989b29-569a-4345-858a-8779ccbc69f9/"
	sso_key       = "ac6bbeb2-0116-425e-b689-f364adc278b1"
	sso_secret    = "REPLACE_WITH_SSO_CLIENT_SECRET"
	sso_token_url = "https://login.partner.microsoftonline.cn/b1989b29-569a-4345-858a-8779ccbc69f9/oauth2/v2.0/token"
	sso_user_url  = "https://microsoftgraph.chinacloudapi.cn/v1.0/me"
)

type SSOLoginDTO struct {
	Code string `json:"code"`
}

type SSOLoginRetDTO struct {
	Token string `json:"token"`
}

type SSOTokenResp struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	Error        string `json:"error"`
	ErrorCodes   []int  `json:"error_codes"`
}

type SSOUserInfo struct {
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Mail              string `json:"mail"`
	SurName           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
}

func SSOUserLogin(redirect_uri, code string) (string, error) {
	accessToken, err := GetSSOAccessToken(redirect_uri, code)
	if err != nil {
		return accessToken, err
	}
	//获取microsoft用户
	ssoUser, err := GetSSOUserInfo(accessToken)
	if err != nil || ssoUser.Mail == "" {
		return "", err
	}
	//获得uid
	var user = dao.User{
		LoginName: ssoUser.Mail,
		Password:  UserPasswordInit(""),
		UserName:  ssoUser.DisplayName,
		Email:     ssoUser.Mail,
		Status:    config.USER_STATUS_NORMAL,
		Roles:     "",
	}
	userID, err := user.CreateOrSelectUserInfo()
	seelog.Infof("SSOUserLogin %d, %v, %d, %v", user.UserId, user, userID, err)
	if err != nil || userID == 0 {
		return "", err
	}
	if user.Status != config.USER_STATUS_NORMAL {
		return "", errors.New("user is forbid")
	}
	//生成jwt token
	token, err := SetUserToken(user)
	return token, err
}

func GetSSOAccessToken(redirect_uri, code string) (string, error) {
	req := url.Values{
		"scope":         []string{"openid"},
		"redirect_uri":  []string{redirect_uri},
		"grant_type":    []string{"authorization_code"},
		"client_secret": []string{sso_secret},
		"client_id":     []string{sso_key},
		"code":          []string{code},
	}
	resp := SSOTokenResp{}
	err := utils.DoHttpPostFormURLencode(sso_token_url, req, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	return resp.AccessToken, nil

}

func GetSSOUserInfo(accessToken string) (SSOUserInfo, error) {
	resp := SSOUserInfo{}
	err := utils.DoHttpGetJson(sso_user_url, &resp, http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	})
	return resp, err
}
