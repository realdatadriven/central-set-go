package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/realdatadriven/central-set-go/internal/password"

	"github.com/pascaldekloe/jwt"
	"github.com/realdatadriven/etlx"
)

// AuthenticateAD connects to AD and validates a user's credentials.
func AuthenticateAD(ldapURL, baseDN, serviceUser, servicePass, username, password string, skipVerify bool) (bool, Dict, error) {
	//fmt.Println("AuthenticateAD:", ldapURL, baseDN, serviceUser, servicePass, username, password)
	// Connect to LDAP server (use LDAPS if possible)
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: skipVerify}))
	if err != nil {
		return false, nil, fmt.Errorf("failed to connect LDAP: %w", err)
	}
	defer l.Close()
	// 1 Bind as service account (to search)
	err = l.Bind(serviceUser, servicePass)
	if err != nil {
		return false, nil, fmt.Errorf("service bind failed: %w", err)
	}
	// 2 Search for user by sAMAccountName, userPrincipalName, or mail
	//filter := fmt.Sprintf("(|(sAMAccountName=%[1]s)(userPrincipalName=%[1]s)(mail=%[1]s))", username)
	filter := fmt.Sprintf("(|(uid=%[1]s)(cn=%[1]s)(mail=%[1]s))", username)
	if os.Getenv("LDAP_SEARCHREQ_FILTER") != "" {
		filter = fmt.Sprintf(os.Getenv("LDAP_SEARCHREQ_FILTER"), username)
	}
	//fmt.Println(baseDN, 1, os.Getenv("LDAP_SEARCHREQ_FILTER"), 2, filter)
	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", "cn", "mail", "displayName", "givenName", "sn", "sAMAccountName", "userPrincipalName", "memberOf"},
		nil,
	)
	sr, err := l.Search(searchReq)
	if err != nil {
		return false, nil, fmt.Errorf("LDAP search failed: %w", err)
	}
	if len(sr.Entries) != 1 {
		return false, nil, fmt.Errorf("user not found or multiple matches for %s", username)
	}
	entry := sr.Entries[0]
	userDN := entry.DN
	//fmt.Println("LDAP search entries:", len(sr.Entries), userDN, password)
	// 3 Verify password by binding as user
	if err := l.Bind(userDN, password); err != nil {
		return false, nil, fmt.Errorf("invalid credentials %s: %v", username, err) // invalid credentials
	}
	displayName := entry.GetAttributeValue("displayName")
	firstName := entry.GetAttributeValue("givenName")
	lastName := entry.GetAttributeValue("sn")
	// Fallback: split displayName if givenName/sn not available
	if firstName == "" && displayName != "" {
		parts := strings.Fields(displayName)
		if len(parts) > 0 {
			firstName = parts[0]
		}
		if len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}
	}
	userInfo := make(Dict)
	email := entry.GetAttributeValue("mail")
	if entry.GetAttributeValue("mail") == "" {
		email = username
	}
	// 4 Extract user attributes into a map
	userInfo["dn"] = userDN
	userInfo["cn"] = entry.GetAttributeValue("cn")
	userInfo["displayName"] = entry.GetAttributeValue("displayName")
	userInfo["givenName"] = entry.GetAttributeValue("givenName")
	userInfo["sn"] = entry.GetAttributeValue("sn")
	userInfo["first_name"] = firstName
	userInfo["last_name"] = lastName
	userInfo["username"] = username
	userInfo["email"] = email
	userInfo["upn"] = entry.GetAttributeValue("userPrincipalName")
	userInfo["groups"] = entry.GetAttributeValues("memberOf")
	return true, userInfo, nil
}

