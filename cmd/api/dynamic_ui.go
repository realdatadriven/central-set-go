package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	texttemplate "text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/realdatadriven/central-set-go/internal/env"
)

func (app *application) getAnonymous() Dict {
	return Dict{"user_id": 2, "username": "anonymous", "role_id": 4, "active": true}
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b) // In production, handle errors appropriately
	return base64.URLEncoding.EncodeToString(b)
}

func (app *application) getUser(r *http.Request) (Dict, error) {
	if os.Getenv("COOKIE_MODE") == "TOKEN" {
		cookie, err := r.Cookie("session")
		if err != nil {
			return nil, err
		}
		return app.verifyTokenString(cookie.Value)
	} else {
		// 1. Retrieve the session cookie from the request
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// http.Error(w, "Unauthorized: No session cookie found", http.StatusUnauthorized)
			return nil, fmt.Errorf("Unauthorized: No session cookie found")
		}
		// 2. Look up the session ID in the server-side store
		user, exists := app.SessionStore.data.Load(cookie.Value)
		if !exists {
			// http.Error(w, "Unauthorized: Invalid or expired session", http.StatusUnauthorized)
			return nil, fmt.Errorf("Unauthorized: Invalid or expired session")
		}
		if _, ok := user.(Dict); !ok {
			return nil, fmt.Errorf("Unauthorized: Unable to get the user data from session")
		}
		// fmt.Println("session_id", cookie.Value, user.(Dict)["username"], user.(Dict)["user_id"], user.(Dict)["role_id"])
		return user.(Dict), nil
	}
}

func (app *application) serve_ui_page(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pathParams := Dict{}
	for k := range q {
		pathParams[k] = q.Get(k)
	}
	//r.URL.Path
	// anonymous
	user := app.getAnonymous()
	userFromSess, err := app.getUser(r)
	if err == nil {
		// fmt.Println("USER:", userFromSess)
		if r.PathValue("page_key") == "login" {
			// w.Header().Set("HX-Redirect", "/ui/"+r.PathValue("ui_slug"))
			// w.WriteHeader(http.StatusOK)
			http.Redirect(w, r, "/ui/"+r.PathValue("ui_slug"), http.StatusSeeOther)
			return
		}
		user = userFromSess
	} else {
		// fmt.Println("ERR:", err)
	}
	params := Dict{
		"lang": "en",
		"user": user,
		"data": Dict{
			"db":          env.GetString("UIDB", "UI"),
			"ui_slug":     r.PathValue("ui_slug"),
			"page_key":    r.PathValue("page_key"),
			"path_params": pathParams,
		},
	}
	params["host"] = getHost(r)
	params["path"] = r.URL.Path
	params["ip"] = ClientIP(r)
	params["loc"] = app.getLocationFromRequest(r, Dict{})
	res := app.RenderUIPage(params)
	if success, _ := res["success"].(bool); !success {
		http.Error(w, fmt.Sprintf("%v", res["msg"]), http.StatusNotFound)
		return
	}
	if secs := app.toInt(res["cache_seconds"]); secs > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", secs))
	}
	w.Header().Set("Content-Type", res["content_type"].(string))
	w.Write([]byte(res["data"].(string)))
}

