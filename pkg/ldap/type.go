package ldap

type UserInfo struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Title       string `json:"title"`
	Mail        string `json:"mail"`
	DisplayName string `json:"displayName"`
}
