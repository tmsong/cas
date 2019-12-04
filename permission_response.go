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
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	EmployeeId        string `json:"employeeId"`
	Type              int    `json:"type"`
	Status            int    `json:"status"`
	ManagerId         int64  `json:"managerId"`
	ManagerName       string `json:"managerName"`
	ManagerAccount    string `json:"managerAccount"`
	ManagerEmail      string `json:"managerEmail"`
	ManagerEmployeeId string `json:"managerEmployeeId"`
	DeptName          string `json:"deptName"`
	DeptCode          string `json:"deptCode"`
	DeptId            string `json:"deptId"`
}

type DepartmentInfoRespose struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	SuperId             int    `json:"superId"`
	SuperName           string `json:"superName"`
	Code                string `json:"code"`
	NcDeptId            string `json:"ncDeptId"`
	NcDeptSuperId       string `json:"ncDeptSuperId"`
	NcDeptSuperName     string `json:"ncDeptSuperName"`
	ManagerAccount      string `json:"managerAccount"`
	ManagerEmail        string `json:"managerEmail"`
	ManagerName         string `json:"managerName"`
	DepartCategoryValue string `json:"departCategoryValue"`
	Status              string `json:"status"`
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