func (app *application) serve_ui_partial(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pathParams := Dict{}
	for k := range q {
		pathParams[k] = q.Get(k)
	}
	user := app.getAnonymous()
	userFromSess, err := app.getUser(r)
	if err == nil {
		user = userFromSess
	}
	params := Dict{
		"lang": "en",
		"user": user,
		"data": Dict{
			"db":           env.GetString("UIDB", "UI"),
			"ui_slug":      r.PathValue("ui_slug"),
			"partial_name": r.PathValue("partial_name"),
			"path_params":  pathParams,
			"raw_query":    r.URL.RawQuery,
		},
	}
	params["host"] = getHost(r)
	params["path"] = r.URL.Path
	params["ip"] = ClientIP(r)
	params["loc"] = app.getLocationFromRequest(r, Dict{})
	// same token/user pattern as read_odata, if the partial is auth-gated
	res := app.RenderUIPartial(params)
	if success, _ := res["success"].(bool); !success {
		http.Error(w, fmt.Sprintf("%v", res["msg"]), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", res["content_type"].(string))
	w.Write([]byte(res["data"].(string)))
}

func (app *application) RenderUIPage(params Dict) Dict {
	data, ok := params["data"].(Dict)
	if !ok {
		return Dict{"success": false, "msg": `missing "data" params`}
	}
	uiSlug, _ := data["ui_slug"].(string)
	if uiSlug == "" {
		uiSlug, _ = data["ui_name"].(string)
	}
	if uiSlug == "" {
		return Dict{"success": false, "msg": `"ui_slug" or "ui_name" is required`}
	}
	pageKey, _ := data["page_key"].(string)
	pathParams, _ := data["path_params"].(Dict)
	if pathParams == nil {
		pathParams = Dict{}
	}
	// 1) resolve the website (`ui` row)
	sql := `select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`
	ui, err := app.GetRowByFilter(sql, params, []any{uiSlug, uiSlug})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui: %s", err)}
	}
	if len(ui) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("ui %q not found", uiSlug)}
	}
	uiID := ui["ui_id"]
	// 2) resolve the page: explicit page_key, or the default_page = true one
	var page Dict
	if pageKey != "" {
		sql = `select * from "ui_page" where ui_id = ? and page_key = ? and active = true and excluded = false`
		page, err = app.GetRowByFilter(sql, params, []any{uiID, pageKey})
	} else {
		sql = `select * from "ui_page" where ui_id = ? and default_page = true and active = true and excluded = false`
		page, err = app.GetRowByFilter(sql, params, []any{uiID, pageKey})
	}
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_page: %s", err)}
	}
	if len(page) == 0 {
		msg := fmt.Sprintf("no default page configured for ui %q", uiSlug)
		if pageKey != "" {
			msg = fmt.Sprintf("page %q not found for ui %q", pageKey, uiSlug)
		}
		return Dict{"success": false, "msg": msg}
	}
	// login_required
	if app.toBool(page["login_required"]) && app.toInt(params["user"].(Dict)["role_id"]) == 4 {
		msg := fmt.Sprintf("page %q requires login", pageKey)
		return Dict{"success": false, "msg": msg, "login_required": true}
	}
	pageID := page["ui_page_id"]
	// 3) active partials for this ui, in a stable order
	sql = `select up.*
	from "ui_page_partial" upp
	join "ui_partial" up on upp.ui_partial_id = up.ui_partial_id
	where upp.ui_page_id = ? 
		and upp.ui_id = ? 
		and upp.active = true 
		and upp.excluded = false 
	order by up.ui_partial_id asc`
	partials, err := app.GetRowsByFilter(sql, params, []any{pageID, uiID})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_page_partial: %s", err)}
	}
	// 4) build template data: resolve ui_page_data first, then every active
	// partial's ui_partial_data, each keyed by its own data name.
	tmplData := Dict{
		"Host":       params["host"],
		"Path":       params["path"],
		"UI":         ui,
		"Page":       page,
		"PathParams": pathParams,
		"config":     Dict{"frontend_url": app.config.frontend_url},
		"user":       params["user"],
	}
	sql = `select * from "ui_page_data" where ui_page_id = ? and active = true and excluded = false`
	pageDataRows, err := app.GetRowsByFilter(sql, params, []any{pageID})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_page_data: %s", err)}
	}
	if msg := app.resolveUIData(params, pageDataRows, "ui_page_data", pathParams, tmplData); msg != "" {
		return Dict{"success": false, "msg": msg}
	}
	for _, p := range partials {
		sql = `select * from "ui_partial_data" where ui_partial_id = ? and active = true and excluded = false`
		partialDataRows, perr := app.GetRowsByFilter(sql, params, []any{p["ui_partial_id"]})
		if perr != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_partial_data: %s", perr)}
		}
		if msg := app.resolveUIData(params, partialDataRows, "ui_partial_data", pathParams, tmplData); msg != "" {
			return Dict{"success": false, "msg": msg}
		}
	}
	// 5) parse the page template plus every active partial as named templates
	// so the page template can invoke them, e.g. {{template "header" .}}.
	pageTemplate, _ := page["page_template"].(string)
	tset, err := template.New("__page__").Funcs(sprig.FuncMap()).Parse(pageTemplate)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to parse page_template: %s", err)}
	}
	for _, p := range partials {
		name, _ := p["ui_partial"].(string)
		body, _ := p["partial_template"].(string)
		if name == "" {
			continue
		}
		if _, err := tset.New(name).Funcs(sprig.FuncMap()).Parse(body); err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to parse partial %q: %s", name, err)}
		}
	}
	var out bytes.Buffer
	if err := tset.ExecuteTemplate(&out, "__page__", tmplData); err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to render page: %s", err)}
	}
	return Dict{
		"success":       true,
		"data":          out.String(),
		"content_type":  "text/html; charset=utf-8",
		"cache_seconds": page["cache_seconds"],
		"ui":            ui,
		"page":          page,
	}
}

