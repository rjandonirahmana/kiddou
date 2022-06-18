package main

import (
	"database/sql"
	"fmt"
	"kiddou/base"
	"kiddou/cron"
	grpcVideo "kiddou/grpc/videos"
	"kiddou/handler"
	"kiddou/repo"
	"kiddou/usecase"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	dbStr := fmt.Sprintf("host=localhost port=5444 user=rjandoni password=12345 dbname=kiddou sslmode=disable")
	fmt.Println(dbStr)
	db, err := sql.Open("pgx", dbStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Println(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err == nil {
		mig, err := migrate.NewWithDatabaseInstance("file://migration", "postgres", driver)
		if err == nil {
			err = mig.Up()
			if err != nil && err != migrate.ErrNoChange {
				version, _, _ := mig.Version()
				ver := int(version) - 1
				mig.Force(ver)
				log.Println(err)
				panic(err)

			} else {
				log.Println("Migrate database success")

			}
		} else {
			log.Println(err)

		}
	} else {
		log.Println(err)
	}

	configGoogle := &oauth2.Config{
		ClientID:     "1080827334930-bc6b9e6psejpuds8sk483fq0l8eetihi.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8484/callback-google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	configGithub := &oauth2.Config{
		ClientID:     "eb6ef5373a783a802e95",
		ClientSecret: "9e6f30aa83ae2034e62b690c4f7399b75cc01f4b",
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8484/callback-github",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	redis, err := base.RedisConnection("", "12345", 3)
	if err != nil {
		panic(err)
	}
	defer redis.Close()

	authentication := base.NewRedisAuth(redis)
	repoUser := repo.NewUserRepo(db)
	videoRepo := repo.NewRepositoryVideos(db)
	repoSubcribe := repo.NewRepositorySSub(db)

	usecaseVideo := usecase.NewUsecaseVideos(videoRepo, db, repoSubcribe)
	usecaseUser := usecase.NewUsecaseUser(repoUser, "secretbangett", db, authentication)
	handlerUser := handler.NewUserHandler(usecaseUser)
	videohandler := handler.NewHandlerVideo(usecaseVideo)

	LoginSSO := handler.NewLoginSSO(configGoogle, configGithub, usecaseUser)
	middleware := handler.NewMiddleware(authentication)

	go func() {

		for {
			err := cron.TaskMonitoring(db, repoSubcribe, videoRepo)
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Minute * 10)
			log.Println("sleep for 10 minutes to loop again")
		}
	}()

	go func() {
		listen, err := net.Listen("tcp", ":56767")
		if err != nil {
			panic(err)
		}

		grpcServer := grpc.NewServer()
		grpcVideo.RegisterVideosStreanServer(grpcServer, handler.NewGrpcVideos(usecaseVideo))

		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	app := gin.Default()
	app.LoadHTMLFiles("template")

	app.POST("/register", handlerUser.Register)
	app.POST("/login", handlerUser.Login)

	app.POST("/create/video", middleware.GetTokenFromHeaderBearer(videohandler.CreateVideosAdmin))
	app.POST("/subscription/subscribe", middleware.GetTokenFromHeaderBearer(videohandler.SubscribersVideo))
	app.GET("/subscription/status/:id", middleware.GetTokenFromHeaderBearer(videohandler.StatusSUbscribe))
	app.POST("/subscription/renew", middleware.GetTokenFromHeaderBearer(videohandler.RenewSubscribe))
	app.GET("/api/v1/home", LoginSSO.HomeLogin)
	app.GET("/login-google", LoginSSO.LoginGoogle)
	app.GET("/callback-google", LoginSSO.CallbackGoogleLogin)
	app.Run(":8484")

}
