package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"slices"
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
	/*cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}
	return app.verifyTokenString(cookie.Value)*/
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
	return user.(Dict), nil
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
		fmt.Println("ERR:", err)
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
	// same token/user pattern as read_odata, if the partial is auth-gated
	res := app.RenderUIPartial(params)
	if success, _ := res["success"].(bool); !success {
		http.Error(w, fmt.Sprintf("%v", res["msg"]), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", res["content_type"].(string))
	w.Write([]byte(res["data"].(string)))
}

// RenderUIPage resolves a `ui` website + one of its `ui_page` rows, fetches
// every data source needed by the page and by its active `ui_partial` rows
// (via app.ODataRead), and renders the final HTML by parsing the page
// template together with all active partial templates.
//
// params["data"] is expected to contain:
//
//	ui_slug / ui_name : string  - required, identifies the `ui` row.
//	page_key          : string  - optional. When empty, the `ui_page` row
//	                    with default_page = true is served instead.
//	path_params       : Dict    - optional. Route/path values made available
//	                    to odata_path templates as {{.PathParams.xxx}},
//	                    e.g. odata_path:
//	                    "STORE/product?$filter=product_id eq {{.PathParams.product_id}}"
//
// Any other keys already on params (user, lang, location, ...) are copied,
// unchanged, into every app.ODataRead call so permissions/RLA behave exactly
// like they do on the regular /odata and /read_odata endpoints.
//
// On success, returns:
//
//	Dict{
//	  "success":      true,
//	  "data":         "<rendered html>",
//	  "content_type": "text/html; charset=utf-8",
//	  "cache_seconds": <ui_page.cache_seconds>,
//	  "ui":           <ui row>,
//	  "page":         <ui_page row>,
//	}
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
		"UI":         ui,
		"Page":       page,
		"PathParams": pathParams,
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

// RenderUIPartial resolves a single `ui_partial` by name for a given `ui`
// and renders it standalone (no page_template wrapper) — intended for htmx
// requests that fetch or refresh just one fragment of a page, e.g.
// GET /ui/{ui_slug}/partial/{partial_name}.
//
// params["data"] is expected to contain:
//
//	ui_slug / ui_name : string  - required, identifies the `ui` row.
//	partial_name / ui_partial : string - required, identifies the `ui_partial`
//	                    row (the ui_partial column, e.g. "header").
//	path_params       : Dict    - optional, same as RenderUIPage: made
//	                    available to odata_path templates as
//	                    {{.PathParams.xxx}}.
//
// Every other active partial for the same `ui` is parsed into the same
// template set (and its own ui_partial_data resolved) so that, if the
// requested partial invokes another one via {{template "name" .}}, that
// nested lookup still resolves. Only the requested partial is executed as
// the response body.
//
// On success, returns:
//
//	Dict{
//	  "success":      true,
//	  "data":         "<rendered html fragment>",
//	  "content_type": "text/html; charset=utf-8",
//	  "ui":           <ui row>,
//	  "partial":      <ui_partial row>,
//	}
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
		"UI":         ui,
		"PathParams": pathParams,
		"query":      data["raw_query"],
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

// resolveUIData fetches, for every row in `rows`, the data source named by
// nameCol ("ui_page_data" or "ui_partial_data") through app.ODataRead, after
// resolving any {{.PathParams.xxx}} placeholders in its odata_path. The
// result is stored in tmplData under the row's own data-name, as a single
// Dict when single_row_obj is true, or as a []Dict otherwise. Returns a
// non-empty error message on failure, "" on success.
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
		fmt.Println("DATA:", tmplData[name].(Dict)["msg"], tmplData[name].(Dict)["data"])
		fmt.Println("ODataRead result for", name, ":", slices.Collect(maps.Keys(tmplData[name].(Dict))))
	}
	return ""
}

// renderODataPath executes odata_path as a text/template, giving stored
// odata_path values access to path/route params, e.g.
//
//	STORE/product?$filter=product_id eq {{.PathParams.product_id}}
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
	return _path, nil
}