func (app *application) RenderUIPartial(params Dict) Dict {
	data, ok := params["data"].(Dict)
	if !ok {
		return Dict{"success": false, "msg": `missing "data" params`}
	}
	uiSlug, _ := data["ui_slug"].(string)
	if uiSlug == "" {
		uiSlug, _ = data["ui_name"].(string)
	}
	if uiSlug == "" {
		return Dict{"success": false, "msg": `"ui_slug" or "ui_name" is required`}
	}
	partialName, _ := data["partial_name"].(string)
	if partialName == "" {
		partialName, _ = data["ui_partial"].(string)
	}
	if partialName == "" {
		return Dict{"success": false, "msg": `"partial_name" is required`}
	}
	pathParams, _ := data["path_params"].(Dict)
	if pathParams == nil {
		pathParams = Dict{}
	}
	// 1) resolve the website (`ui` row)
	sql := `select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`
	ui, err := app.GetRowByFilter(sql, params, []any{uiSlug, uiSlug})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui: %s", err)}
	}
	if len(ui) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("ui %q not found", uiSlug)}
	}
	uiID := ui["ui_id"]
	// 2) every active partial for this ui, so cross-partial {{template}} calls
	// keep working; the requested one must be among them.
	sql = `select * from "ui_partial" where ui_id = ? and active = true and excluded = false order by ui_partial_id asc`
	partials, err := app.GetRowsByFilter(sql, params, []any{uiID})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_partial: %s", err)}
	}
	var target Dict
	for _, p := range partials {
		if name, _ := p["ui_partial"].(string); name == partialName {
			target = p
			break
		}
	}
	if target == nil {
		return Dict{"success": false, "msg": fmt.Sprintf("partial %q not found for ui %q", partialName, uiSlug)}
	}
	// 3) resolve ui_partial_data for every active partial (not just the
	// requested one), so nested partial calls have their data available too.
	tmplData := Dict{
		"Host":       params["host"],
		"Path":       params["path"],
		"UI":         ui,
		"PathParams": pathParams,
		"query":      data["raw_query"],
		"config":     Dict{"frontend_url": app.config.frontend_url},
		"user":       params["user"],
	}
	for _, p := range partials {
		sql = `select * from "ui_partial_data" where ui_partial_id = ? and active = true and excluded = false`
		partialDataRows, perr := app.GetRowsByFilter(sql, params, []any{p["ui_partial_id"]})
		if perr != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_partial_data: %s", perr)}
		}
		if msg := app.resolveUIData(params, partialDataRows, "ui_partial_data", pathParams, tmplData); msg != "" {
			return Dict{"success": false, "msg": msg}
		}
	}
	// 4) parse every active partial as a named template, then execute only
	// the requested one.
	tset := template.New(partialName).Funcs(sprig.FuncMap())
	for _, p := range partials {
		name, _ := p["ui_partial"].(string)
		body, _ := p["partial_template"].(string)
		if name == "" {
			continue
		}
		var perr error
		if name == partialName {
			tset, perr = tset.Parse(body)
		} else {
			_, perr = tset.New(name).Parse(body)
		}
		if perr != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to parse partial %q: %s", name, perr)}
		}
	}
	var out bytes.Buffer
	if err := tset.ExecuteTemplate(&out, partialName, tmplData); err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to render partial %q: %s", partialName, err)}
	}
	return Dict{
		"success":      true,
		"data":         out.String(),
		"content_type": "text/html; charset=utf-8",
		"ui":           ui,
		"partial":      target,
	}
}

