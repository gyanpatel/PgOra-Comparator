package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	store           = sessions.NewCookieStore([]byte("asdaskdhasdhgsajdgasdsadksakdhasidoajsdousahdopj"))
	authenticatedYN = "authenticated-orapgcomp"
	sessoinKeyID    = "user-authenticated-orapgcomp"
	sessionUser     = "username-orapgcomp"
	pageErr         = "Internal server error"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//setupLogFile()
	code := r.URL.Query().Get("code")
	if len(code) > 0 {
		ssoLoginVal(w, r, code)
	}
	e := ErrorMessages{LoginError: "", SessionTimeOutMins: secretDetails.SessionTimeOutMins}
	t := template.Must(template.ParseFS(templates, "templates/login.html"))
	err := t.Execute(w, e)
	if err != nil {
		log.Println("ERROR:handleLogin Error occured Parsing - templates/login ", err)
		renderErrorPage(w, fmt.Errorf(pageErr+" Error occured Parsing - templates/login "))
		return
	}
}
func handleLogout(w http.ResponseWriter, r *http.Request) {
	//setupLogFile()
	session, _ := store.Get(r, sessoinKeyID)
	log.Println("Info : handleLogout", session)
	session.Values[authenticatedYN] = false
	errs := session.Save(r, w)
	if errs != nil {
		log.Println("ERROR:handleLogout Error occured session.Save ", errs)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func handleLoginVal(w http.ResponseWriter, r *http.Request) {
	errp := r.ParseForm()
	if errp != nil {
		log.Println("ERROR:handleLoginVal Error occured r.ParseForm() ", errp)

	}
	params := r.PostForm
	userName := params.Get("username")
	passWord := params.Get("password")
	log.Println("Info : handleLoginVal", userName, "- login Attempt")
	auth, err := AuthenticateUser(userName, passWord)
	if err != nil || auth == "N" {
		log.Println("ERROR : handleLoginVal", userName, "- login attempt failed ")
		e := ErrorMessages{LoginError: "Invalid login details", SessionTimeOutMins: secretDetails.SessionTimeOutMins}
		t := template.Must(template.ParseFS(templates, "templates/login.html"))
		errt := t.Execute(w, e)
		if errt != nil {
			log.Println("ERROR:handleLoginVal Error occured Parsing - templates/login ", errt)
		}
		return
	} else if auth == "R" {
		log.Println("Info : handleLoginVal", userName, "- login attempt redirected to  change password ")
		//http.Redirect(w, r, resetPassURL, http.StatusSeeOther)
		http.Redirect(w, r, secretDetails.UserPassResetURL, http.StatusSeeOther)

	} else if auth == "Y" {
		log.Println("Info : handleLoginVal", userName, "- login attempt successful ")
		log.Println("Info : handleLoginVal", userName, "Session timeout mins ", secretDetails.SessionTimeOutMins)
		session, _ := store.Get(r, sessoinKeyID)
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   60 * secretDetails.SessionTimeOutMins * 12,
			HttpOnly: true,
		}
		session.Values[authenticatedYN] = true
		session.Values[sessionUser] = userName
		errs := session.Save(r, w)
		if errs != nil {
			log.Println("ERROR:handleLoginVal Error occured session.Save ", errs)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		log.Println("ERROR : handleLoginVal", err)
		renderErrorPage(w, err)
		return

	}

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	userName, err := userSessionValidation(w, r)
	log.Println("Info : handleHome - after userSessionValidation")
	if err != nil {
		log.Println("ERROR : handleHome", err)
		return
	}
	errp := r.ParseForm()
	if errp != nil {
		log.Println("ERROR:handleHome Error occured r.ParseForm() ", errp)
	}
	homePageItems := CommonPageItems{UserName: userName, TableQueryList: tableQueryList, PageDesc: "Table Selection"}
	t := template.Must(template.ParseFS(templates, "templates/home.html", "templates/_menu.html", "templates/_sidenav.html", "templates/_footer.html"))
	errt := t.Execute(w, homePageItems)
	if errt != nil {
		log.Println("ERROR:handleHome Error occured Parsing - templates/home.html ", errt)
	}
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	userName, err := userSessionValidation(w, r)
	log.Println("Info : handleHome - after userSessionValidation")
	if err != nil {
		log.Println("ERROR : handleHome", err)
		return
	}
	errp := r.ParseForm()
	if errp != nil {
		log.Println("ERROR:handleHome Error occured r.ParseForm() ", errp)
	}
	selTableList := r.PostForm["tablelist"]
	log.Println("selTableList", selTableList)
	var selTableQueryList []TableQueryList
	if (r.URL.Path[1:]) == "compall" {
		selTableQueryList = tableQueryList
	} else {
		for _, sel := range tableQueryList {
			if contains(selTableList, sel.TableName) {
				selTableQueryList = append(selTableQueryList, sel)
			}
		}
	}
	comparisonResult, err := compareOraPgTables(selTableQueryList)
	if err != nil {
		log.Println("ERROR: Error occured whilst comparison ", err)
		renderErrorPage(w, fmt.Errorf("ERROR: Error occured whilst comparison %v", err))
		return
	}

	homePageItems := CommonPageItems{ComparisonResult: *comparisonResult, UserName: userName, TableQueryList: tableQueryList, PageDesc: "Comparison result"}
	t := template.Must(template.ParseFS(templates, "templates/home.html", "templates/_menu.html", "templates/_sidenav.html", "templates/_footer.html"))
	errt := t.Execute(w, homePageItems)
	if errt != nil {
		log.Println("ERROR:handleHome Error occured Parsing - templates/home.html ", errt)
	}
}

//history to display  comparison history
func handleHist(w http.ResponseWriter, r *http.Request) {
	userName, err := userSessionValidation(w, r)
	log.Println("Info : handleHist - after userSessionValidation")
	if err != nil {
		log.Println("ERROR : handleHist", err)
		return
	}
	errp := r.ParseForm()
	if errp != nil {
		log.Println("ERROR:handleHist Error occured r.ParseForm() ", errp)
	}

	homePageItems := CommonPageItems{UserName: userName, TableQueryList: tableQueryList, PageDesc: "Comparison History"}
	t := template.Must(template.ParseFS(templates, "templates/history.html", "templates/_menu.html", "templates/_sidenav.html", "templates/_footer.html"))
	errt := t.Execute(w, homePageItems)
	if errt != nil {
		log.Println("ERROR:handleHist Error occured Parsing - templates/home.html ", errt)
	}
}

func renderErrorPage(w http.ResponseWriter, errorMsg error) {
	errPageItems := CommonPageItems{PageDesc: errorMsg.Error()}
	t := template.Must(template.ParseFS(templates, "templates/error.html", "templates/_menu.html", "templates/_sidenav.html", "templates/_footer.html"))
	errt := t.Execute(w, errPageItems)
	if errt != nil {
		log.Println("ERROR:renderErrorPage Error occured Parsing - templates/error.html ", errt)
	}
}

func userSessionValidation(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := store.Get(r, sessoinKeyID)
	userName, _ := session.Values[sessionUser].(string)
	log.Println("Info: userSessionValidation ", userName)

	if err != nil {
		log.Println("ERROR: userSessionValidation ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	if auth, ok := session.Values[authenticatedYN].(bool); !ok || !auth || len(userName) == 0 {
		e := ErrorMessages{LoginError: "Your sesssion has expired or you haven't logged in, please login .", SessionTimeOutMins: secretDetails.SessionTimeOutMins}
		t := template.Must(template.ParseFS(templates, "templates/login.html"))
		errt := t.Execute(w, e)
		if errt != nil {
			log.Println("ERROR:userSessionValidation Error occured Parsing - templates/login.html ", errt)
		}
		return "", fmt.Errorf(userName, " Your sesssion has expired, please login again")
	}
	return userName, nil
}

func handleAwsSsoLogin(w http.ResponseWriter, r *http.Request) {
	//setupLogFile()
	code := r.URL.Query().Get("code")
	if len(code) > 0 {
		ssoLoginVal(w, r, code)
		return
	}
	conf := &oauth2.Config{
		ClientID:     secretDetails.CognitoUserPoolClientID,
		ClientSecret: secretDetails.CognitoUserPoolClientSecret,
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			TokenURL: secretDetails.TokenURL,
			AuthURL:  secretDetails.AuthURL,
		},
	}
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline) + "&redirect_uri=" + secretDetails.RedirectURI

	http.Redirect(w, r, url, http.StatusSeeOther)
	e := ErrorMessages{LoginError: "", SessionTimeOutMins: secretDetails.SessionTimeOutMins}
	t := template.Must(template.ParseFS(templates, "templates/login.html"))
	errt := t.Execute(w, e)
	if errt != nil {
		log.Println("ERROR:handleAwsSsoLogin Error occured Parsing - templates/login.html ", errt)
	}
}
func ssoLoginVal(w http.ResponseWriter, r *http.Request, code string) {
	log.Println("Code ", code)

	if len(code) > 0 {
		secretHash := "Basic " + base64.StdEncoding.EncodeToString([]byte(secretDetails.CognitoUserPoolClientID+":"+secretDetails.CognitoUserPoolClientSecret))
		tokenURL := secretDetails.TokenURL
		authDetails := url.Values{}
		authDetails.Set("grant_type", "authorization_code")
		authDetails.Set("client_id", secretDetails.CognitoUserPoolClientID)
		authDetails.Set("redirect_uri", secretDetails.RedirectURI)
		authDetails.Set("code", code)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(authDetails.Encode())) // URL-encoded payload
		if err != nil {
			log.Println("ERROR: Error occured http.NewRequest", err)
			renderErrorPage(w, fmt.Errorf("ERROR: Error occured http.NewRequest, Unable to auhtenticate "))
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", secretHash)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("ERROR: Error occured client.Do", err, "req", req)
			renderErrorPage(w, fmt.Errorf("ERROR: Error client.Do, Unable to auhtenticate "))
			return
		}
		log.Println("resp.Status", resp.Status)
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ERROR: Read Token Response  ", err)
			renderErrorPage(w, fmt.Errorf("ERROR: Read Token Response, Unable to auhtenticate "))
			return

		}
		var accessTokenM map[string]interface{}
		errmap := json.Unmarshal([]byte(string(respBody)), &accessTokenM)
		if errmap != nil {
			log.Println("ERROR: Parsing Token Response  ", errmap)
			renderErrorPage(w, fmt.Errorf("ERROR: Parsing Token Response, Unable to auhtenticate "))
			return
		}
		accessToken := accessTokenM["access_token"].(string)
		// Use the custom HTTP client when requesting a token.
		userIfoURL := secretDetails.UserIfoURL

		requ, err := http.NewRequest(http.MethodGet, userIfoURL, nil)

		if err != nil {
			log.Println("ERROR: Error occured http.NewRequest userInfo", err)
			renderErrorPage(w, fmt.Errorf("ERROR: Unable to retrieve user details"))
			return
		}
		requ.Header.Add("Authorization", "Bearer "+accessToken)
		requ.Header.Add("Content-Type", "application/json")
		resu, err := client.Do(requ)
		if err != nil {
			log.Println("ERROR: client.Do(requ)  ", err)
			renderErrorPage(w, fmt.Errorf("ERROR: Unable to retrieve user details"))
			return
		}
		defer resu.Body.Close()

		respBodyU, err := ioutil.ReadAll(resu.Body)
		if err != nil {
			log.Println("ERROR: Parsing respBody UserName  ", errmap)
			renderErrorPage(w, fmt.Errorf("ERROR: Unable to retrieve user details"))
			return
		}
		var userNameM map[string]interface{}
		erru := json.Unmarshal([]byte(string(respBodyU)), &userNameM)
		if erru != nil {
			log.Println("ERROR: Parsing UserName  ", errmap)
			renderErrorPage(w, fmt.Errorf("ERROR: Unable to retrieve user details"))
			return
		}
		// Key usernname hardcoded as AWS documentaion uses key "username" for user details
		userName := userNameM["username"].(string)
		log.Println("Info : SSOLoginVal", userName, "- login attempt successful ")
		log.Println("Info : SSOLoginVal", userName, "Session timeout mins ", secretDetails.SessionTimeOutMins)
		session, _ := store.Get(r, sessoinKeyID)
		session.Values[authenticatedYN] = true
		session.Values[sessionUser] = userName
		errs := session.Save(r, w)
		if errs != nil {
			log.Println("ERROR:ssoLoginVal Error occured session.Save ", errs)
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// Below custom functions check
//if slice is a subset of another slice
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
func subslice(s1 []string, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !contains(s2, e) {
			return false
		}
	}
	return true
}
*/
