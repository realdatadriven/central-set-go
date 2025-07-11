package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/smtp"
	"github.com/realdatadriven/central-set-go/internal/version"

	"github.com/lmittmann/tint"
	"github.com/realdatadriven/etlx"

	"github.com/joho/godotenv"
	"github.com/yuangwei/go-i18next"
)

var i18n i18next.I18n

func main() {
	// Load .env file
	_err := godotenv.Load()
	if _err != nil {
		slog.Error("Error loading .env file")
	}
	//httpPort := os.Getenv("HTTP_PORT")
	//fmt.Printf("HTTP_PORT: %s\n", httpPort)
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	cookie struct {
		secretKey string
	}
	db struct {
		dsn         string
		driverName  string
		automigrate bool
	}
	jwt struct {
		secretKey        string
		tokenExpireHours int
	}
	notifications struct {
		email string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
	uploadSize            int
	upload_path           string
	core_tables           string
	enable_user           string
	enable_app            string
	actions_not_to_log    string
	broadcast_changes     string
	allow_cli_run_queries bool
	useS3                 bool
	s3Bucket              string
	s3Region              string
	s3ForcePathStyle      bool // Force path-style URLs (necessary for MinIO)
	s3DisableSSL          bool
	s3SkipSSLVerify       bool
}

//type admin struct{}

/*type user struct{
	user_id int
	role_id int
}*/

type application struct {
	config config
	db     etlx.DBInterface //*etlx.DB
	logger *slog.Logger
	mailer *smtp.Mailer
	wg     sync.WaitGroup
	i18n   i18next.I18n
	//user user
	//admin  admin
}

func run(logger *slog.Logger) error {
	var cfg config
	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "f2rkbev2yxhk5viz77ok4rxfip6npjpm")
	cfg.db.driverName = env.GetString("DB_DRIVER_NAME", "sqlite3")
	cfg.db.dsn = env.GetString("DB_DSN", "database/admin.db")
	//fmt.Printf("DB_DRIVER_NAME: %s DB_DRIVER_NAME: %s\n", cfg.db.dsn, cfg.db.driverName)
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "mhaitpm4v3mesosefepyupo6qzpbvidc")
	cfg.jwt.tokenExpireHours = env.GetInt("TOKEN_EXPIRE_HOURS", 24)
	cfg.notifications.email = env.GetString("NOTIFICATIONS_EMAIL", "")
	cfg.smtp.host = env.GetString("SMTP_HOST", "example.smtp.host")
	cfg.smtp.port = env.GetInt("SMTP_PORT", 25)
	cfg.smtp.username = env.GetString("SMTP_USERNAME", "example_username")
	cfg.smtp.password = env.GetString("SMTP_PASSWORD", "pa55word")
	cfg.smtp.from = env.GetString("SMTP_FROM", "Example Name <no_reply@example.org>")
	cfg.uploadSize = env.GetInt("UPLOAD_SIZE", 10<<20)
	cfg.enable_app = env.GetString("ENABLE_APP", "app,role_app,role_app_menu,role_app_menu_table")
	cfg.enable_user = env.GetString("ENABLE_USER", "user_role,column_level_access")
	cfg.core_tables = env.GetString("CORE_TABLES", "user_role,column_level_access")
	cfg.actions_not_to_log = env.GetString("ACTIONS_NOT_TO_LOG", "")
	cfg.broadcast_changes = env.GetString("BROADCAST_CHANGES", "")
	cfg.upload_path = env.GetString("UPLOAD", "static/uploads")
	cfg.allow_cli_run_queries = env.GetBool("ALLOW_CLI_RUN_QUERIES", false)
	cfg.useS3 = env.GetBool("USE_S3_STORAGE", false)
	cfg.s3Bucket = env.GetString("S3_BUCKET", "uploads")
	cfg.s3Region = env.GetString("S3_REGION", "")
	cfg.s3ForcePathStyle = env.GetBool("S3_FORCE_PATH_STYLE", true)
	cfg.s3DisableSSL = env.GetBool("S3_DISABLE_SSL", false)
	cfg.s3SkipSSLVerify = env.GetBool("S3_SKIP_SSL_VERIFY", false)
	showVersion := flag.Bool("version", false, "display version and exit")
	//cli flags
	initdb := flag.Bool("init", false, "initialize the main db")
	dbname := flag.String("dbname", "ADMIN", "initialize the main db")
	embedded := flag.Bool("embedded", true, "use the embedded db")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}
	//db, err := database.New(cfg.db.driverName, cfg.db.dsn, cfg.db.automigrate)
	db, err := etlx.New(cfg.db.driverName, cfg.db.dsn)
	//db, err := etlx.GetDB(cfg.db.driverName)
	if err != nil {
		return err
	}
	defer db.Close()
	mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	if err != nil {
		return err
	}
	// i18n
	i18n, err = i18next.Init(i18next.I18nOptions{
		Lng:        []string{"en", "pt"},
		DefaultLng: "en",
		Ns:         "yaml",
		Backend: i18next.Backend{
			LoadPath: []string{"./locales/{{.Lng}}.yml"},
		},
	})
	if err != nil {
		fmt.Println("Err: i18n: ", err)
	}
	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		mailer: mailer,
		i18n:   i18n,
		//admin:  admin{},
	}
	// err = db.Ping()
	if *initdb /*&& err != nil*/ {
		fname := fmt.Sprintf(`%s.%s.sql`, *dbname, db.GetDriverName())
		err := app.setupDB(fname, *dbname, *embedded)
		if err != nil {
			return fmt.Errorf("error setingup the DB: %v\n", err)
		}
	}
	app.CronJobs()
	return app.serveHTTP()
}

// go run ./cmd/api
