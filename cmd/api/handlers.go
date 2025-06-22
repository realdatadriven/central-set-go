package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/realdatadriven/central-set-go/internal/password"
	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"
	"github.com/realdatadriven/central-set-go/internal/validator"

	"github.com/pascaldekloe/jwt"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	/*_, found, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}*/

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	//input.Validator.CheckField(!found, "Email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	/*hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.InsertUser(input.Email, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}*/

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	/*user, found, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}*/

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	/*input.Validator.CheckField(found, "Email", "Email address could not be found")

	if found {
		passwordMatches, err := password.Matches(input.Password, user.Password)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}*/

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	var claims jwt.Claims
	//claims.Subject = strconv.Itoa(user.UserID)

	expiry := time.Now().Add(24 * time.Hour)
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

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(os.IsNotExist(err), err)
		return !os.IsNotExist(err)
	}
	return true
}

// fileExistsInS3 checks if a file exists in the given S3 bucket
func (app *application) fileExistsInS3(svc *s3.S3, bucket, key string) bool {
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	// If no error, the file exists
	return err == nil
}

// AWS SESSION
func (app *application) awsSession() (*session.Session, error) {
	awsConfig := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"), // Optional
		),
		Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")), // Optional custom endpoint,
		S3ForcePathStyle: aws.Bool(app.config.s3ForcePathStyle), // Force path-style URLs (necessary for MinIO)
		DisableSSL:       aws.Bool(app.config.s3DisableSSL),     // MinIO often runs without SSL locally
	}
	// Create a custom HTTP client that skips SSL certificate verification
	if app.config.s3SkipSSLVerify && !app.config.s3DisableSSL {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // Disable SSL verification
		awsConfig = &aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				os.Getenv("AWS_SESSION_TOKEN"), // Optional
			),
			Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),    // Optional custom endpoint,
			S3ForcePathStyle: aws.Bool(app.config.s3ForcePathStyle),    // Force path-style URLs (necessary for MinIO)
			DisableSSL:       aws.Bool(app.config.s3DisableSSL),        // MinIO often runs without SSL locally
			HTTPClient:       &http.Client{Transport: customTransport}, // Use the custom transport with TLS config
		}
	}
	// Create AWS session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}
	return sess, nil
}

// uploadToS3 uploads a file to the configured S3 bucket
func (app *application) uploadToS3(file io.Reader, fileName string) (string, error) {
	// Create AWS session
	sess, err := app.awsSession()
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}
	// Create S3 service client
	svc := s3.New(sess)
	// Define the S3 bucket and key
	bucket := os.Getenv("S3_BUCKET")
	originalKey := fileName
	ext := filepath.Ext(originalKey)
	baseName := fileName[:len(fileName)-len(ext)]
	// Check if the file already exists and modify the file name if necessary
	key := originalKey
	for i := 1; app.fileExistsInS3(svc, bucket, key); i++ {
		key = fmt.Sprintf("%s_%d%s", baseName, i, ext)
	}
	// Read file into a buffer to allow seeking
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return "", fmt.Errorf("failed to read file into buffer: %v", err)
	}
	// Convert buffer into a ReadSeeker
	fileReader := bytes.NewReader(buffer.Bytes())
	// Upload file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileReader,
		//ACL:    aws.String("public-read"), // Optional: Set ACL for public access if needed
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}
	return key, nil
}

func (app *application) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Get the multipart form
	err := r.ParseMultipartForm(int64(app.config.uploadSize)) // 100MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	// is to save in temp folder
	lang := r.FormValue("lang")
	err = app.i18n.ChangeLanguage(lang)
	if err != nil {
		fmt.Println("Err Load Lang:", err)
	}
	tmp := r.FormValue("tmp")
	//print("Tmp: ", tmp, "lang: ", lang, "\n")
	// get the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Get the file name and extension
	fileName := filepath.Base(handler.Filename)
	fileExt := filepath.Ext(handler.Filename)
	fileNameNoExt := fileName[:len(fileName)-len(fileExt)]
	//fmt.Printf("Name: %s Ext: %s \n", fileNameNoExt, fileExt)
	var dst *os.File
	// var err error
	// create the fale
	boolTemp, _ := strconv.ParseBool(tmp)
	// USE S3
	var data map[string]interface{}
	if app.config.useS3 && !boolTemp {
		fmt.Println("IS S3")
		// Upload to S3
		fname, err := app.uploadToS3(file, fileNameNoExt+fileExt)
		if err != nil {
			fmt.Println(err.Error())
			data = map[string]interface{}{
				"success": false,
				"msg":     "Failed to upload to S3: " + err.Error(),
			}
			err = response.JSON(w, http.StatusOK, data)
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}
		text, _ := i18n.T("file-success", struct{ File string }{File: fname})
		data = map[string]interface{}{
			"success": true,
			"msg":     text,
			"file":    fname,
		}
		err = response.JSON(w, http.StatusOK, data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}
	if boolTemp {
		//print("parsed temp bool:", boolTemp, "\n")
		dst, err = os.CreateTemp("", fileNameNoExt+"-*"+fileExt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
	} else {
		//dst, err = os.Create("static/uploads/" + handler.Filename)
		//_path := "static/uploads/" + fileNameNoExt + "" + fileExt
		_path := fmt.Sprintf("%s/%s%s", app.config.upload_path, fileNameNoExt, fileExt)
		if app.fileExists(_path) {
			for i := 1; i <= 100; i++ {
				_path := fmt.Sprintf("%s/%s_%d%s", app.config.upload_path, fileNameNoExt, i, fileExt)
				if !app.fileExists(_path) {
					dst, err = os.Create(_path)
					break
				}
			}
		} else {
			dst, err = os.Create(_path)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
	}
	//print(dst.Name(), "\n")
	// Save the file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//user := *(contextGetAuthenticatedUser(r))
	//print(user["username"].(string), "->", int(user["user_id"].(float64)), "->", int(user["role_id"].(float64)), "\n")
	text, _ := i18n.T("file-success", struct{ File string }{File: filepath.Base(dst.Name())})
	data = map[string]interface{}{
		"success": true,
		"msg":     text,
		"file":    filepath.Base(dst.Name()),
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
