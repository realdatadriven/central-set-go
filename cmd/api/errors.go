package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"runtime/debug"

	"strings"

	"github.com/realdatadriven/central-set-go/internal/response"
	"github.com/realdatadriven/central-set-go/internal/validator"
)

func GetRelativePath(fullPath, baseDir string) string {
	if !strings.HasPrefix(fullPath, baseDir) {
		return fullPath
	}

	relativePath := strings.TrimPrefix(fullPath, baseDir)
	if relativePath == "" || relativePath == "." {
		return "."
	}

	return relativePath
}

// Function to extract file and line number
func extractFileLine(stack string) string {
	// Regular expression to capture the file and line number from the stack trace
	re := regexp.MustCompile(`\n\t([^\n]+\.go):(\d+)`)
	matches := re.FindAllStringSubmatch(stack, -1)

	var simplifiedStack []string
	for _, match := range matches {
		// Each match will have [full match, file path, line number]
		if len(match) == 3 {
			// We just take the file name and line number
			file := match[1]
			line := match[2]
			simplifiedStack = append(simplifiedStack, fmt.Sprintf("%s:%s", file, line))
		}
	}

	return strings.Join(simplifiedStack, "\n")
}

func (app *application) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)
	fmt.Println(trace)
	requestAttrs := slog.Group("request", "method", method, "url", url)
	//app.logger.Error(message, requestAttrs, "trace", fmt.Sprintf("%s", trace))
	app.logger.Error(message, requestAttrs)

	if app.config.notifications.email != "" {
		data := app.newEmailData()
		data["Message"] = message
		data["RequestMethod"] = method
		data["RequestURL"] = url
		data["Trace"] = trace

		err := app.mailer.Send(app.config.notifications.email, data, "error-notification.tmpl")
		if err != nil {
			trace = string(debug.Stack())
			app.logger.Error(err.Error(), requestAttrs, "trace", trace)
		}
	}
}

func (app *application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]
	err := response.JSONWithHeaders(w, status, map[string]interface{}{"success": false, "msg": fmt.Sprintf("Server Err: %s", message)}, headers)
	if err != nil {
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) invalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	//headers := make(http.Header)
	//headers.Set("WWW-Authenticate", "Bearer")

	//app.errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers)
	data := map[string]interface{}{
		"success": false,
		"msg":     "Token is invalid!",
	}
	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) expiredAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"success": false,
		"msg":     "Token is expired!",
	}
	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) authenticationRequired(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"success": false,
		"msg":     "You must be authenticated to access this resource!",
	}
	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	message := "You must be authenticated to access this resource"
	app.errorMessage(w, r, http.StatusUnauthorized, message, headers)
}
