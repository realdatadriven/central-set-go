package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/realdatadriven/central-set-go/internal/response"
	"github.com/realdatadriven/etlx"
)

func (app *application) crudRequestData(r *http.Request) (Dict, error) {
	data := Dict{}
	// JSON request
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("invalid JSON body: %w", err)
		}
		return data, nil
	}
	// HTMX / HTML form request.
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}
	for key, values := range r.PostForm {
		if len(values) == 0 {
			data[key] = ""
			continue
		}
		// Preserve multiple values as []any.
		if len(values) == 1 {
			data[key] = values[0]
			continue
		}
		items := make([]any, len(values))
		for i, value := range values {
			items[i] = value
		}
		data[key] = items
	}
	return data, nil
}

func (app *application) crudPrimaryKey(params Dict, db, table string) (string, error) {
	dsn, _, err := app.GetDBNameFromParams(params)
	if err != nil {
		return "", fmt.Errorf("failed to query for PK Field %s %s: %w", db, table, err)
	}
	dbConn, err := etlx.GetDB(dsn)
	if err != nil {
		return "", fmt.Errorf("failed to query for PK Field %s %s: %w", db, table, err)
	}
	etlx_engine := &etlx.ETLX{}
	dialect := etlx_engine.GetDialect(dbConn.GetDriverName())
	pksql := dialect.GetPrimaryKeyAutoIncrementQuery(table)
	pk, _, err := dbConn.QuerySingleRow(pksql, []any{}...)
	if err != nil {
		return "", fmt.Errorf("failed to query for PK Field %s %s: %w", table, pksql, err)
	}
	if len(*pk) == 0 {
		return "", fmt.Errorf("failed to get the PK Field %s %s", table, pksql)
	}
	pkfield, ok := (*pk)["column_name"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get the PK Field %s %s", table, pksql)
	}
	return pkfield, nil
}

func (app *application) getPKFromSchema(db, table string) (string, error) {
	sql := `select field as column_name from table_schema where db = ? and "table" = ? and pk = true`
	pk, err := app.AdminGetRowByFilter(sql, []any{db, table})
	if err != nil {
		return "", fmt.Errorf("failed to query for PK Field %s %s: %w", db, table, err)
	}
	if len(pk) == 0 {
		return "", fmt.Errorf("failed to get the PK Field %s %s", db, table)
	}
	pkfield, ok := pk["column_name"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get the PK Field %s %s", db, table)
	}
	return pkfield, nil
}

func (app *application) crud_api_handler(w http.ResponseWriter, r *http.Request) {
	crud := app.crud_api(w, r)
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		if err := response.JSON(w, http.StatusOK, crud); err != nil {
			app.serverError(w, r, err)
		}
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if success, ok := crud["success"].(bool); ok && !success {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		fmt.Fprintf(w, `%s`, template.HTMLEscapeString(crud["msg"].(string)))
	}
	return
}

func (app *application) crud_api(w http.ResponseWriter, r *http.Request) Dict {
	db := r.PathValue("db")
	table := r.PathValue("table")
	id := r.PathValue("id")
	q := r.URL.Query()
	pathParams := Dict{}
	for k := range q {
		pathParams[k] = q.Get(k)
	}
	user := app.getAnonymous()
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		token := app.verifyToken(r)
		if !token["success"].(bool) {
			// return token
		} else if _user := contextGetAuthenticatedUser(r); _user != nil {
			user = *_user
		}
	} else {
		userFromSess, err := app.getUser(r)
		if err != nil {
			user = userFromSess
		}
	}
	sql := `select * from app where (app = ? or db = ?) and excluded = false`
	_app, err := app.AdminGetRowByFilter(sql, []any{db, db})
	if err != nil {
		return Dict{
			"success": false,
			"msg":     "Failed to fetch APP!",
		}
	}
	if len(_app) == 0 {
		return Dict{
			"success": false,
			"msg":     fmt.Sprintf("APP associated with db %s not found!", db),
		}
	}
	params := Dict{
		"lang": "en",
		"user": user,
		"app":  _app,
		"data": Dict{
			"db":         db,
			"table":      table,
			"pathParams": pathParams,
		},
	}
	loc := app.getLocationFromRequest(r, params)
	params["location"] = loc
	lang := "en"
	if _, ok := params["lang"]; ok {
		lang = params["lang"].(string)
	}
	err = app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println(err)
	}
	// GET follows the existing OData read pattern.
	if r.Method == http.MethodGet {
		odataPath := fmt.Sprintf("%s/%s", db, table)
		if id != "" {
			// Resolve the primary key and turn /{id} into an OData filter.
			pk, err := app.getPKFromSchema(db, table)
			if err != nil {
				return Dict{
					"success": false,
					"msg":     err.Error(),
				}
			}
			odataPath = fmt.Sprintf("%s/%s/?$filter=%s eq %s", db, table, pk, id)
		} else if r.URL.RawQuery != "" {
			odataPath = fmt.Sprintf("%s/%s/?%s", db, table, r.URL.RawQuery)
		}
		data := app.ODataRead(params, odataPath)
		return data
	}
	// Everything other than GET becomes create_update.
	data, err := app.crudRequestData(r)
	if err != nil && r.Method != http.MethodDelete {
		return Dict{
			"success": false,
			"msg":     err.Error(),
		}
	}
	// For PUT/PATCH/DELETE /crud/{db}/{table}/{id},
	// put the URL id into the table's actual primary-key field.
	if id != "" {
		pk, err := app.getPKFromSchema(db, table)
		if err != nil {
			return Dict{
				"success": false,
				"msg":     err.Error(),
			}
		}
		data[pk] = id
		if r.Method == http.MethodDelete {
			data["exclided"] = true
		}
	} else if r.Method == http.MethodDelete || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		// For DELETE /crud/{db}/{table}, require the primary key in the request body.
		return Dict{
			"success": false,
			"msg":     "Missing primary key id for the action.",
		}
	}
	// DELETE uses the existing C7 CRUD deletion mechanism.
	if r.Method == http.MethodDelete {
		data["_to_delete"] = true
	}
	params["data"].(Dict)["data"] = data
	return app.create_update(params)
}
