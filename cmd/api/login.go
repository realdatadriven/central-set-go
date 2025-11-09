package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/realdatadriven/central-set-go/internal/password"
	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"
	"github.com/realdatadriven/central-set-go/internal/validator"

	"github.com/pascaldekloe/jwt"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Lang string `json:"lang"`
		Data struct {
			Username  string              `json:"username"`
			Password  string              `json:"password"`
			Validator validator.Validator `json:"-"`
		}
	}

	err := request.DecodeJSON(w, r, &params)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if params.Data.Username == "" {
		err = response.JSON(w, http.StatusOK, map[string]any{
			"success": false,
			"msg":     "Email is required",
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	if params.Data.Password == "" {
		err = response.JSON(w, http.StatusOK, map[string]any{
			"success": false,
			"msg":     "Password is required",
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	user, found, err := app.db.GetUserByNameOrEmail(params.Data.Username)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if !found || len(user) == 0 {
		err = response.JSON(w, http.StatusOK, map[string]any{
			"success": false,
			"msg":     "Email or Password incorrect",
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	if found {
		passwordMatches, err := password.Matches(params.Data.Password, user["password"].(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		if !passwordMatches {
			err = response.JSON(w, http.StatusOK, map[string]any{
				"success": false,
				"msg":     "Email or Password incorrect",
			})
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}
	}

	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		fmt.Print(err)
	}
	claims.Subject = string(json_user) //strconv.Itoa(user["user_id"].(strings))
	//fmt.Print(claims.Subject)
	//claims.Subject = user["user_id"].(string) // strconv.Itoa(int64(user["user_id"])
	expiry := time.Now().Add(time.Duration(app.config.jwt.tokenExpireHours) * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"success": true,
		"msg":     "Loged in successfully!",
		"data":    user,
		"token":   string(jwtBytes),
		"expiry":  expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// AuthenticateAD connects to AD and validates a user's credentials.
func AuthenticateAD(ldapURL, baseDN, serviceUser, servicePass, username, password string, skipVerify bool) (bool, map[string]any, error) {
	//fmt.Println("AuthenticateAD:", ldapURL, baseDN, serviceUser, servicePass, username, password)
	// Connect to LDAP server (use LDAPS if possible)
	l, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: skipVerify}))
	if err != nil {
		return false, nil, fmt.Errorf("failed to connect LDAP: %w", err)
	}
	defer l.Close()

	// 1️⃣ Bind as service account (to search)
	err = l.Bind(serviceUser, servicePass)
	if err != nil {
		return false, nil, fmt.Errorf("service bind failed: %w", err)
	}

	// 2️⃣ Search for user by sAMAccountName, userPrincipalName, or mail
	filter := fmt.Sprintf("(|(sAMAccountName=%[1]s)(userPrincipalName=%[1]s)(mail=%[1]s))", username)
	fmt.Println(baseDN, filter)
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

	fmt.Println("LDAP search entries:", len(sr.Entries))

	if len(sr.Entries) != 1 {
		return false, nil, fmt.Errorf("user not found or multiple matches for %s", username)
	}

	entry := sr.Entries[0]
	userDN := entry.DN

	// 3️⃣ Verify password by binding as user
	if err := l.Bind(userDN, password); err != nil {
		return false, nil, nil // invalid credentials
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
	userInfo := make(map[string]any)
	// 4️⃣ Extract user attributes into a map
	userInfo["dn"] = userDN
	userInfo["cn"] = entry.GetAttributeValue("cn")
	userInfo["displayName"] = entry.GetAttributeValue("displayName")
	userInfo["givenName"] = entry.GetAttributeValue("givenName")
	userInfo["sn"] = entry.GetAttributeValue("sn")
	userInfo["first_name"] = firstName
	userInfo["last_name"] = lastName
	userInfo["username"] = entry.GetAttributeValue("sAMAccountName")
	userInfo["email"] = entry.GetAttributeValue("mail")
	userInfo["upn"] = entry.GetAttributeValue("userPrincipalName")
	userInfo["groups"] = entry.GetAttributeValues("memberOf")

	return true, userInfo, nil
}

func (app *application) _login(params map[string]any) map[string]any {
	_data := map[string]any{}
	if _, ok := params["data"]; ok {
		_data = params["data"].(map[string]any)
	}
	if app.IsEmpty(_data) {
		msg, _ := app.i18n.T("no-data", map[string]any{})
		return map[string]any{
			"success": true,
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
	var user map[string]any
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
			return map[string]any{
				"success": false,
				"msg":     err.Error(),
			}
		}
		user, found, err = app.db.GetUserByNameOrEmail(username)
		if err != nil || len(user) == 0 {
			user = map[string]any{
				"username":   userInfo["username"],
				"password":   "LDAP_USER", //userInfo["password"],
				"first_name": userInfo["first_name"],
				"last_name":  userInfo["last_name"],
				"email":      userInfo["email"],
				"lang_id":    1,
				"role_id":    2,
				"active":     true,
				"excluded":   false,
				"create_at":  time.Now(),
				"updated_at": time.Now(),
			}
			err := app.AdminInsertData("users", user)
			if err != nil {
				return map[string]any{
					"success": false,
					"msg":     err.Error(),
				}
			} else {
				user, found, err = app.db.GetUserByNameOrEmail(username)
				if err != nil {
					return map[string]any{
						"success": false,
						"msg":     err.Error(),
					}
				}
			}
		}
	} else {
		user, found, err = app.db.GetUserByNameOrEmail(username)
		if err != nil || len(user) == 0 {
			return map[string]any{
				"success": false,
				"msg":     err.Error(),
			}
		}
		if found {
			//_hash, _ := password.Hash(pass)
			//fmt.Println(pass, _hash, user["password"].(string))
			match, err := password.Matches(pass, user["password"].(string))
			if err != nil {
				return map[string]any{
					"success": false,
					"msg":     err.Error(),
				}
			}
			if !match {
				msg, _ := app.i18n.T("user-pass-incorrect", map[string]any{})
				return map[string]any{
					"success": false,
					"msg":     msg,
				}
			}
		}
	}
	var claims jwt.Claims
	json_user, err := json.Marshal(user)
	if err != nil {
		return map[string]any{
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
		return map[string]any{
			"success": false,
			"msg":     err.Error(),
		}
	}
	data := map[string]any{
		"success": true,
		"msg":     "Loged in successfully!",
		"data":    user,
		"token":   string(jwtBytes),
		"expiry":  expiry.Format(time.RFC3339),
	}
	return data
}

func (app *application) alter_pass(params map[string]any) map[string]any {
	_check_login := app._login(params)
	if success, ok := _check_login["success"].(bool); ok && success {
		_data := map[string]any{}
		if _, ok := params["data"]; ok {
			_data = params["data"].(map[string]any)
		}
		newPassword, _ := _data["new_password"].(string)
		oldPassword, _ := _data["password"].(string)
		if newPassword == "" {
			msg, _ := app.i18n.T("new_pass_is_required", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		if newPassword == oldPassword {
			msg, _ := app.i18n.T("new_pass_old_pass", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		if len(newPassword) < 8 {
			msg, _ := app.i18n.T("password_min_length", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		hasUpper, _ := regexp.MatchString(`[A-Z]`, newPassword)
		if !hasUpper {
			msg, _ := app.i18n.T("pass_must_have_upper", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		hasNumber, _ := regexp.MatchString(`[0-9]`, newPassword)
		if !hasNumber {
			msg, _ := app.i18n.T("pass_must_have_number", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		hasSpecial, _ := regexp.MatchString(`[$&+,:;=?@#!*ªº.-]`, newPassword)
		if !hasSpecial {
			msg, _ := app.i18n.T("pass_must_have_special", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		query := `UPDATE users 
			SET password = :password 
		WHERE email = :username
			OR username = :username`
		pass, err := password.Hash(newPassword)
		if err != nil {
			msg, _ := app.i18n.T("pass_must_have_special", map[string]any{})
			return map[string]any{"success": false, "msg": msg}
		}
		username := ""
		if _, ok := _data["username"].(string); ok {
			username = _data["username"].(string)
		} else if _, ok := _data["user"].(string); ok {
			username = _data["user"].(string)
		} else if _, ok := _data["u"].(string); ok {
			username = _data["u"].(string)
		}
		_data = map[string]any{"username": username, "password": pass}
		_, err = app.db.ExecuteNamedQuery(query, _data)
		if err != nil {
			msg, _ := app.i18n.T("unexpected-error", map[string]any{"err": err.Error()})
			return map[string]any{
				"success": false,
				"msg":     msg,
			}
		}
		msg, _ := app.i18n.T("alter-pass-success", map[string]any{})
		return map[string]any{
			"success": true,
			"msg":     msg,
		}
	} else {
		return _check_login
	}
}
