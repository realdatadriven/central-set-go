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

func (app *application) contains(slice []any, element any) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func (app *application) joinSlice_(slice []any) string {
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

func (app *application) joinSlice(slice []any, sep string) string {
	var parts []string
	for _, v := range slice {
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, sep)
}

func (app *application) filter(slice []map[string]any, fn func(map[string]any) bool) []map[string]any {
	filtered := []map[string]any{}
	for _, element := range slice {
		if fn(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}

func (app *application) filterAny(slice []any, fn func(any) bool) []any {
	filtered := []any{}
	for _, element := range slice {
		if fn(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}

func (app *application) String2Any(s string) any {
	return s
}

func (app *application) filterInterface(slice []any, fn func(map[string]any) bool) []map[string]any {
	filtered := []map[string]any{}
	for _, element := range slice {
		if fn(element.(map[string]any)) {
			filtered = append(filtered, element.(map[string]any))
		}
	}
	return filtered
}

func (app *application) _map(slice []map[string]any, fn func(map[string]any) map[string]any) []map[string]any {
	result := make([]map[string]any, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func (app *application) _map2(slice []any, fn func(any) any) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func (app *application) sliceStrs2SliceInterfaces(strs []string) []any {
	interfaces := make([]any, len(strs))
	for i, v := range strs {
		interfaces[i] = v
	}
	return interfaces
}

func (app *application) sliceInterfaces2SliceStrs(strs []any) []string {
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

func (app *application) IsEmpty(value any) bool {
	switch v := value.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	default:
		return false
	}
}
