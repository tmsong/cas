package cas

type PolicyResponse struct {
	Success bool
	Code    int64
	Message string
	Data    []PolicyData
}

type PolicyData struct {
	Id         int64
	AppId      int64
	PolicyKey  string
	PolicyName string
	Remark     string
	Policys    []Policy
	Status     int64
	CreateTime int64
	UpdateTime int64
}

type Policy struct {
	Effect        string
	Actions       []string
	Resources     []Resource
	IpCondition   []IpCondition
	TimeCondition []TimeCondition
}

type Resource struct {
	FlagKey       string
	FlagOptionKey string
}

type IpCondition struct {
	Equal     string
	NotBelong string
	Belong    string
}

type TimeCondition struct {
	Equal   int64
	Greater int64
	Less    int64
}

type PolicyByUrlResponse struct {
	Success bool
	Code    int64
	Message string
	Data    []Policy
}

type FlagListResponse struct {
	Success bool
	Code    int64
	Message string
	Data    []FlagDetail
}

type FlagDetail struct {
	Id             int64
	AppId          int64
	FlagId         int64
	FlagoptionKey  string
	FlagoptionName string
	CustDesc       string
	Remark         string
	CreateTime     int64
	UpdateTime     int64
}
