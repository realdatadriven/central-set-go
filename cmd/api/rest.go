package main

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/assets"
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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/realdatadriven/central-set-go/internal/response"
)

func (app *application) crud_api(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	id := r.PathValue("id")
	user := app.getAnonymous()
	userFromSess, err := app.getUser(r)
	if err != nil {
		user = userFromSess
	}
	params := Dict{
		"lang": "en",
		"user": user,
		"data": Dict{
			"db":    db,
			"table": table,
		},
	}
	// GET follows the existing OData read pattern.
	if r.Method == http.MethodGet {
		odataPath := fmt.Sprintf("%s/%s/", db, table)
		if id != "" {
			// Resolve the primary key and turn /{id} into an OData filter.
			pk, err := app.crudPrimaryKey(db, table)
			if err != nil {
				_ = response.JSON(w, http.StatusBadRequest, Dict{
					"success": false,
					"msg":     err.Error(),
				})
				return
			}
			odataPath = fmt.Sprintf("%s/%s/?$filter=%s eq %s",db,table,pk,id)
		} else if r.URL.RawQuery != "" {
			odataPath = fmt.Sprintf("%s/%s/?%s",db,table,r.URL.RawQuery)
		}
		data := app.ODataRead(params, odataPath)
		_ = response.JSON(w, http.StatusOK, data)
		return
	}
	// Everything other than GET becomes create_update.
	data, err := app.crudRequestData(r)
	if err != nil {
		_ = response.JSON(w, http.StatusBadRequest, Dict{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}
	// For PUT/PATCH/DELETE /crud/{db}/{table}/{id},
	// put the URL id into the table's actual primary-key field.
	if id != "" {
		pk, err := app.crudPrimaryKey(db, table)
		if err != nil {
			_ = response.JSON(w, http.StatusBadRequest, Dict{
				"success": false,
				"msg":     err.Error(),
			})
			return
		}
		data[pk] = id
	}
	// DELETE uses the existing C7 CRUD deletion mechanism.
	if r.Method == http.MethodDelete {
		data["_to_delete"] = true
	}
	params["data"].(Dict)["data"] = data
	result := app.create_update(params)
	status := http.StatusOK
	if success, ok := result["success"].(bool); ok && !success {
		status = http.StatusBadRequest
	}
	_ = response.JSON(w, status, result)
}