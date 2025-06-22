package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (app *application) newEmailData() map[string]any {
	data := map[string]any{
		"BaseURL": app.config.baseURL,
	}

	return data
}

func (app *application) backgroundTask(r *http.Request, fn func() error) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.reportServerError(r, fmt.Errorf("%s", err))
			}
		}()

		err := fn()
		if err != nil {
			app.reportServerError(r, err)
		}
	}()
}

func (app *application) contains(slice []interface{}, element interface{}) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func (app *application) joinSlice_(slice []interface{}) string {
	var sb strings.Builder
	for _, v := range slice {
		switch t := v.(type) {
		case int:
			sb.WriteString(strconv.Itoa(t))
		case float64:
			sb.WriteString(strconv.FormatFloat(t, 'f', -1, 64))
		case string:
			sb.WriteString(t)
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return sb.String()
}

func (app *application) joinSlice(slice []interface{}, sep string) string {
	var parts []string
	for _, v := range slice {
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, sep)
}

func (app *application) filter(slice []map[string]interface{}, fn func(map[string]interface{}) bool) []map[string]interface{} {
	filtered := []map[string]interface{}{}
	for _, element := range slice {
		if fn(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}

func (app *application) filterInterface(slice []interface{}, fn func(map[string]interface{}) bool) []map[string]interface{} {
	filtered := []map[string]interface{}{}
	for _, element := range slice {
		if fn(element.(map[string]interface{})) {
			filtered = append(filtered, element.(map[string]interface{}))
		}
	}
	return filtered
}

func (app *application) _map(slice []map[string]interface{}, fn func(map[string]interface{}) map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func (app *application) _map2(slice []interface{}, fn func(interface{}) interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func (app *application) sliceStrs2SliceInterfaces(strs []string) []interface{} {
	interfaces := make([]interface{}, len(strs))
	for i, v := range strs {
		interfaces[i] = v
	}
	return interfaces
}

func (app *application) sliceInterfaces2SliceStrs(strs []interface{}) []string {
	_strings := make([]string, len(strs))
	for i, v := range strs {
		_strings[i] = v.(string)
	}
	return _strings
}

// Create a temporary file in the default temporary directory
func (app *application) tempFIle(content string, name string) (string, error) {
	// Create a temporary file in the default temporary directory
	tempFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %s", err)
	}
	// Defer closing the file to ensure it's closed even if an error occurs
	defer tempFile.Close()

	// Write the content to the file
	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("error writing to temporary file: %s", err)
	}
	// Get the name of the temporary file
	tempFileName := tempFile.Name()
	return tempFileName, nil
}

func (app *application) IsEmpty(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case map[interface{}]interface{}:
		return len(v) == 0
	default:
		return false
	}
}