func (app *application) resolveUIData(base Dict, rows []Dict, nameCol string, pathParams Dict, tmplData Dict) string {
	for _, row := range rows {
		name, _ := row[nameCol].(string)
		odataPath, _ := row["odata_path"].(string)
		if name == "" || odataPath == "" {
			continue
		}
		fmt.Println("ODATA PATH:", odataPath, pathParams)
		resolvedPath, err := app.renderODataPath(odataPath, pathParams)
		if err != nil {
			return fmt.Sprintf("failed to resolve odata_path for %q: %s", name, err)
		}
		fmt.Println("ODATA PATH:", resolvedPath)
		// fresh params per call: app.ODataRead writes db/table/filters into
		// params["data"], so each data source needs its own Dict.
		callParams := Dict{"data": Dict{}}
		if l, ok := base["lang"]; ok {
			callParams["lang"] = l
		}
		if u, ok := base["user"]; ok {
			callParams["user"] = u
		}
		if l, ok := base["location"]; ok {
			callParams["location"] = l
		}
		res := app.ODataRead(callParams, resolvedPath)
		if success, ok := res["success"].(bool); !ok || !success {
			return fmt.Sprintf("failed to load %q: %v", name, res["msg"])
		}
		single := app.toBool(row["single_row_obj"])
		switch d := res["data"].(type) {
		case Dict:
			if single {
				tmplData[name] = res
			} else {
				tmplData[name] = res
				tmplData[name].(Dict)["data"] = []Dict{d}
			}
		case []Dict:
			if single {
				if len(d) > 0 {
					tmplData[name] = res
					tmplData[name].(Dict)["data"] = d[0]
				} else {
					tmplData[name] = res
					tmplData[name].(Dict)["data"] = Dict{}
				}
			} else {
				tmplData[name] = res
			}
		default:
			tmplData[name] = res
		}
		// fmt.Println("DATA:", tmplData[name].(Dict)["msg"], tmplData[name].(Dict)["data"])
		// fmt.Println("ODataRead result for", name, ":", slices.Collect(maps.Keys(tmplData[name].(Dict))))
	}
	return ""
}

// STORE/product?$filter=product_id eq {{.PathParams.product_id}}
func (app *application) renderODataPath(odataPath string, pathParams Dict) (string, error) {
	t, err := texttemplate.New("odata_path").Parse(odataPath)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := t.Execute(&out, Dict{"PathParams": pathParams}); err != nil {
		return "", err
	}
	_path := strings.ReplaceAll(out.String(), "\r\n", "")
	_path = strings.ReplaceAll(_path, "\n", "")
	_path = strings.ReplaceAll(_path, "\r", "")
	re := regexp.MustCompile(`\s+`)
	_path = strings.TrimSpace(re.ReplaceAllString(_path, " "))
	re = regexp.MustCompile(`\s+&|&\s+`)
	_path = strings.TrimSpace(re.ReplaceAllString(_path, "&"))
	re = regexp.MustCompile(`\?&`)
	_path = strings.TrimSpace(re.ReplaceAllString(_path, "?"))
	return _path, nil
}

func (app *application) RenderUIAsset(params Dict) Dict {
	data, ok := params["data"].(Dict)
	if !ok {
		return Dict{"success": false, "msg": `missing "data" params`}
	}
	uiSlug, _ := data["ui_slug"].(string)
	if uiSlug == "" {
		uiSlug, _ = data["ui_name"].(string)
	}
	if uiSlug == "" {
		return Dict{"success": false, "msg": `"ui_slug" or "ui_name" is required`}
	}
	assetPath, _ := data["asset_path"].(string)
	if assetPath == "" {
		return Dict{"success": false, "msg": `"asset_path" is required`}
	}
	sql := `select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`
	ui, err := app.GetRowByFilter(sql, params, []any{uiSlug, uiSlug})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui: %s", err)}
	}
	if len(ui) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("ui %q not found", uiSlug)}
	}
	uiID := ui["ui_id"]
	sql = `select * from "ui_asset" where ui_id = ? and asset_path = ? and active = true and excluded = false`
	asset, err := app.GetRowByFilter(sql, params, []any{uiID, assetPath})
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_asset: %s", err)}
	}
	if len(asset) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("asset %q not found for ui %q", assetPath, uiSlug)}
	}
	content, _ := asset["asset_content"].(string)
	encoding, _ := asset["content_encoding"].(string)
	var raw []byte
	if strings.EqualFold(encoding, "base64") {
		decoded, derr := base64.StdEncoding.DecodeString(content)
		if derr != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to decode asset %q: %s", assetPath, derr)}
		}
		raw = decoded
	} else {
		raw = []byte(content)
	}
	mimeType, _ := asset["mime_type"].(string)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	cacheSeconds := app.toInt(asset["cache_seconds"])
	checksum, _ := asset["checksum"].(string)
	lastModified, _ := app.toTime(asset["updated_at"])
	return Dict{
		"success":       true,
		"data":          raw,
		"mime_type":     mimeType,
		"cache_seconds": cacheSeconds,
		"checksum":      checksum,
		"last_modified": lastModified,
		"asset":         asset,
	}
}

