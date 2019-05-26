package scim2

type EnterpriseUser struct {
	EmployeeNumber string
	CostCenter     string
	Organization   string
	Division       string
	Department     string
	Manager        Manager
}

type Manager struct {
	DisplayName string
	Ref         string
	Value       string
}