// dynamic login function where a table cudo be set, the username / email is set and the passwor field to like for exemple, i have a table with name tennant, has field tennant = username and has email = email and has field password = password
func (app *application) dynamic_login(params Dict) Dict {
	login_table := "users"
	if _, ok := params["login_table"].(string); ok {
		login_table = params["login_table"].(string)
	}
	if login_table == "" {
		msg, _ := app.i18n.T("login-table-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	user_id_field := "user_id"
	if _, ok := params["user_id_field"].(string); ok {
		user_id_field = params["user_id_field"].(string)
	}
	if user_id_field == "" {
		msg, _ := app.i18n.T("user-id-field-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	username_field := "username"
	if _, ok := params["username_field"].(string); ok {
		username_field = params["username_field"].(string)
	}
	if username_field == "" {
		msg, _ := app.i18n.T("username-field-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	password_field := "password"
	if _, ok := params["password_field"].(string); ok {
		password_field = params["password_field"].(string)
	}
	if password_field == "" {
		msg, _ := app.i18n.T("password-field-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	email_field := "email"
	if _, ok := params["email_field"].(string); ok {
		email_field = params["email_field"].(string)
	}
	if email_field == "" {
		msg, _ := app.i18n.T("email-field-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	active_field := "active"
	if _, ok := params["active_field"].(string); ok {
		active_field = params["active_field"].(string)
	}
	_data := Dict{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(Dict)
	}
	if app.IsEmpty(_data) {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	username := ""
	if _, ok := _data[username_field].(string); ok {
		username = _data[username_field].(string)
	} else if _, ok := _data[email_field].(string); ok {
		username = _data[email_field].(string)
	} else if _, ok := _data["username"].(string); ok {
		username = _data["username"].(string)
	} else if _, ok := _data["user"].(string); ok {
		username = _data["user"].(string)
	} else if _, ok := _data["email"].(string); ok {
		username = _data["email"].(string)
	}
	pass := ""
	if _, ok := _data[password_field].(string); ok {
		pass = _data[password_field].(string)
	} else if _, ok := _data["password"].(string); ok {
		pass = _data["password"].(string)
	} else if _, ok := _data["pass"].(string); ok {
		pass = _data["pass"].(string)
	}
	var user Dict
	//var found bool
	var err error
	sql := fmt.Sprintf(`select * from "%s" where ("%s" = ? OR "%s" = ?) and "%s" = true and excluded = false`, login_table, username_field, email_field, active_field)
	user, err = app.GetRowByFilter(sql, params, []any{username, username})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	} else if len(user) == 0 {
		msg, _ := app.i18n.T("user-pass-incorrect", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	if len(user) > 0 {
		//_hash, _ := password.Hash(pass)
		//fmt.Println(pass, _hash, user["password"].(string))
		match, err := password.Matches(pass, user[password_field].(string))
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		if !match {
			msg, _ := app.i18n.T("user-pass-incorrect", Dict{})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
	}
	delete(user, password_field)
	user["user_id"] = user[user_id_field]
	user["username"] = user[username_field]
	user["email"] = user[email_field]
	user["role_id"] = params["dyn_login_role_id"]
	user["login_table"] = login_table
	user["user_id_field"] = user_id_field
	user["username_field"] = username_field
	user["password_field"] = password_field
	user["email_field"] = email_field
	user["is_dynamic"] = true
	delete(user, "created_at")
	delete(user, "updated_at")
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	claims.Subject = string(json_user)
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	data := Dict{
		"success": true,
		"msg":     "Loged in successfully!",
		"data":    user,
		"token":   string(jwtBytes),
		"expiry":  expiry.Format(time.RFC3339),
	}
	return data
}

// signup function, gets username/email and password from params, creates a new user
func (app *application) signup(params Dict) Dict {
	msg, _ := app.i18n.T("not-implemented-yet", Dict{})
	return Dict{
		"success": false,
		"msg":     msg,
	}
}

// login function, gets username/email and password from params, validates user and returns JWT token

func (app *application) _login(params Dict) Dict {
	_data := Dict{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(Dict)
	}
	if app.IsEmpty(_data) {
		msg, _ := app.i18n.T("no-data", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	username := ""
	if _, ok := _data["username"].(string); ok {
		username = _data["username"].(string)
	} else if _, ok := _data["user"].(string); ok {
		username = _data["user"].(string)
	} else if _, ok := _data["u"].(string); ok {
		username = _data["u"].(string)
	}
	pass := ""
	if _, ok := _data["password"].(string); ok {
		pass = _data["password"].(string)
	} else if _, ok := _data["pass"].(string); ok {
		pass = _data["pass"].(string)
	} else if _, ok := _data["p"].(string); ok {
		pass = _data["p"].(string)
	}
	var user Dict
	var found bool
	var err error
	// fmt.Println(_data)
	if os.Getenv("USE_LDAP_AUTH") == "true" && username != "root" {
		ldapURL := os.Getenv("LDAP_URL")
		serviceUser := os.Getenv("LDAP_BIND_USER")
		servicePass := os.Getenv("LDAP_PASSWORD")
		skipVerify := false
		if os.Getenv("LDAP_SKIP_VERIFY_CERT") == "true" {
			skipVerify = true
		}
		baseDN := os.Getenv("LDAP_BASE_DN")
		_, userInfo, err := AuthenticateAD(ldapURL, baseDN, serviceUser, servicePass, username, pass, skipVerify)
		//fmt.Println("LDAP", userInfo, err)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		user, found, err = app.db.GetUserByNameOrEmail(username)
		if err != nil || len(user) == 0 {
			user = Dict{
				"username":   userInfo["username"],
				"password":   "LDAP_USER", //userInfo["password"],
				"first_name": userInfo["first_name"],
				"last_name":  userInfo["last_name"],
				"email":      userInfo["email"],
				"lang_id":    1,
				"role_id":    2,
				"active":     true,
				"excluded":   false,
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}
			err := app.AdminInsertData("users", user)
			if err != nil {
				return Dict{
					"success": false,
					"msg":     err.Error(),
				}
			} else {
				user, found, err = app.db.GetUserByNameOrEmail(username)
				if err != nil {
					return Dict{
						"success": false,
						"msg":     err.Error(),
					}
				}
			}
		}
	} else {
		user, found, err = app.db.GetUserByNameOrEmail(username)
		if err != nil || len(user) == 0 {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		if found {
			//_hash, _ := password.Hash(pass)
			//fmt.Println(pass, _hash, user["password"].(string))
			match, err := password.Matches(pass, user["password"].(string))
			if err != nil {
				return Dict{
					"success": false,
					"msg":     err.Error(),
				}
			}
			if !match {
				msg, _ := app.i18n.T("user-pass-incorrect", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
		}
	}
	delete(user, "password")
	delete(user, "created_at")
	delete(user, "updated_at")
	delete(user, "phone")
	delete(user, "email")
	delete(user, "timezone")
	delete(user, "attach_profile_pic")
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	claims.Subject = string(json_user)
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	data := Dict{
		"success": true,
		"msg":     "Loged in successfully!",
		"data":    user,
		"token":   string(jwtBytes),
		"expiry":  expiry.Format(time.RFC3339),
	}
	return data
}

func (app *application) alter_pass(params Dict) Dict {
	_check_login := app._login(params)
	if success, ok := _check_login["success"].(bool); ok && success {
		_data := Dict{}
		if _, ok := params["data"]; ok {
			_data = params["data"].(Dict)
		}
		newPassword, _ := _data["new_password"].(string)
		oldPassword, _ := _data["password"].(string)
		if newPassword == "" {
			msg, _ := app.i18n.T("new_pass_is_required", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		if newPassword == oldPassword {
			msg, _ := app.i18n.T("new_pass_old_pass", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		if len(newPassword) < 8 {
			msg, _ := app.i18n.T("password_min_length", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		hasUpper, _ := regexp.MatchString(`[A-Z]`, newPassword)
		if !hasUpper {
			msg, _ := app.i18n.T("pass_must_have_upper", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		hasNumber, _ := regexp.MatchString(`[0-9]`, newPassword)
		if !hasNumber {
			msg, _ := app.i18n.T("pass_must_have_number", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		hasSpecial, _ := regexp.MatchString(`[$&+,:;=?@#!*ªº.-]`, newPassword)
		if !hasSpecial {
			msg, _ := app.i18n.T("pass_must_have_special", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		query := `UPDATE users 
			SET password = :password 
		WHERE email = :username
			OR username = :username`
		pass, err := password.Hash(newPassword)
		if err != nil {
			msg, _ := app.i18n.T("pass_must_have_special", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		username := ""
		if _, ok := _data["username"].(string); ok {
			username = _data["username"].(string)
		} else if _, ok := _data["user"].(string); ok {
			username = _data["user"].(string)
		} else if _, ok := _data["u"].(string); ok {
			username = _data["u"].(string)
		}
		_data = Dict{"username": username, "password": pass}
		_, err = app.db.ExecuteNamedQuery(query, _data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
		msg, _ := app.i18n.T("alter-pass-success", Dict{})
		return Dict{
			"success": true,
			"msg":     msg,
		}
	} else {
		return _check_login
	}
}

func (app *application) access_key(params Dict) Dict {
	_data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		if _, ok := params["data"].(Dict)["data"]; !ok {
			msg, _ := app.i18n.T("no-data", Dict{})
			return Dict{"success": false, "msg": msg}
		}
		_data = params["data"].(Dict)["data"].(Dict)
	}
	// fmt.Println(_data)
	if _, ok := _data["access_key_desc"].(string); !ok {
		msg, _ := app.i18n.T("access-key-no-description", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	if _, ok := _data["expires_at"].(string); !ok {
		msg, _ := app.i18n.T("access-key-no-expires-dt", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	if _, ok := _data["for_user_id"]; !ok {
		msg, _ := app.i18n.T("access-key-no-user", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	sql := `select * from "users" where "user_id" = ? and  "excluded" = false`
	user, err := app.AdminGetRowByID(sql, _data["for_user_id"])
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	if len(user) == 0 {
		msg, _ := app.i18n.T("user-incorrect", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	if _, ok := user["user_id"]; !ok {
		msg, _ := app.i18n.T("user-incorrect", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	//time.Parse("", )
	var expires_at time.Time
	expires_at, err = time.Parse("2006-01-02T15:04:05", _data["expires_at"].(string))
	if err != nil {
		expires_at, err = time.Parse("2006-01-02 15:04:05", _data["expires_at"].(string))
		if err != nil {
			expires_at, err = time.Parse("2006-01-02 15:04", _data["expires_at"].(string))
			if err != nil {
				msg, _ := app.i18n.T("incorret-expiration-dt", Dict{"dt": _data["expires_at"].(string)})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
		}
	}
	delete(user, "password")
	delete(user, "created_at")
	delete(user, "updated_at")
	delete(user, "phone")
	delete(user, "email")
	delete(user, "timezone")
	delete(user, "attach_profile_pic")
	// fmt.Println("EXPIRATION DATE:", _data["expires_at"].(string))
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	claims.Subject = string(json_user)
	expiry := time.Now().Add(time.Until(expires_at))
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	_data["access_token"] = string(jwtBytes)
	params["data"].(Dict)["data"] = _data
	upsert := app.create_update(params)
	// fmt.Println(upsert)
	if _, ok := upsert["success"]; !ok {
		return upsert
	} else if _, ok := upsert["success"].(bool); !ok {
		return upsert
	} else if ok, _ := upsert["success"].(bool); !ok {
		return upsert
	}
	// fmt.Println(expiry.Format(time.RFC3339), _data["access_token"])
	msg, _ := app.i18n.T("success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
		"expiry":  expiry.Format(time.RFC3339),
	}
}

// recover password, gets params with email, checks if is a valid user, creates a token and send email with link to reset password
func (app *application) recover_pass(params Dict) Dict {
	_data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		_data = params["data"].(Dict)
	}
	email := ""
	if _, ok := _data["email"].(string); ok {
		email = _data["email"].(string)
	}
	if email == "" {
		msg, _ := app.i18n.T("email-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	user, found, err := app.db.GetUserByNameOrEmail(email)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	if !found || len(user) == 0 {
		msg, _ := app.i18n.T("user-not-found", Dict{"email": email})
		fmt.Println("USER NOT FOUND!", msg, email)
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	delete(user, "password")
	delete(user, "created_at")
	delete(user, "updated_at")
	delete(user, "phone")
	delete(user, "timezone")
	delete(user, "attach_profile_pic")
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	claims.Subject = string(json_user)
	expiry := time.Now().Add(1 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.frontend_url}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", app.config.frontend_url, string(jwtBytes))
	_etlx := etlx.ETLX{}
	bodyTemplate := `
		<p>Hi {{.first_name}},</p>
		<p>You have requested to reset your password. Please click the link below to reset your password:</p>
		<p><a href="{{.reset_link}}">{{.reset_link}}</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request a password reset, please ignore this email.</p>
		<p>Best regards,<br/>The Team</p>
	`
	emailParams := Dict{
		"to":      []any{user["email"]},
		"subject": "Password Recovery",
		"body":    bodyTemplate,
		"data": Dict{
			"first_name": user["first_name"],
			"reset_link": resetLink,
		},
	}
	err = _etlx.SendEmail(emailParams)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	msg, _ := app.i18n.T("recover-pass-email-sent", Dict{"email": email})
	return Dict{
		"success": true,
		"msg":     msg,
	}
}

// reset password, gets params with token and new_password, validates token and updates user password
func (app *application) reset_pass(params Dict) Dict {
	_data := Dict{}
	//fmt.Println("reset_pass params:", params)
	if _, ok := params["data"].(Dict); ok {
		_data = params["data"].(Dict)
	}
	tokenStr := ""
	if _, ok := _data["token"].(string); ok {
		tokenStr = _data["token"].(string)
	}
	newPassword := ""
	if _, ok := _data["new_password"].(string); ok {
		newPassword = _data["new_password"].(string)
	}
	confirmPassword := ""
	if _, ok := _data["confirm_password"].(string); ok {
		confirmPassword = _data["confirm_password"].(string)
	}
	if newPassword != confirmPassword {
		msg, _ := app.i18n.T("new_pass_diff_confirm_pass", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	if len(newPassword) < 8 {
		msg, _ := app.i18n.T("password_min_length", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	hasUpper, _ := regexp.MatchString(`[A-Z]`, newPassword)
	if !hasUpper {
		msg, _ := app.i18n.T("pass_must_have_upper", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	hasNumber, _ := regexp.MatchString(`[0-9]`, newPassword)
	if !hasNumber {
		msg, _ := app.i18n.T("pass_must_have_number", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	hasSpecial, _ := regexp.MatchString(`[$&+,:;=?@#!*ªº.-]`, newPassword)
	if !hasSpecial {
		msg, _ := app.i18n.T("pass_must_have_special", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	//fmt.Println(tokenStr)
	if tokenStr == "" || newPassword == "" {
		msg, _ := app.i18n.T("token-and-new-pass-required", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	token, err := jwt.HMACCheck([]byte(tokenStr), []byte(app.config.jwt.secretKey))
	if err != nil {
		msg, _ := app.i18n.T("invalid-token", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	if token.Expires.Time().Before(time.Now()) {
		msg, _ := app.i18n.T("token-expired", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	subject := token.Subject
	var user Dict
	err = json.Unmarshal([]byte(subject), &user)
	if err != nil {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	if _, ok := user["username"]; !ok {
		msg, _ := app.i18n.T("invalid-token-no-username", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	query := `UPDATE users 
			SET password = :password 
		WHERE email = :username
			OR username = :username`
	pass, err := password.Hash(newPassword)
	if err != nil {
		msg, _ := app.i18n.T("password-hash-error", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	data := Dict{"username": user["username"], "password": pass}
	_, err = app.db.ExecuteNamedQuery(query, data)
	if err != nil {
		msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	msg, _ := app.i18n.T("reset-pass-success", Dict{})
	return Dict{
		"success": true,
		"msg":     msg,
	}
}

// oauth_login handles  google, github, etc. oauth login
func (app *application) oauth_login(params Dict) Dict {
	msg, _ := app.i18n.T("not-implemented-yet", Dict{})
	return Dict{
		"success": false,
		"msg":     msg,
	}
}
