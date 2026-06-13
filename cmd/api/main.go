package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/realdatadriven/central-set-go/internal/auth"
	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/smtp"
	"github.com/realdatadriven/central-set-go/internal/version"
	"github.com/robfig/cron/v3"

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

type app_config struct {
	quackEnabled bool
	baseURL      string
	httpPort     int
	basicAuth    struct {
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
	s3Endpoint            string
	frontend_url          string
	LockoutEnabled        bool
	LockoutThreshold      int
}

//type admin struct{}

/*type user struct{
	user_id int
	role_id int
}*/

type application struct {
	config app_config
	db     etlx.DBInterface //*etlx.DB
	//memdb                          etlx.DBInterface //*etlx.DB
	rateLimitingEnabled bool
	memdb               *sql.DB
	rtRequestLimit      int

	logger                         *slog.Logger
	mailer                         *smtp.Mailer
	wg                             sync.WaitGroup
	i18n                           i18next.I18n
	appType                        string        // can be community, licensor or licensee
	lastLicenseValidation          time.Time     // time of last license validation
	licenceVerificationPeriodicity time.Duration // periodicity of license validation
	SSE_Broker                     *Broker
	WS_ConnectionManager           *ConnectionManager
	//user user
	//admin  admin
	cronScheduler *cron.Cron
	cronEntries   map[any]cron.EntryID
	cronEntriesMu sync.Mutex
	lastCronCheck time.Time

	// Quack server pool
	quackEnabled      bool                     // Whether Quack is enabled
	quackPool         map[int]etlx.DBInterface // In-memory DuckDB instances keyed by quack_server_id
	quackPoolMux      sync.RWMutex             // Protect concurrent access to pool
	quackManager      *QuackManager            // Quack lifecycle manager
	quackInstanciated bool                     // Quack instanciated flag to avoid multiple instantiation of quack manager and pool in case of multiple calls to run function, as it can happen in licensee app where we validate the license on startup and then periodically, and both operations call run function
}

func run(logger *slog.Logger) error {
	var cfg app_config
	// Load environment variables
	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "f2rkbev2yxhk5viz77ok4rxfip6npjpm")
	cfg.db.driverName = env.GetString("DB_DRIVER_NAME", "sqlite3")
	cfg.db.dsn = env.GetString("DB_DSN", "database/ADMIN.db")
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
	cfg.s3Endpoint = env.GetString("AWS_ENDPOINT", "")
	cfg.frontend_url = env.GetString("FRONTEND_URL", "http://localhost:4444")
	cfg.LockoutEnabled = env.GetBool("LOCKOUT_ENABLED", true)
	cfg.LockoutThreshold = env.GetInt("LOCKOUT_THRESHOLD", 3)
	cfg.quackEnabled = env.GetBool("QUACK_ENABLED", false)
	//cli flags
	showVersion := flag.Bool("version", false, "display version and exit")
	initdb := flag.Bool("init", false, "initialize the main db")
	dbname := flag.String("dbname", "ADMIN", "initialize the main db")
	model := flag.String("model", "admin_model.md", "initialize the db with the provided model (only used if init flag is set)")
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
		// CS_LICENCOR_TOKEN and CS_LICENCOR_URL env variables must be set for licensee appType
		appType:                        "community",                     // can be community, licensor or licensee
		lastLicenseValidation:          time.Now().Add(-24 * time.Hour), // in a licencee app, we will validate the license on startup and, as its by default 24 hours periodicity, we set last validation to 24 hours ago
		licenceVerificationPeriodicity: 24 * time.Hour,
		quackEnabled:                   cfg.quackEnabled,
		//admin:  admin{},
	}
	app.rateLimitingEnabled = env.GetBool("RATE_LIMITING", false)
	if app.rateLimitingEnabled {
		app.rtRequestLimit = env.GetInt("RATE_LIMITING_REQUESTS_PER_MINUTE", 100)
		fmt.Printf("Rate limiting is enabled with request limit: %d\n", app.rtRequestLimit)
		//app.memdb, err = etlx.New("duckdb:", ":memory:")
		rtLimitPath := env.GetString("RATE_LIMITING_DB_PATH", "file::memory:?cache=shared")
		app.memdb, err = sql.Open("sqlite3", rtLimitPath)
		if err != nil {
			fmt.Printf("Error setting mem db: %s: %v\n", rtLimitPath, err)
			return err
		}
		defer app.memdb.Close()
		/*/ set tread number to 1 for duckdb to avoid concurrency issues as it's used as in-memory db for license validation and other operations that are not performance critical
		rtThreads := env.GetInt("RATE_LIMITING_DB_THREADS", 1)
		_, err = app.memdb.Exec(fmt.Sprintf("SET threads=%d;", rtThreads))
		if err != nil {
			fmt.Printf("Error setting duckdb threads to 1: %v\n", err)
		}
		// set duckdb memory limit to 1GB to avoid it consuming too much memory as it's used as in-memory db for license validation and other operations that are not performance critical
		memLimit := env.GetString("RATE_LIMITING_DB_MEMORY_LIMIT", "1GB")
		_, err = app.memdb.Exec(fmt.Sprintf("SET memory_limit = '%s';", memLimit))
		if err != nil {
			fmt.Printf("Error setting duckdb memory limit: %v\n", err)
		}*/
		app.memdb.Exec("PRAGMA journal_mode=WAL;")
		busy_timeout := 5_1000
		app.memdb.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", busy_timeout))
		sql := `CREATE TABLE IF NOT EXISTS rate_limits (ip TEXT PRIMARY KEY, request_count INTEGER, last_request_time TIMESTAMP)`
		_, err := app.memdb.Exec(sql)
		if err != nil {
			return err
		}
	}
	// OAUTH INIT IF ENABLED
	if env.GetBool("ENABLE_OAUTH", false) {
		//fmt.Println("ENABLE_OAUTH:", true)
		auth.InitGoth()
	}
	// golang get current time - 24 hours
	//app.lastLicenseValidation = time.Now().Add(-24 * time.Hour)
	// Set tenant environment variables
	sql := `select * from "env" where "active" = ? and "on_srv_start" = ? and "excluded" = ?`
	tenantEnv, err := app.AdminGetRowsByFilter(sql, []any{true, true, false})
	if err != nil {
		fmt.Printf("Error fetching tenant env vars: %v\n", err)
	} else {
		for _, v := range tenantEnv {
			os.Setenv(v["env_name"].(string), v["env_value"].(string))
			//fmt.Printf("Setting env var for admin %s=%s\n", v["env_name"], v["env_value"])
		}
	}
	// err = db.Ping()
	if *initdb /*&& err != nil*/ {
		if *model == "" {
			fname := fmt.Sprintf(`%s.%s.sql`, *dbname, db.GetDriverName())
			err := app.setupDB(fname, *dbname, *embedded)
			if err != nil {
				fmt.Printf("error setingup the DB: %v\n", err)
			}
		} else {
			err := app.setupWithModel(*model)
			if err != nil {
				fmt.Printf("error setingup the DB with model %s: %v\n", *model, err)
			}
		}
		return nil
	}
	app.CronJobs()
	if env.GetBool("ENABLE_ARROW_FLIGHT", false) {
		go func() {
			err := app.serveArrowFlight()
			if err != nil {
				fmt.Printf("Error setting up arrow flight server: %v\n", err)
			}
		}()
	}
	if env.GetBool("SSE_ENABLE", false) {
		go func() {
			err := app.serveSSE()
			if err != nil {
				fmt.Printf("Error setting up SSE server: %v\n", err)
			}
		}()
	}
	return app.serveHTTP()
}

// go run ./cmd/api
