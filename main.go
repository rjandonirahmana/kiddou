package main

import (
	"database/sql"
	"fmt"
	"kiddou/base"
	"kiddou/cron"
	"kiddou/handler"
	"kiddou/repo"
	"kiddou/usecase"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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

			} else {
				log.Println("Migrate database success")

			}
		} else {
			log.Println(err)

		}
	} else {
		log.Println(err)
	}

	redis, err := base.RedisConnection("", "12345", 3)
	if err != nil {
		panic(err)
	}

	authentication := base.NewRedisAuth(redis)
	repoUser := repo.NewUserRepo(db)
	videoRepo := repo.NewRepositoryVideos(db)
	repoSubcribe := repo.NewRepositorySSub(db)

	usecaseVideo := usecase.NewUsecaseVideos(videoRepo, db, repoSubcribe)
	usecaseUser := usecase.NewUsecaseUser(repoUser, "secretbangett", db, authentication)
	handlerUser := handler.NewUserHandler(usecaseUser)
	videohandler := handler.NewHandlerVideo(usecaseVideo)

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

	app := gin.Default()

	app.POST("/register", handlerUser.Register)
	app.POST("/login", handlerUser.Login)

	app.POST("/create/video", middleware.GetTokenFromHeaderBearer(videohandler.CreateVideosAdmin))
	app.POST("/subscription/subscribe", middleware.GetTokenFromHeaderBearer(videohandler.SubscribersVideo))
	app.GET("/subscription/status/:id", middleware.GetTokenFromHeaderBearer(videohandler.StatusSUbscribe))
	app.POST("/subscription/renew", middleware.GetTokenFromHeaderBearer(videohandler.RenewSubscribe))
	app.Run(":8282")

}
