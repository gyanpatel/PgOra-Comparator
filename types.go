package main

//errorMessages to for invalid login
type ErrorMessages struct {
	LoginError         string
	SessionTimeOutMins int
}

//errorPage represents shows an error message
type ErrorPage struct {
	ErrorMsg string
}

// SecretDetails to hold secret deatils required to run the web app
type SecretDetails struct {
	CognitoUserPoolClientID     string `json:"cognitoUserPoolClientID"`
	CognitoUserPoolClientSecret string `json:"cognitoUserPoolClientSecret"`
	Dbname                      string `json:"dbname"`
	Port                        string `json:"port"`
	User                        string `json:"username"`
	Password                    string `json:"password"`
	Host                        string `json:"host"`
	UserPassResetURL            string `json:"userPassResetURL"`
	SessionTimeOutMins          int    `json:"sessionTimeOutMins"`
	WebAppPort                  string `json:"webAppPort"`
	TokenURL                    string `json:"tokenURL"`
	RedirectURI                 string `json:"redirectURI"`
	AuthURL                     string `json:"authURL"`
	UserIfoURL                  string `json:"userIfoURL"`
	DbSslMode                   string `json:"dbsslmode"`
}

type CommonPageItems struct {
	UserName         string
	PageDesc         string
	TableQueryList   []TableQueryList
	ComparisonResult []ComparisonResult
}

type TableList struct {
	TableList string `json:"tablename"`
}

type TableQueryList struct {
	TableName string `json:"tablename"`
	PgQuery   string `json:"pgquery"`
	OraQuery  string `json:"oraquery"`
}

type ComparisonResult struct {
	TableName string
	Date      string
	Result    string
	PGCount   int
	OraCount  int
	DataDiff  []DataDiff
}
type DataDiff struct {
	DbName string
	Data   string
}
