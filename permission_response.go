package cas

// PermissionResponse captures authenticated user information
type PermissionResponse struct {
	Success bool        //请求是否成功
	Code    int64       //响应状态码
	Message string      //响应信息
	Data    interface{} // 响应对象
}

type HasPermissionResponse struct {
}

type UserInfoResponse struct {
	UserId      int64  `json:"userId"`
	Account     string `json:"account"`
	NameZh      string `json:"nameZh"`
	NameEn      string `json:"nameEn"`
	NameDisplay string `json:"nameDisplay"`
	Email       string `json:"email"`
	Type        int    `json:"type"`
	Phone       string `json:"phone"`
	ManagerId   int64  `json:"managerId"`
	ManagerName string `json:"managerName"`
	EmployeeId  string `json:"employeeId"`
	HbUserId    int64  `json:"hbUserId"`
	Department  string `json:"department"`
	Job         string `json:"job"`
	JobDesc     string `json:"jobDesc"`
}

type UserInfoDetailResponse struct {
	UserId            int64  `json:"userId"`
	Account           string `json:"account"`
	NameZh            string `json:"nameZh"`
	NameEn            string `json:"nameEn"`
	NameDisplay       string `json:"displayName"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	EmployeeId        string `json:"employeeId"`
	Type              int    `json:"type"`
	Status            int    `json:"status"`
	ManagerId         int64  `json:"managerId"`
	ManagerName       string `json:"managerName"`
	ManagerNameEn     string `json:"managerNameEn"`
	ManagerNameZh     string `json:"managerNameZh"`
	ManagerAccount    string `json:"managerAccount"`
	ManagerEmail      string `json:"managerEmail"`
	ManagerEmployeeId string `json:"managerEmployeeId"`
	DeptName          string `json:"deptName"`
	DeptCode          string `json:"deptCode"`
	DeptId            int64  `json:"deptId"`
	JoinDate          string `json:"joinDate"`
	BirthDay          string `json:"birthDay"`
	OfficeAddress     string `json:"officeAddress"`
}

type UserInfoVagueResponse struct {
	UserId       int64  `json:"userId"`
	Account      string `json:"account"`
	NameZh       string `json:"nameZh"`
	NameEn       string `json:"nameEn"`
	NameDisplay  string `json:"nameDisplay"`
	Email        string `json:"email"`
	EmployeeId   string `json:"employeeId"`
	DepartmentId int64  `json:"departmentId"`
}

type DepartmentInfoResponse struct {
	Id                  int64  `json:"id"`
	Name                string `json:"name"`
	SuperId             int64  `json:"superId"`
	SuperName           string `json:"superName"`
	Code                string `json:"code"`
	NcDeptId            string `json:"ncDeptId"`
	NcDeptSuperId       string `json:"ncDeptSuperId"`
	NcDeptSuperName     string `json:"ncDeptSuperName"`
	ManagerId           int64  `json:"managerId"`
	ManagerEmployeeId   string `json:"managerEmployeeId"`
	ManagerAccount      string `json:"managerAccount"`
	ManagerEmail        string `json:"managerEmail"`
	ManagerName         string `json:"managerName"`
	ManagerNameZh       string `json:"managerNameZh"`
	ManagerNameEn       string `json:"managerNameEn"`
	DepartCategoryValue string `json:"departCategoryValue"`
	Status              int    `json:"status"`
}

type PermissionListResponse struct {
	PermissionId   int64  `json:"permissionId"`
	PermissionKey  string `json:"permissionKey"`
	Url            string `json:"url"`
	PermissionName string `json:"permissionName"`
}

type RoleListResponse struct {
	RoleId    int64  `json:"roleId"`
	RoleKey   string `json:"roleKey"`
	RoleName  string `json:"roleName"`
	ValidTime string `json:"validTime"`
}

type UserPermissionListResponse struct {
	AppId          int64  `json:"appId"`
	PermissionKey  string `json:"permissionKey"`
	PermissionName string `json:"permissionName"`
	Url            string `json:"url"`
	Method         int64  `json:"method"`
	ParentId       int64  `json:"parentId"`
	IsMenu         int    `json:"isMenu"`
	OrderBy        int    `json:"orderBy"`
	Remark         string `json:"remark"`
	Status         int    `json:"status"`
	CreateTime     int64  `json:"createTime"`
	UpdateTime     int64  `json:"updateTime"`
}

type GetSsoUserByDDInfoResponse struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	NameZh      string `json:"nameZh"`
	NameEn      string `json:"nameEn"`
	NameDisplay string `json:"nameDisplay"`
	Department  string `json:"department"`
	Email       string `json:"email"`
	EmployeeId  string `json:"employeeId"`
	DingdingUid string `json:"dingdingUid"`
	DeviceId    string `json:"deviceId"`
}

type OfficeSiteDetailResponse struct {
	OfficeAddressId string `json:"officeAddressId"`
	OfficeAddress   string `json:"officeAddress"`
}
