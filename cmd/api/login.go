package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/markbates/goth/gothic"
	"github.com/realdatadriven/central-set-go/assets"
	"github.com/realdatadriven/central-set-go/internal/env"
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
			if app.config.LockoutEnabled {
				_ = app.updateFailedLoginAttempts(username)
			}
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
func (app *application) dynamic_signup(params Dict) Dict {
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
	email := ""
	if _, ok := _data[email_field].(string); ok {
		email = _data[email_field].(string)
	} else if _, ok := _data["email"].(string); ok {
		email = _data["email"].(string)
	}
	pass := ""
	if _, ok := _data[password_field].(string); ok {
		pass = _data[password_field].(string)
	} else if _, ok := _data["password"].(string); ok {
		pass = _data["password"].(string)
	} else if _, ok := _data["pass"].(string); ok {
		pass = _data["pass"].(string)
	}
	query := `INSERT INTO %s (%s, %s, %s) VALUES (:username, :email, :password)`
	password_hashed, err := password.Hash(pass)
	if err != nil {
		msg, _ := app.i18n.T("password-hash-error", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	data := Dict{}
	data[username_field] = username
	data[email_field] = email
	data[password_field] = password_hashed
	data[active_field] = true
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
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

// update failed_login_attmpt and last_failed_login for table users
func (app *application) updateFailedLoginAttempts(username string) error {
	query := `UPDATE users 
		SET failed_login_attmpt = coalesce(failed_login_attmpt, 0) + 1, 
		last_failed_login = :last_failed_login,
		active = CASE
			WHEN coalesce(failed_login_attmpt, 0) + 1 >= :lockout_threshold THEN false
			ELSE true
		END,
		updated_at = :last_failed_login
	WHERE username = :username`
	data := Dict{
		"username":          username,
		"last_failed_login": time.Now(),
		"lockout_threshold": app.config.LockoutThreshold,
	}
	_, err := app.db.ExecuteNamedQuery(query, data)
	return err
}

// reset failed_login_attmpt and last_failed_login for table users
func (app *application) resetFailedLoginAttempts(username string) error {
	query := `UPDATE users
		SET failed_login_attmpt = 0, 
		last_failed_login = NULL,
		updated_at = :updated_at
	WHERE username = :username`
	data := Dict{
		"username":   username,
		"updated_at": time.Now(),
	}
	_, err := app.db.ExecuteNamedQuery(query, data)
	return err
}

// login function, gets username/email and password from params, validates user and returns JWT token
func (app *application) _login(params Dict) Dict {
	_data := Dict{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(Dict)
	} else {
		return Dict{
			"success": false,
			"msg":     "No \"data\" parameter provided in your message body!",
		}
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
	} else if _, ok := _data["email"].(string); ok {
		username = _data["email"].(string)
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
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		//fmt.Println(user)
		if len(user) == 0 {
			msg, _ := app.i18n.T("user-pass-incorrect", Dict{})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
		// check if active is true
		if active, ok := user["active"].(bool); ok {
			if !active {
				msg, _ := app.i18n.T("account-inactive", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
		} else if active, ok := user["active"]; ok {
			if app.toBool(active) == false {
				msg, _ := app.i18n.T("account-inactive", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
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
				if app.config.LockoutEnabled {
					// fmt.Println("app.updateFailedLoginAttempts(username):", username, app.config.LockoutThreshold)
					err = app.updateFailedLoginAttempts(username)
					if err != nil {
						fmt.Println("Error updating failed login attempts:", err)
					}
				}
				msg, _ := app.i18n.T("user-pass-incorrect", Dict{})
				return Dict{
					"success": false,
					"msg":     msg,
				}
			}
		}
	}
	// TWO FACTOR
	if enable_2f_auth, ok := user["enable_2f_auth"].(bool); ok && enable_2f_auth {
		code := app.randomString(8)
		nxt_code_2f_auth, _ := password.Hash(code)
		_etlx := etlx.ETLX{}
		bodyTemplate := `
			<p>Hi {{.first_name}},</p>
			<p>You here is the two factor code generate in your login: {{.code}}</p>
			<p>It will expire in 5 min.</p>
			<p>If you did not request a login, please ignore this email.</p>
			<p>Best regards,<br/>The Team</p>
		`
		tmplPath := "templates/two-factor-code.html"
		if os.Getenv("TWO_FACTOR_CODE_TMPL_PATH") != "" {
			tmplPath = os.Getenv("TWO_FACTOR_CODE_TMPL_PATH")
		}
		content, err := os.ReadFile(tmplPath)
		if err != nil {
			content, err = assets.EmbeddedFiles.ReadFile(tmplPath)
		}
		if content != nil {
			bodyTemplate = string(content)
		}
		emailParams := Dict{
			"to":      []any{user["email"]},
			"subject": "Two Factor Auth",
			"body":    bodyTemplate,
			"data": Dict{
				"first_name": user["first_name"],
				"code":       code,
			},
		}
		err = _etlx.SendEmail(emailParams)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		query := `UPDATE users 
			SET nxt_code_2f_auth = :nxt_code_2f_auth
			  , code_2f_expires_at = :code_2f_expires_at
			  , updated_at = :updated_at
		WHERE user_id = :user_id`
		min_2_expire := env.GetInt("TWO_FACTOR_CODE_EXP_IN_MIN", 5)
		_data = Dict{
			"user_id":            user["user_id"],
			"nxt_code_2f_auth":   nxt_code_2f_auth,
			"code_2f_expires_at": time.Now().Add(time.Duration(min_2_expire) * time.Minute),
			"updated_at":         time.Now(),
		}
		_, err = app.db.ExecuteNamedQuery(query, _data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
		msg, _ := app.i18n.T("two-factor-code-email-sent", Dict{"email": user["email"]})
		return Dict{
			"success":    true,
			"msg":        msg,
			"two_factor": true,
			"username":   username,
		}
	}
	// reset failed login attempts
	if app.config.LockoutEnabled {
		_ = app.resetFailedLoginAttempts(username)
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

// login function, gets username/email and password from params, validates user and returns JWT token
func (app *application) two_factor_code_valid(params Dict) Dict {
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
	if _, ok := _data["code"].(string); ok {
		pass = _data["code"].(string)
	}
	var user Dict
	var found bool
	//var err Error
	user, found, err := app.db.GetUserByNameOrEmail(username)
	if err != nil || len(user) == 0 {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	if !found || len(user) == 0 {
		msg, _ := app.i18n.T("user-not-found", Dict{"email": user["email"]})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	if _, ok := user["nxt_code_2f_auth"].(string); !ok {
		msg, _ := app.i18n.T("code-not-saved-or-already-used", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	} else if cd, ok := user["nxt_code_2f_auth"].(string); ok && cd == "" {
		msg, _ := app.i18n.T("code-not-saved-or-already-used", Dict{})
		return Dict{
			"success": false,
			"msg":     msg,
		}
	}
	if found {
		//_hash, _ := password.Hash(pass)
		//fmt.Println(pass, _hash, user["password"].(string))
		match, err := password.Matches(pass, user["nxt_code_2f_auth"].(string))
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		if !match {
			if app.config.LockoutEnabled {
				_ = app.updateFailedLoginAttempts(username)
			}
			msg, _ := app.i18n.T("two-factor-code-incorrect", Dict{})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
		code_2f_expires_at, _ := user["code_2f_expires_at"].(time.Time)
		if code_2f_expires_at.Before(time.Now()) {
			msg, _ := app.i18n.T("two-factor-code-expired", Dict{})
			return Dict{
				"success": false,
				"msg":     msg,
			}
		}
		query := `UPDATE users 
			SET nxt_code_2f_auth = :nxt_code_2f_auth
			  , code_2f_expires_at = :code_2f_expires_at
			  , updated_at = :updated_at
		WHERE user_id = :user_id`
		_data = Dict{
			"user_id":            user["user_id"],
			"nxt_code_2f_auth":   nil,
			"code_2f_expires_at": nil,
			"updated_at":         time.Now(),
		}
		_, err = app.db.ExecuteNamedQuery(query, _data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
			return Dict{
				"success": false,
				"msg":     msg,
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
				, alter_pass_nxt_login = false
				, updated_at = :updated_at
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
		_data = Dict{"username": username, "password": pass, "updated_at": time.Now()}
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
	_table := params["data"].(Dict)["table"]
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
		dyn_login_map_2users := env.GetBool("DYN_LOGIN_TABLE_MAP_TO_USERS", false)
		dyn_login_table := env.GetString("DYN_LOGIN_TABLE", "")
		dyn_login_table_id := env.GetString("DYN_LOGIN_USER_ID_FIELD", "")
		dyn_login_table_email := env.GetString("DYN_LOGIN_EMAIL_FIELD", "")
		dyn_login_table_active := env.GetString("DYN_LOGIN_ACTIVE_FIELD", "")
		dyn_login_id_ok := false
		dyn_id := any(nil)
		if _, ok := _data[dyn_login_table_id]; ok {
			dyn_login_id_ok = true
			dyn_id = any(app.toInt(_data[dyn_login_table_id]))
		}
		if dyn_login_map_2users && dyn_login_table != "" && dyn_login_table_id != "" && dyn_login_table_email != "" && dyn_login_id_ok {
			dyn_login_data, err := app.ODataGetRow(params, dyn_login_table, dyn_login_table_id, dyn_id)
			if err != nil {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting %s/%d: %s", dyn_login_table, dyn_id, err.Error()),
				}
			}
			if len(dyn_login_data) == 0 {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error, no data returned for %s/%d", dyn_login_table, dyn_id),
				}
			}
			if email, ok := dyn_login_data[dyn_login_table_email].(string); !ok || email == "" {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error, no email returned for %s/%d, please contact the admin", dyn_login_table, dyn_id),
				}
			}
			if active, ok := dyn_login_data[dyn_login_table_active]; ok {
				if !app.toBool(active) {
					return Dict{
						"success": false,
						"msg":     fmt.Sprintf("%s/%s: does not seem to be active, please contact the admin", dyn_login_table, dyn_login_data[dyn_login_table_email]),
					}
				}
			} else {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("%s/%s: does not seem to be active, please contact the admin", dyn_login_table, dyn_login_data[dyn_login_table_email]),
				}
			}
			sql := `select * from users where email = ?`
			dyn_login_eqv_user, err := app.AdminGetRowByFilter(sql, []any{dyn_login_data[dyn_login_table_email]})
			if err != nil {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error getting the user for %s/%s: %s", dyn_login_table, dyn_login_data[dyn_login_table_email], err.Error()),
				}
			}
			if len(dyn_login_eqv_user) == 0 {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("Error, no user data returned for %s/%s", dyn_login_table, dyn_login_data[dyn_login_table_email]),
				}
			}
			if active, ok := dyn_login_eqv_user["active"]; ok {
				if !app.toBool(active) {
					return Dict{
						"success": false,
						"msg":     fmt.Sprintf("%s/%s: user's does not seem to be active, please contact the admin", dyn_login_table, dyn_login_data[dyn_login_table_email]),
					}
				}
			} else {
				return Dict{
					"success": false,
					"msg":     fmt.Sprintf("%s/%s: users does not seem to be active, please contact the admin", dyn_login_table, dyn_login_data[dyn_login_table_email]),
				}
			}
			_data["for_user_id"] = dyn_login_eqv_user["user_id"]
			params["data"].(Dict)["table"] = _table
		} else {
			msg, _ := app.i18n.T("access-key-no-user", Dict{})
			return Dict{"success": false, "msg": msg}
		}
	}
	sql := `select * from "users" where "user_id" = ? and "excluded" = false`
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
	_, db, _ := app.GetDBNameFromParams(params)
	bdc := Dict{
		"type":     "data_change",
		"database": db,
		"table":    params["data"].(Dict)["table"],
	}
	app.BroadCastChange(bdc)
	//fmt.Println(expiry.Format(time.RFC3339), _data["access_token"])
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
		// fmt.Println("USER NOT FOUND!", msg, email)
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
	tmplPath := "templates/recover-pass.html"
	if os.Getenv("RECOVER_PASS_TMPL_PATH") != "" {
		tmplPath = os.Getenv("RECOVER_PASS_TMPL_PATH")
	}
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		content, err = assets.EmbeddedFiles.ReadFile(tmplPath)
	}
	if content != nil {
		bodyTemplate = string(content)
	}
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
		"success":      true,
		"msg":          msg,
		"recover_pass": true,
	}
}

func (app *application) confirm_emmail(params Dict) Dict {
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
	resetLink := fmt.Sprintf("%s/confirm-email?token=%s", app.config.frontend_url, string(jwtBytes))
	_etlx := etlx.ETLX{}
	bodyTemplate := `
		<p>Hi {{.first_name}},</p>
		<p>You have requested to confirm your email. Please click the link below to confirm your email:</p>
		<p><a href="{{.reset_link}}">{{.reset_link}}</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request a user account, please ignore this email.</p>
		<p>Best regards,<br/>The Team</p>
	`
	tmplPath := "templates/confirm-email.html"
	if os.Getenv("CONFIRM_EMAIL_TMPL_PATH") != "" {
		tmplPath = os.Getenv("CONFIRM_EMAIL_TMPL_PATH")
	}
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		content, err = assets.EmbeddedFiles.ReadFile(tmplPath)
	}
	if content != nil {
		bodyTemplate = string(content)
	}
	emailParams := Dict{
		"to":      []any{user["email"]},
		"subject": "Confirm Email",
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
	msg, _ := app.i18n.T("confirm-email-sent", Dict{"email": email})
	return Dict{
		"success":       true,
		"msg":           msg,
		"confirm_email": true,
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
				, alter_pass_nxt_login = false
				, updated_at = :updated_at
		WHERE email = :username
			OR username = :username`
	pass, err := password.Hash(newPassword)
	if err != nil {
		msg, _ := app.i18n.T("password-hash-error", Dict{})
		return Dict{"success": false, "msg": msg}
	}
	data := Dict{"username": user["username"], "password": pass, "updated_at": time.Now()}
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

// OAUTH 2
func (app *application) GothLoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider required", http.StatusBadRequest)
		return
	}

	// Goth handles state, redirect, PKCE (for supported providers), etc.
	// You can pass extra oauth2 options if needed
	q := r.URL.Query()
	if q.Get("prompt") != "" {
		// example: force consent screen
		//r = gothic.SetState(r, "state-with-prompt") // optional custom state
	}

	url, err := gothic.GetAuthURL(w, r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*for _, c := range w.Header().Values("Set-Cookie") {
		fmt.Println("gothic.GetAuthURL Cookie:", c)
	}*/
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (app *application) GothCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider missing from path", http.StatusBadRequest)
		return
	}
	// fmt.Println("GothCallbackHandler Cookie:", r.Header.Get("Cookie"))
	//fmt.Println("CALLBACK PROVIDER:", provider)
	// Completes the flow: exchanges code, fetches user info
	gu, err := gothic.CompleteUserAuth(w, r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider)))
	if err != nil {
		fmt.Println(11, "gothic.CompleteUserAuth:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(gu.Name, gu.Email, gu.FirstName, gu.LastName, gu.ExpiresAt)

	user, found, err := app.db.GetUserByNameOrEmail(gu.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": err.Error()})
		return
	}
	if !found || len(user) == 0 {
		query := `INSERT INTO 
		"users" ("username", "first_name" , "last_name", "email", "password", "role_id", "lang_id", "active", "alter_pass_nxt_login", "created_at", "updated_at", "excluded") 
		VALUES (:username, :first_name, :last_name, :email, :password, :role_id, :lang_id, :active, :alter_pass_nxt_login, :created_at, :updated_at, :excluded)`
		pass, err := password.Hash(app.randomString(8))
		if err != nil {
			msg, _ := app.i18n.T("password-hash-error", Dict{})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": msg})
			return
		}
		role_id := env.GetInt("OAUTH_DEFAULT_ROLE_ID", 2)
		username := gu.Email
		if len(strings.Split(gu.Email, "@")) > 1 {
			username = strings.Split(gu.Email, "@")[0]
		}
		first_name := gu.FirstName
		if first_name == "" && len(strings.Split(gu.Name, " ")) > 1 {
			first_name = strings.Split(gu.Name, " ")[0]
		} else if first_name == "" {
			first_name = gu.Name
		}
		last_name := gu.LastName
		if last_name == "" && len(strings.Split(gu.Name, " ")) > 1 {
			last_name = strings.Split(gu.Name, " ")[1]
		} else if last_name == "" {
			last_name = gu.Name
		}
		data := Dict{
			"username":             username,
			"first_name":           first_name,
			"last_name":            last_name,
			"email":                gu.Email,
			"password":             pass,
			"role_id":              role_id,
			"lang_id":              1,
			"active":               true,
			"alter_pass_nxt_login": false,
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
			"excluded":             false,
		}
		_, err = app.db.ExecuteNamedQuery(query, data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": msg})
			return
		}
		user, found, err = app.db.GetUserByNameOrEmail(gu.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": err.Error()})
			return
		} else if !found || len(user) == 0 {
			msg, _ := app.i18n.T("user-not-found", Dict{"email": gu.Name})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": msg})
			return
		}
	}
	user["oauth"] = true
	user["provider"] = provider
	delete(user, "password")
	delete(user, "created_at")
	delete(user, "updated_at")
	delete(user, "phone")
	delete(user, "timezone")
	delete(user, "attach_profile_pic")
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": err.Error()})
		return
	}
	claims.Subject = string(json_user)
	//expiry := time.Now().Add(8 * time.Hour)
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	if !gu.ExpiresAt.IsZero() {
		// expiry = gu.ExpiresAt
	}
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.frontend_url}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    string(jwtBytes),
		HttpOnly: false,
		Secure:   false, // true
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	w.Header().Set("X-Auth-Provider", provider)      // ← optional
	w.Header().Set("X-Auth-Token", string(jwtBytes)) // ← optional
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// Clone of GothCallbackHandler but for hypermedia response insted of json, to work with things like htmx
func (app *application) HyperMGothCallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	provider := r.PathValue("provider")
	// fmt.Println("GothCallbackHandler Cookie:", r.Header.Get("Cookie"))
	if provider == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString("provider missing from path"))
		return
	}
	// fmt.Println("CALLBACK PROVIDER:", provider, gothic.ProviderParamKey)
	// Completes the flow: exchanges code, fetches user info
	gu, err := gothic.CompleteUserAuth(w, r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider)))
	if err != nil {
		fmt.Println(12, "gothic.CompleteUserAuth:", provider, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString(err.Error()))
		return
	}

	loc := app.getLocationFromRequest(r, Dict{})

	//fmt.Println(gu.Name, gu.Email, gu.FirstName, gu.LastName, gu.ExpiresAt)
	user, found, err := app.db.GetUserByNameOrEmail(gu.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString(err.Error()))
		return
	}
	if !found || len(user) == 0 {
		query := `INSERT INTO 
		"users" ("username", "first_name" , "last_name", "email", "password", "role_id", "lang_id", "active", "alter_pass_nxt_login", "created_at", "updated_at", "excluded") 
		VALUES (:username, :first_name, :last_name, :email, :password, :role_id, :lang_id, :active, :alter_pass_nxt_login, :created_at, :updated_at, :excluded)`
		pass, err := password.Hash(app.randomString(8))
		if err != nil {
			msg, _ := app.i18n.T("password-hash-error", Dict{})
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `%s`, template.HTMLEscapeString(msg))
			return
		}
		role_id := env.GetInt("OAUTH_DEFAULT_ROLE_ID", 2)
		username := gu.Email
		if len(strings.Split(gu.Email, "@")) > 1 {
			username = strings.Split(gu.Email, "@")[0]
		}
		first_name := gu.FirstName
		if first_name == "" && len(strings.Split(gu.Name, " ")) > 1 {
			first_name = strings.Split(gu.Name, " ")[0]
		} else if first_name == "" {
			first_name = gu.Name
		}
		last_name := gu.LastName
		if last_name == "" && len(strings.Split(gu.Name, " ")) > 1 {
			last_name = strings.Split(gu.Name, " ")[1]
		} else if last_name == "" {
			last_name = gu.Name
		}
		data := Dict{
			"username":             username,
			"first_name":           first_name,
			"last_name":            last_name,
			"email":                gu.Email,
			"password":             pass,
			"role_id":              role_id,
			"lang_id":              1,
			"active":               true,
			"alter_pass_nxt_login": false,
			"created_at":           time.Now().In(loc),
			"updated_at":           time.Now().In(loc),
			"excluded":             false,
		}
		_, err = app.db.ExecuteNamedQuery(query, data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", Dict{"err": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `%s`, template.HTMLEscapeString(msg))
			return
		}
		user, found, err = app.db.GetUserByNameOrEmail(gu.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `%s`, template.HTMLEscapeString(err.Error()))
			return
		} else if !found || len(user) == 0 {
			msg, _ := app.i18n.T("user-not-found", Dict{"email": gu.Name})
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `%s`, template.HTMLEscapeString(msg))
			return
		}
	}
	if env.GetBool("OAUTH_DYN_LOGIN_TABLE_FILL", false) && env.GetBool("DYN_LOGIN_TABLE_MAP_TO_USERS", false) {
		dyn_db := env.GetString("DYN_LOGIN_DB", "")
		dyn_login_table := env.GetString("DYN_LOGIN_TABLE", "")
		dyn_login_table_map_to_users := env.GetBool("DYN_LOGIN_TABLE_MAP_TO_USERS", false)
		dyn_login_table_user_id := env.GetString("DYN_LOGIN_USER_ID_FIELD", "")
		dyn_login_table_email := env.GetString("DYN_LOGIN_EMAIL_FIELD", "")
		dyn_login_table_username := env.GetString("DYN_LOGIN_USERNAME_FIELD", "")
		dyn_login_table_active := env.GetString("DYN_LOGIN_ACTIVE_FIELD", "")
		// fmt.Println("DYN_LOGIN_TABLE_MAP_TO_USERS:", dyn_login_table_map_to_users, dyn_db, dyn_login_table, dyn_login_table_user_id, dyn_login_table_email, dyn_login_table_username, dyn_login_table_active)
		if dyn_db != "" && dyn_login_table != "" && dyn_login_table_map_to_users && dyn_login_table_user_id != "" && dyn_login_table_email != "" && dyn_login_table_username != "" && dyn_login_table_active != "" {
			sql := fmt.Sprintf(`select * from %s where %s = ?`, dyn_login_table, dyn_login_table_email)
			dyn_login_data, err := app.GetRowByFilter(sql, Dict{"db": dyn_db}, []any{gu.Email})
			if err != nil {
				fmt.Println("Error getting the user for", dyn_login_table, gu.Email, err.Error())
			} else if len(dyn_login_data) == 0 {
				insert := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, user_id, created_at, updated_at) VALUES (:name, :email, :active, :user_id, :created_at, :updated_at)`, dyn_login_table, dyn_login_table_username, dyn_login_table_email, dyn_login_table_active)
				name := gu.Name
				if name == "" {
					name = fmt.Sprintf("%s %s", gu.FirstName, gu.LastName)
				}
				data := Dict{
					"name":       gu.Name,
					"email":      gu.Email,
					"active":     true,
					"user_id":    user["id"],
					"created_at": time.Now().In(loc),
					"updated_at": time.Now().In(loc),
				}
				err = app.ExecuteQuery(insert, Dict{"db": dyn_db}, data)
				if err != nil {
					fmt.Println("Error inserting the user for", dyn_login_table, gu.Email, err.Error())
				}
			}
		}
	}
	user["oauth"] = true
	user["provider"] = provider
	delete(user, "password")
	delete(user, "created_at")
	delete(user, "updated_at")
	delete(user, "phone")
	delete(user, "timezone")
	delete(user, "attach_profile_pic")
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString(err.Error()))
		return
	}
	claims.Subject = string(json_user)
	//expiry := time.Now().Add(8 * time.Hour)
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	if !gu.ExpiresAt.IsZero() {
		// expiry = gu.ExpiresAt
	}
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)
	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.frontend_url}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString(err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    string(jwtBytes),
		HttpOnly: true,
		Secure:   true, // true
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	w.Header().Set("X-Auth-Provider", provider)      // ← optional
	w.Header().Set("X-Auth-Token", string(jwtBytes)) // ← optional
	/*if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}*/
	redirectRootHandler(w, r)
}

func redirectRootHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Determine scheme (http vs https)
	provider := r.PathValue("provider")
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else if forward := r.Header.Get("X-Forwarded-Proto"); forward != "" {
		scheme = forward
	}
	// 2. Build root-only URL (Scheme + Host + trailing slash)
	rootURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
	rootURL = strings.TrimSuffix(rootURL, "/oauth/"+provider+"/callback")
	// fmt.Println("Redirect URL:", rootURL)
	// 3. Issue a 303 See Other redirect (cleanest for refreshing post-actions)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", rootURL)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Redirect(w, r, rootURL, http.StatusSeeOther)
	}
}