// serve_ui_asset is the HTTP handler for GET /ui/{ui_slug}/static/{asset...}
func (app *application) serve_ui_asset(w http.ResponseWriter, r *http.Request) {
	params := Dict{
		"data": Dict{
			"db":         env.GetString("UIDB", "UI"),
			"ui_slug":    r.PathValue("ui_slug"),
			"asset_path": r.PathValue("asset"),
		},
	}
	params["ip"] = ClientIP(r)
	res := app.RenderUIAsset(params)
	if success, _ := res["success"].(bool); !success {
		http.Error(w, fmt.Sprintf("%v", res["msg"]), http.StatusNotFound)
		return
	}
	if mimeType, _ := res["mime_type"].(string); mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
	if secs := app.toInt(res["cache_seconds"]); secs > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", secs))
	}
	if checksum, _ := res["checksum"].(string); checksum != "" {
		w.Header().Set("ETag", fmt.Sprintf(`"%s"`, checksum))
	}
	// ServeContent reads any ETag already set above to answer If-None-Match,
	// compares lastModified against If-Modified-Since, and handles Range
	// requests - all for free.
	lastModified, _ := res["last_modified"].(time.Time)
	raw, _ := res["data"].([]byte)
	http.ServeContent(w, r, r.PathValue("asset"), lastModified, bytes.NewReader(raw))
}

func (app *application) toTime(v any) (time.Time, bool) {
	switch val := v.(type) {
	case time.Time:
		return val, true
	case string:
		layouts := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006-01-02",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, val); err == nil {
				return t, true
			}
		}
	case int64:
		return time.Unix(val, 0), true
	case float64:
		return time.Unix(int64(val), 0), true
	}
	return time.Time{}, false
}