// RenderUIAsset resolves a single `ui_asset` row for a given `ui` by its
// relative asset_path (e.g. "assets/store.css"), decodes its content
// according to content_encoding, and returns it ready to be written to an
// http.ResponseWriter (or fed to http.ServeContent) — including a parsed
// last_modified time (from updated_at) for browser caching.
//
// params["data"] is expected to contain:
//
//	ui_slug / ui_name : string - required, identifies the `ui` row.
//	asset_path        : string - required, matches ui_asset.asset_path.
//
// On success, returns:
//
//	Dict{
//	  "success":       true,
//	  "data":          []byte,          // decoded asset content
//	  "mime_type":     string,
//	  "cache_seconds": int,
//	  "checksum":      string,          // "" if not stored
//	  "last_modified": time.Time,       // zero value if updated_at is absent/unparsable
//	  "asset":         Dict,            // the raw ui_asset row
//	}
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
// (register the route with the trailing "..." wildcard so asset paths like
// "assets/img/logo.png" match as a single value). It defers conditional-GET
// (If-None-Match / If-Modified-Since) and Range handling to
// http.ServeContent: the ETag comes from the asset's stored checksum when
// present, and Last-Modified from updated_at.
func (app *application) serve_ui_asset(w http.ResponseWriter, r *http.Request) {
	params := Dict{
		"data": Dict{
			"db":         env.GetString("UIDB", "UI"),
			"ui_slug":    r.PathValue("ui_slug"),
			"asset_path": r.PathValue("asset"),
		},
	}
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
	// requests — all for free.
	lastModified, _ := res["last_modified"].(time.Time)
	raw, _ := res["data"].([]byte)
	http.ServeContent(w, r, r.PathValue("asset"), lastModified, bytes.NewReader(raw))
}

// toTime best-effort converts common datetime representations returned by
// different DB drivers (time.Time, RFC3339/"YYYY-MM-DD HH:MM:SS" strings,
// unix epoch numbers) into a time.Time. Returns the zero time and false when
// the value is absent or unrecognized — callers (e.g. serve_ui_asset) should
// treat the zero value as "no Last-Modified available" rather than an error.
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

// serve_ui_login is the HTTP handler for POST /ui/{ui_slug}/login. It's the
// endpoint the login page's htmx form posts to (application/x-www-form-
// urlencoded, not JSON), and it responds accordingly: no JSON in, no JSON
// out.
//
// It confirms the `ui` exists, then delegates the actual credential check to
// the existing app.dynamic_login (same username/password/JWT-issuing logic
// used elsewhere in the app), against the default "users" table.
//
//   - On success: sets an HttpOnly "session" cookie holding the JWT, and
//     sends an HX-Redirect response header - htmx does a full browser
//     navigation to that URL (which is what you want here anyway, since the
//     cookie only takes effect on the next request/navigation).
//   - On failure: writes a normal HTML error fragment (a DaisyUI alert) with
//     a non-2xx status - no JSON envelope. The login page's form is set up
//     with hx-target="#login-error" hx-swap="innerHTML", plus a small
//     htmx.config.responseHandling tweak so htmx still swaps the body in
//     even though the status is 4xx (by default it wouldn't).
//
// Wire it up alongside the read-only /ui routes - note it's a plain handler
// (not RenderUIPage/RenderUIPartial), since it performs a real auth check
// and sets a cookie rather than rendering a stored template:
//
//	mux.HandleFunc("POST /ui/{ui_slug}/login", app.serve_ui_login)
//
// For the cookie to actually authenticate subsequent requests, uncomment the
// cookie->Authorization fallback already sitting in app.authenticate
// (middleware.go) - it's currently commented out, reading a "session"
// cookie exactly like the one set below.
func (app *application) serve_ui_login(w http.ResponseWriter, r *http.Request) {
	uiSlug := r.PathValue("ui_slug")
	_, err := app.getUser(r)
	if err != nil {
		// w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
		// w.WriteHeader(http.StatusOK)
		http.Redirect(w, r, "/ui/"+r.PathValue("ui_slug"), http.StatusSeeOther)
		return
	}
	params := Dict{
		"lang": "en",
		"data": Dict{
			"db":         env.GetString("UIDB", "UI"),
			"ui_slug":    r.PathValue("ui_slug"),
			"asset_path": r.PathValue("asset"),
		},
	}
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

	if err := r.ParseForm(); err != nil {
		app.writeHTMLError(w, http.StatusBadRequest, "Invalid form submission.")
		return
	}
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if email == "" || password == "" {
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
		app.writeHTMLError(w, http.StatusUnauthorized, msg)
		return
	}
	// 2. Generate a unique session ID
	sessionID := generateSessionID()
	// 3. Save user data to the server-side store
	token, _ := res["token"].(string)
	_data := res["data"].(Dict)
	_data["token"] = token
	app.SessionStore.data.Store(sessionID, _data)
	// 4. Issue the session cookie to the client
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true, // Prevents XSS attacks
		Secure:   true, // Forces HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	// w.Header().Set("HX-Redirect", "/ui/"+uiSlug)
	// w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/ui/"+r.PathValue("ui_slug"), http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintln(w, "Logged out successfully!")
}

// writeHTMLError writes a small HTML fragment (a DaisyUI alert) with the
// given status code - used instead of a JSON body since the login form is
// posted by htmx and expects HTML back, not JSON.
func (app *application) writeHTMLError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `<div class="alert alert-error"><span>%s</span></div>`, template.HTMLEscapeString(msg))
}
