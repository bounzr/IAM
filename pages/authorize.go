package pages

//AuthorizePage contains data for authorize.html
type AuthorizePage struct {
	ClientName, ClientID, ClientURI string
	ScopesList []string
}
