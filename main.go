package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bars-squad/ais-user-command-service/config"
	"github.com/bars-squad/ais-user-command-service/databases/mongodb"
	"github.com/bars-squad/ais-user-command-service/jwt"
	"github.com/bars-squad/ais-user-command-service/middleware"
	"github.com/bars-squad/ais-user-command-service/modules/admin"
	"github.com/bars-squad/ais-user-command-service/responses"
	"github.com/bars-squad/ais-user-command-service/server"
	"github.com/bars-squad/ais-user-command-service/session"
	"github.com/go-playground/validator"
	rv8 "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" //for development
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cfg                 *config.Config
	httpResponse        = responses.HttpResponseStatusCodesImpl{}
	healthCheckMessage  = "Application running properly"
	pageNotFoundMessage = "You're lost, double check the endpoint"
)

func init() {
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	validate := validator.New()

	mongoClient, err := mongo.NewClient(cfg.Mongodb.ClientOptions)
	if err != nil {
		logger.Fatal(err)
	}
	mongoClientAdapter := mongodb.NewClientAdapter(mongoClient)
	if err := mongoClientAdapter.Connect(context.Background()); err != nil {
		logger.Fatal(err)
	}

	// set redis object
	rdb := rv8.NewClient(cfg.Redis.Options)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Fatal(err)
	}
	// rdb.AddHook(hook.NewAPMRedisHook())

	// set session object
	sess := session.NewRedisSessionStoreAdapter(rdb, time.Hour*24*7, "user.profile")
	// set mongodb
	mongodb := mongoClientAdapter.Database(cfg.Mongodb.Database)

	// set jwt object
	privateKey := jwt.GetRSAPrivateKey("./secret/private.pem")
	publicKey := jwt.GetRSAPublicKey("./secret/public.pem")
	jsonWebToken := jwt.NewJSONWebToken(privateKey, publicKey)

	// set basic auth
	basicAuth := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)

	sessionMiddleware := middleware.NewSessionMiddleware(jsonWebToken)

	router := mux.NewRouter()
	router.HandleFunc("/user-command", index)
	// http.Handle("/", router)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	adminRepository := admin.NewRepository(logger, mongodb)

	adminUsecase := admin.NewUsecase(&admin.Property{
		ServiceName:  cfg.Application.Name,
		Logger:       logger,
		Repository:   adminRepository,
		JSONWebToken: jsonWebToken,
		Session:      sess,
		// Publisher:                publisher,
	})

	admin.NewHTTPHandler(logger, validate, router, basicAuth, adminUsecase, sessionMiddleware)

	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization", "X-RECAPTCHA-TOKEN"},
		AllowCredentials: true,
	}).Handler(router)

	server := server.NewServer(logger, handler, cfg.Application.Port)
	server.Start()

	// When we run this program it will block waiting for a signal. By typing ctrl-C, we can send a SIGINT signal, causing the program to print interrupt and then exit.
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	// closing service for a gracefull shutdown.
	server.Close()
	// redis.Close()
	// db.Close()
	// publisher.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	responses.REST(w, httpResponse.Ok("").NewResponses(nil, healthCheckMessage))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	responses.REST(w, httpResponse.NotFound("PAGES_NOT_FOUND").NewResponses(nil, pageNotFoundMessage))
}
