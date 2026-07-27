package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	texttemplate "text/template"
	"time"
	// OPEN TELEMETRY
	//"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

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
	ui, err := app.AdminGetRowByFilter(
		`select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`,
		[]any{uiSlug, uiSlug},
	)
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
		page, err = app.AdminGetRowByFilter(
			`select * from "ui_page" where ui_id = ? and page_key = ? and active = true and excluded = false`,
			[]any{uiID, pageKey},
		)
	} else {
		page, err = app.AdminGetRowByFilter(
			`select * from "ui_page" where ui_id = ? and default_page = true and active = true and excluded = false`,
			[]any{uiID},
		)
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
	pageID := page["ui_page_id"]

	// 3) active partials for this ui, in a stable order
	partials, err := app.AdminGetRowsByFilter(
		`select * from "ui_partial" where ui_id = ? and active = true and excluded = false order by ui_partial_id asc`,
		[]any{uiID},
	)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_partial: %s", err)}
	}

	// 4) build template data: resolve ui_page_data first, then every active
	// partial's ui_partial_data, each keyed by its own data name.
	tmplData := Dict{
		"UI":         ui,
		"Page":       page,
		"PathParams": pathParams,
	}

	pageDataRows, err := app.AdminGetRowsByFilter(
		`select * from "ui_page_data" where ui_page_id = ? and active = true and excluded = false`,
		[]any{pageID},
	)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_page_data: %s", err)}
	}
	if msg := app.resolveUIData(params, pageDataRows, "ui_page_data", pathParams, tmplData); msg != "" {
		return Dict{"success": false, "msg": msg}
	}

	for _, p := range partials {
		partialDataRows, perr := app.AdminGetRowsByFilter(
			`select * from "ui_partial_data" where ui_partial_id = ? and active = true and excluded = false`,
			[]any{p["ui_partial_id"]},
		)
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
	tset, err := template.New("__page__").Parse(pageTemplate)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to parse page_template: %s", err)}
	}
	for _, p := range partials {
		name, _ := p["ui_partial"].(string)
		body, _ := p["partial_template"].(string)
		if name == "" {
			continue
		}
		if _, err := tset.New(name).Parse(body); err != nil {
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
	ui, err := app.AdminGetRowByFilter(
		`select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`,
		[]any{uiSlug, uiSlug},
	)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui: %s", err)}
	}
	if len(ui) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("ui %q not found", uiSlug)}
	}
	uiID := ui["ui_id"]

	// 2) every active partial for this ui, so cross-partial {{template}} calls
	// keep working; the requested one must be among them.
	partials, err := app.AdminGetRowsByFilter(
		`select * from "ui_partial" where ui_id = ? and active = true and excluded = false order by ui_partial_id asc`,
		[]any{uiID},
	)
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
	}
	for _, p := range partials {
		partialDataRows, perr := app.AdminGetRowsByFilter(
			`select * from "ui_partial_data" where ui_partial_id = ? and active = true and excluded = false`,
			[]any{p["ui_partial_id"]},
		)
		if perr != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui_partial_data: %s", perr)}
		}
		if msg := app.resolveUIData(params, partialDataRows, "ui_partial_data", pathParams, tmplData); msg != "" {
			return Dict{"success": false, "msg": msg}
		}
	}

	// 4) parse every active partial as a named template, then execute only
	// the requested one.
	tset := template.New(partialName)
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
// Dict when sigle_row_obj is true, or as a []Dict otherwise. Returns a
// non-empty error message on failure, "" on success.
func (app *application) resolveUIData(base Dict, rows []Dict, nameCol string, pathParams Dict, tmplData Dict) string {
	for _, row := range rows {
		name, _ := row[nameCol].(string)
		odataPath, _ := row["odata_path"].(string)
		if name == "" || odataPath == "" {
			continue
		}

		resolvedPath, err := app.renderODataPath(odataPath, pathParams)
		if err != nil {
			return fmt.Sprintf("failed to resolve odata_path for %q: %s", name, err)
		}

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

		single := app.toBool(row["sigle_row_obj"])
		switch d := res["data"].(type) {
		case Dict:
			if single {
				tmplData[name] = d
			} else {
				tmplData[name] = []Dict{d}
			}
		case []Dict:
			if single {
				if len(d) > 0 {
					tmplData[name] = d[0]
				} else {
					tmplData[name] = Dict{}
				}
			} else {
				tmplData[name] = d
			}
		default:
			if single {
				tmplData[name] = Dict{}
			} else {
				tmplData[name] = []Dict{}
			}
		}
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
	return out.String(), nil
}

func (app *application) serve_ui_page(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pathParams := Dict{}
	for k := range q {
		pathParams[k] = q.Get(k)
	}
	params := Dict{
		"lang": "en",
		"data": Dict{
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
	params := Dict{
		"lang": "en",
		"data": Dict{
			"ui_slug":      r.PathValue("ui_slug"),
			"partial_name": r.PathValue("partial_name"),
			"path_params":  pathParams,
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

	ui, err := app.AdminGetRowByFilter(
		`select * from "ui" where (ui_slug = ? or ui_name = ?) and active = true and excluded = false`,
		[]any{uiSlug, uiSlug},
	)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to fetch ui: %s", err)}
	}
	if len(ui) == 0 {
		return Dict{"success": false, "msg": fmt.Sprintf("ui %q not found", uiSlug)}
	}
	uiID := ui["ui_id"]

	asset, err := app.AdminGetRowByFilter(
		`select * from "ui_asset" where ui_id = ? and asset_path = ? and active = true and excluded = false`,
		[]any{uiID, assetPath},
	)
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