func (app *application) ui_login(w http.ResponseWriter, r *http.Request) {
	uiSlug := r.PathValue("ui_slug")
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "login"
	}
	_, err := app.getUser(r)
	if err == nil {
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Redirect(w, r, "/ui/"+r.PathValue("ui_slug"), http.StatusSeeOther)
		}
		return
	}
	params := Dict{
		"lang": "en",
		"data": Dict{
			"db":      env.GetString("UIDB", "UI"),
			"ui_slug": r.PathValue("ui_slug"),
		},
	}
	params["ip"] = ClientIP(r)
	sql := `select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`
	ui, err := app.GetRowByFilter(sql, params, []any{uiSlug, uiSlug})
	if err != nil {
		app.writeHTMLError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch ui: %s", err))
		return
	}
	if len(ui) == 0 {
		app.writeHTMLError(w, http.StatusNotFound, fmt.Sprintf("ui %q not found", uiSlug))
		return
	}
	sql = `select p.* 
	from ui_page p
	where p.page_key = ? and p.active = true and p.excluded = false and p.ui_id = ?`
	pageData, err := app.GetRowByFilter(sql, params, []any{page, ui["ui_id"]})
	if err != nil {
		fmt.Println("Error geting page response template:", err)
	}
	response_tmpl, _ := pageData["response_tmpl"].(string)
	if err := r.ParseForm(); err != nil {
		if response_tmpl != "" {
			res, err := app.RenderTemplate(response_tmpl, Dict{"success": false, "msg": "Invalid form submission."})
			if err != nil {
				fmt.Println("Error rendering page response template:", err)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, res)
				return
			}
		}
		app.writeHTMLError(w, http.StatusBadRequest, "Invalid form submission.")
		return
	}
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if email == "" || password == "" {
		if response_tmpl != "" {
			res, err := app.RenderTemplate(response_tmpl, Dict{"success": false, "msg": "Email and password are required."})
			if err != nil {
				fmt.Println("Error rendering page response template:", err)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, res)
				return
			}
		}
		app.writeHTMLError(w, http.StatusBadRequest, "Email and password are required.")
		return
	}
	res := app._login(Dict{
		"data": Dict{
			"username": email,
			"password": password,
		},
	})
	success, _ := res["success"].(bool)
	if !success {
		msg, _ := res["msg"].(string)
		// fmt.Println(msg)
		if msg == "" {
			msg = "Invalid email or password."
		}
		if response_tmpl != "" {
			res, err := app.RenderTemplate(response_tmpl, Dict{"success": false, "msg": msg})
			if err != nil {
				fmt.Println("Error rendering page response template:", err)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, res)
				return
			}
		}
		app.writeHTMLError(w, http.StatusUnauthorized, msg)
		return
	}
	// TWO FACTOR AUTH
	if two_factor, ok := res["two_factor"]; ok && app.toBool(two_factor) {
		_link := fmt.Sprintf("/ui/%s/two-factor?username=%s", uiSlug, email)
		fmt.Println("2FactorRedirect:", _link)
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", _link)
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Redirect(w, r, _link, http.StatusSeeOther)
		}
	}
	app.startUISession(w, res)
	// w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
	// w.WriteHeader(http.StatusOK)
	_link := fmt.Sprintf("/ui/%s", uiSlug)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", _link)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Redirect(w, r, _link, http.StatusSeeOther)
	}
}

