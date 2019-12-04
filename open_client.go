/**
 * @note
 * open_client.go
 *
 * @author	songtianming
 * @date 	2019-12-04
 */
package cas

import (
	"errors"
	"github.com/tmsong/hlog"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type OpenClient struct {
	appId   int64
	appKey  string
	client  *http.Client
	logger  *hlog.Logger
	openUrl *url.URL
}

func NewOpenClient(client *http.Client, openURL *url.URL, l *hlog.Logger) *OpenClient {
	return &OpenClient{
		client:  client,
		openUrl: openURL,
		logger:  l,
	}
}

func (c *OpenClient) UserInfoDetailUrl(userId int64, employeeId string) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/get_user_info_detail"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	if len(employeeId) == 0 {
		params["userId"] = userId
	} else {
		params["employeeId"] = employeeId
	}
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) DepartmentInfoUrl(departmentId int64) (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/dept/get_dept"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	params["id"] = strconv.FormatInt(departmentId, 10)
	return u.String(), JsonEncode(params), nil
}

func (c *OpenClient) AllDepartmentInfoUrl() (string, string, error) {
	u, err := c.openUrl.Parse(path.Join(c.openUrl.Path, "api/open/sso/dept/get_all_dept"))
	if err != nil {
		return "", "", err
	}
	params := CreateBaseParams(c.appId, c.appKey)
	return u.String(), JsonEncode(params), nil
}

///////////////////////////////////////////////////////

func (c *OpenClient) UserInfoDetail(userId int64, employeeId string) (*UserInfoDetailResponse, error) {
	u, body, err := c.UserInfoDetailUrl(userId, employeeId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := &UserInfoDetailResponse{}
	InterfaceToStruct(r.Data, re)
	return re, nil
}

func (c *OpenClient) DepartmentInfo(departmentId int64) (*DepartmentInfoRespose, error) {
	u, body, err := c.DepartmentInfoUrl(departmentId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := &DepartmentInfoRespose{}
	InterfaceToStruct(r.Data, re)
	return re, nil
}

func (c *OpenClient) AllDepartmentInfo() ([]*DepartmentInfoRespose, error) {
	u, body, err := c.AllDepartmentInfoUrl()
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body, c.logger)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := []*DepartmentInfoRespose{}
	InterfaceToStruct(r.Data, re)
	return re, nil
}
