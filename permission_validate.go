package cas

import (
	"errors"
	"net/http"
	"net/url"
	"path"
)

// NewServiceTicketValidator create a new *ServiceTicketValidator
func NewPermissionValidator(client *http.Client, permissionURL *url.URL, baseParams map[string]interface{}) *PermissionValidator {
	return &PermissionValidator{
		client:        client,
		permissionURL: permissionURL,
		baseParams:    baseParams,
	}
}

// ServiceTicketValidator is responsible for the validation of a service ticket
type PermissionValidator struct {
	client        *http.Client
	permissionURL *url.URL
	baseParams    map[string]interface{}
}

// ValidateTicket validates the service ticket for the given server. The method will try to use the service validate
// endpoint of the cas >= 2 protocol, if the service validate endpoint not available, the function will use the cas 1
// validate endpoint.
func (validator *PermissionValidator) HasPermission(userId int64, url string) error {
	u, body, err := validator.HasPermissionUrl(userId, url)
	if err != nil {
		return err
	}
	ret := PostByJson(u, body)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return errors.New("no permission")
	}
	return nil
}

func (validator *PermissionValidator) RoleList(userId int64) ([]RoleListResponse, error) {
	u, body, err := validator.RoleListUrl(userId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := []RoleListResponse{}
	InterfaceToStruct(r.Data, &re)
	return re, nil
}
func (validator *PermissionValidator) PermissionList(userId, roleId int64) ([]PermissionListResponse, error) {
	u, body, err := validator.PermissionListUrl(userId, roleId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := []PermissionListResponse{}
	InterfaceToStruct(r.Data, &re)
	return re, nil
}

func (validator *PermissionValidator) UserInfo(userId int64) (*UserInfoResponse, error) {
	u, body, err := validator.UserInfoUrl(userId)
	if err != nil {
		return nil, err
	}
	ret := PostByJson(u, body)
	r := PermissionResponse{}
	_ = JsonDecode(ret, &r)
	if r.Code != 200 {
		return nil, errors.New("error")
	}
	re := &UserInfoResponse{}
	InterfaceToStruct(r.Data, re)
	return re, nil
}

func (validator *PermissionValidator) HasPermissionUrl(userId int64, url string) (string, string, error) {
	u, err := validator.permissionURL.Parse(path.Join(validator.permissionURL.Path, "api/open/sso/has_premission"))
	if err != nil {
		return "", "", err
	}
	params := validator.baseParams
	params["userId"] = userId
	params["url"] = url
	return u.String(), JsonEncode(params), nil
}

//func (validator *PermissionValidator) HasPermissionForManagerUrl(serviceURL *url.URL, ticket string) (string, error) {
//	u, err := validator.permissionURL.Parse(path.Join(validator.permissionURL.Path, "api/open/sso/has_premission_for_manager"))
//	if err != nil {
//		return "", err
//	}
//	q := u.Query()
//	q.Add("service", sanitisedURLString(serviceURL))
//	q.Add("ticket", ticket)
//	u.RawQuery = q.Encode()
//	return u.String(), nil
//}
func (validator *PermissionValidator) RoleListUrl(userId int64) (string, string, error) {
	u, err := validator.permissionURL.Parse(path.Join(validator.permissionURL.Path, "api/open/upm/user/role_list"))
	if err != nil {
		return "", "", err
	}
	params := validator.baseParams
	params["userId"] = userId
	return u.String(), JsonEncode(params), nil
}
func (validator *PermissionValidator) PermissionListUrl(userId, roleId int64) (string, string, error) {
	u, err := validator.permissionURL.Parse(path.Join(validator.permissionURL.Path, "api/open/upm/role/permission_list"))
	if err != nil {
		return "", "", err
	}
	params := validator.baseParams
	params["userId"] = userId
	params["roleId"] = roleId
	return u.String(), JsonEncode(params), nil
}
func (validator *PermissionValidator) UserInfoUrl(userId int64) (string, string, error) {
	u, err := validator.permissionURL.Parse(path.Join(validator.permissionURL.Path, "api/open/sso/user_info"))
	if err != nil {
		return "", "", err
	}
	params := validator.baseParams
	params["userId"] = userId
	return u.String(), JsonEncode(params), nil
}