// startUISession persists the authenticated result returned by _login or
// two_factor_code_valid and issues the matching UI session cookie.
func (app *application) startUISession(w http.ResponseWriter, res Dict) {
	token, _ := res["token"].(string)
	if os.Getenv("COOKIE_MODE") == "TOKEN" {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		return
	}
	sessionID := generateSessionID()
	user, _ := res["data"].(Dict)
	user["token"] = token
	app.SessionStore.data.Store(sessionID, user)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// uiFormParams turns an HTML form submission into the params shape consumed by
// the existing authentication actions. r.Form deliberately includes query
// parameters so reset-password links can keep their token in the URL.
func (app *application) uiFormParams(w http.ResponseWriter, r *http.Request) (Dict, bool) {
	if err := r.ParseForm(); err != nil {
		app.writeHTMLError(w, http.StatusBadRequest, "Invalid form submission.")
		return nil, false
	}
	data := Dict{
		"db":      env.GetString("UIDB", "UI"),
		"ui_slug": r.PathValue("ui_slug"),
	}
	for key, values := range r.Form {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}
	lang, _ := data["lang"].(string)
	if lang == "" {
		lang = "en"
	}
	params := Dict{"lang": lang, "data": data}
	params["host"] = getHost(r)
	params["path"] = r.URL.Path
	params["ip"] = ClientIP(r)
	params["loc"] = app.getLocationFromRequest(r, data)
	return params, true
}

// ensureUIExists keeps auth form endpoints scoped to an active UI, just as
// ui_login does before it invokes the login action.
func (app *application) ensureUIExists(w http.ResponseWriter, uiSlug string, params Dict) bool {
	sql := `select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`
	ui, err := app.GetRowByFilter(sql, params, []any{uiSlug, uiSlug})
	if err != nil {
		app.writeHTMLError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch ui: %s", err))
		return false
	}
	if len(ui) == 0 {
		app.writeHTMLError(w, http.StatusNotFound, fmt.Sprintf("ui %q not found", uiSlug))
		return false
	}
	return true
}

func (app *application) writeHTMLSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<div class="alert alert-success"><span>%s</span></div>`, template.HTMLEscapeString(msg))
}

func (app *application) uiAction(w http.ResponseWriter, r *http.Request, action func(Dict) Dict) {
	params, ok := app.uiFormParams(w, r)
	if !ok {
		return
	}
	if !app.ensureUIExists(w, r.PathValue("ui_slug"), params) {
		return
	}
	res := action(params)
	if success, _ := res["success"].(bool); !success {
		msg, _ := res["msg"].(string)
		if msg == "" {
			msg = "Unable to complete the request."
		}
		app.writeHTMLError(w, http.StatusUnprocessableEntity, msg)
		return
	}
	msg, _ := res["msg"].(string)
	if msg == "" {
		msg = "Request completed successfully."
	}
	app.writeHTMLSuccess(w, msg)
}

// ui_validate_code completes a two-factor login and creates the regular UI
// session only after the supplied code has been validated.
func (app *application) ui_validate_code(w http.ResponseWriter, r *http.Request) {
	params, ok := app.uiFormParams(w, r)
	if !ok {
		return
	}
	if !app.ensureUIExists(w, r.PathValue("ui_slug"), params) {
		return
	}
	data := params["data"].(Dict)
	if data["username"] == "" && data["email"] != "" {
		data["username"] = data["email"]
	}
	ui := r.URL.Query().Get("ui")
	page := r.URL.Query().Get("page")
	sql := `select p.* 
	from ui_page p
	join ui on ui.ui_id = p.ui_id
	where p.page_key = ? and p.active = true and p.excluded = false
		and (ui.ui_slug = ? or ui.ui_name = ?) and ui.active = true and ui.excluded = false`
	pageData, err := app.GetRowByFilter(sql, params, []any{page, ui, ui})
	if err != nil {
		fmt.Println("Error geting page response template:", err)
	}
	response_tmpl, _ := pageData["response_tmpl"].(string)
	res := app.two_factor_code_valid(params)
	if success, _ := res["success"].(bool); !success {
		if response_tmpl != "" {
			res, err := app.RenderTemplate(response_tmpl, res)
			if err != nil {
				fmt.Println("Error rendering page response template:", err)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, res)
				return
			}
		}
		msg, _ := res["msg"].(string)
		if msg == "" {
			msg = "Invalid two-factor code."
		}
		app.writeHTMLError(w, http.StatusUnauthorized, msg)
		return
	}
	app.startUISession(w, res)
	uiSlug := r.PathValue("ui_slug")
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/ui/"+uiSlug, http.StatusSeeOther)
}

func (app *application) ui_signup(w http.ResponseWriter, r *http.Request) {
	app.uiAction(w, r, func(params Dict) Dict {
		params["login_table"] = os.Getenv("DYN_LOGIN_TABLE")
		params["user_id_field"] = os.Getenv("DYN_LOGIN_USER_ID_FIELD")
		params["dyn_login_role_id"] = os.Getenv("DYN_LOGIN_ROLE_ID")
		params["username_field"] = os.Getenv("DYN_LOGIN_USERNAME_FIELD")
		params["email_field"] = os.Getenv("DYN_LOGIN_EMAIL_FIELD")
		params["password_field"] = os.Getenv("DYN_LOGIN_PASSWORD_FIELD")
		params["active_field"] = os.Getenv("DYN_LOGIN_ACTIVE_FIELD")
		return app.dynamic_signup(params)
	})
}

func (app *application) ui_recover_pass(w http.ResponseWriter, r *http.Request) {
	// app.uiAction(w, r, app.recover_pass)
	params, ok := app.uiFormParams(w, r)
	if !ok {
		return
	}
	if !app.ensureUIExists(w, r.PathValue("ui_slug"), params) {
		return
	}
	ui := r.URL.Query().Get("ui")
	page := r.URL.Query().Get("page")
	sql := `select p.* 
	from ui_page p
	join ui on ui.ui_id = p.ui_id
	where p.page_key = ? and p.active = true and p.excluded = false
		and (ui.ui_slug = ? or ui.ui_name = ?) and ui.active = true and ui.excluded = false`
	pageData, err := app.GetRowByFilter(sql, params, []any{page, ui, ui})
	if err != nil {
		fmt.Println("Error geting page response template:", err)
	}
	response_tmpl, _ := pageData["response_tmpl"].(string)
	res := app.recover_pass(params)
	if response_tmpl != "" {
		res, err := app.RenderTemplate(response_tmpl, res)
		if err != nil {
			fmt.Println("Error rendering page response template:", err)
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, res)
			return
		}
	}
	msg, _ := res["msg"].(string)
	if msg == "" {
		msg = "Unexpected Error"
	}
	app.writeHTMLError(w, http.StatusUnauthorized, msg)
	return
}

func (app *application) ui_reset_pass(w http.ResponseWriter, r *http.Request) {
	// app.uiAction(w, r, app.reset_pass)
	params, ok := app.uiFormParams(w, r)
	if !ok {
		return
	}
	if !app.ensureUIExists(w, r.PathValue("ui_slug"), params) {
		return
	}
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = getHost(r)
	}
	ui := r.URL.Query().Get("ui")
	page := r.URL.Query().Get("page")
	sql := `select p.* 
	from ui_page p
	join ui on ui.ui_id = p.ui_id
	where p.page_key = ? and p.active = true and p.excluded = false
		and (ui.ui_slug = ? or ui.ui_name = ?) and ui.active = true and ui.excluded = false`
	pageData, err := app.GetRowByFilter(sql, params, []any{page, ui, ui})
	if err != nil {
		fmt.Println("Error geting page response template:", err)
	}
	response_tmpl, _ := pageData["response_tmpl"].(string)
	res := app.reset_pass(params)
	if success, _ := res["success"].(bool); !success {
		if response_tmpl != "" {
			res, err := app.RenderTemplate(response_tmpl, res)
			if err != nil {
				fmt.Println("Error rendering page response template:", err)
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, res)
				return
			}
		}
		msg, _ := res["msg"].(string)
		if msg == "" {
			msg = "Unexpected Error"
		}
		app.writeHTMLError(w, http.StatusUnauthorized, msg)
		return
	}
	// fmt.Println(res)
	app.logoutHandler(w, r)
	/*if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirect)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)*/
}

func (app *application) ui_alter_pass(w http.ResponseWriter, r *http.Request) {
	user, err := app.getUser(r)
	if err != nil {
		app.writeHTMLError(w, http.StatusUnauthorized, "Authentication is required.")
		return
	}
	app.uiAction(w, r, func(params Dict) Dict {
		data := params["data"].(Dict)
		if username, _ := user["username"].(string); username != "" {
			data["username"] = username
		}
		return app.alter_pass(params)
	})
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	uiSlug := r.PathValue("ui_slug")
	_, err := app.getUser(r)
	if err != nil {
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/ui/"+uiSlug+"/login")
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Redirect(w, r, "/ui/"+uiSlug+"/login", http.StatusSeeOther)
		}
		return
	}
	if os.Getenv("COOKIE_MODE") == "TOKEN" {
		_, err := r.Cookie("session")
		if err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				MaxAge:   -1, // Forces immediate deletion
				HttpOnly: true,
			})
		}
	} else {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			// 1. Delete session from the server store
			app.SessionStore.data.Delete(cookie.Value)
		}
		// 2. Expire the cookie on the client browser immediately
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Forces immediate deletion
			HttpOnly: true,
		})
	}
	if r.Method == http.MethodPost {
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Redirect(w, r, "/ui/"+uiSlug, http.StatusSeeOther)
		}
		return
	}
	fmt.Fprintln(w, "Logged out successfully!")
}

// writeHTMLError writes a small HTML fragment (a DaisyUI alert) with the
func (app *application) writeHTMLError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `<div class="alert alert-error"><span>%s</span></div>`, template.HTMLEscapeString(msg))
}

func (app *application) handleConfirmEmail(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("token")
	if code == "" {
		code = r.URL.Query().Get("code")
	}
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = getHost(r)
	}
	redirect = fmt.Sprintf("%s?token=%s", redirect, code)
	params := Dict{"lang": "en", "data": Dict{"token": code}}
	params["host"] = getHost(r)
	params["path"] = r.URL.Path
	params["ip"] = ClientIP(r)
	params["loc"] = app.getLocationFromRequest(r, Dict{})
	// params["token"] = code
	res := app.handle_confirm_email(params)
	if success, _ := res["success"].(bool); !success {
		msg, _ := res["msg"].(string)
		if msg == "" {
			msg = "Unexpected Error"
		}
		app.writeHTMLError(w, http.StatusUnauthorized, msg)
		return
	}
	// fmt.Println(res)
	app.startUISession(w, res)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirect)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
