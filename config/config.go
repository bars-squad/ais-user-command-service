package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Logger struct {
		Formatter logrus.Formatter
	}
	Application struct {
		Port           string
		Name           string
		AllowedOrigins []string
	}
	BasicAuth struct {
		Username string
		Password string
	}
	Mongodb struct {
		ClientOptions *options.ClientOptions
		Database      string
	}
}

func (cfg *Config) mongodb() {
	appName := os.Getenv("APP_NAME")
	uri := os.Getenv("MONGODB_URL")
	db := os.Getenv("MONGODB_DATABASE")
	minPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MIN_POOL_SIZE"), 10, 64)
	maxPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MAX_POOL_SIZE"), 10, 64)
	maxConnIdleTime, _ := strconv.ParseInt(os.Getenv("MONGODB_MAX_IDLE_CONNECTION_TIME_MS"), 10, 64)

	// fmt.Printf("MONGODB_URL\n%s\n\n", uri)
	// fmt.Printf("MONGODB_DATABASE\n%s\n\n", db)

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(time.Millisecond * time.Duration(maxConnIdleTime)).
		SetServerAPIOptions(serverAPIOptions)

	cfg.Mongodb.ClientOptions = opts
	cfg.Mongodb.Database = db
}

func (cfg *Config) basicAuth() {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	cfg.BasicAuth.Username = username
	cfg.BasicAuth.Password = password
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}

func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	rawAllowedOrigins := strings.Trim(os.Getenv("ALLOWED_ORIGINS"), " ")

	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Application.Port = port
	cfg.Application.Name = appName
	cfg.Application.AllowedOrigins = allowedOrigins
}

func Load() *Config {
	cfg := new(Config)
	cfg.app()
	cfg.basicAuth()
	cfg.logFormatter()
	cfg.mongodb()
	return cfg
}
