package main

import (
	"database/sql"
	"fmt"
	"kiddou/base"
	"kiddou/handler"
	"kiddou/repo"
	"kiddou/usecase"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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
	usecaseUser := usecase.NewUsecaseUser(repoUser, "secretbangett", db, authentication)
	handlerUser := handler.NewUserHandler(usecaseUser)

	fiber := fiber.New(fiber.Config{})

	fiber.Post("/register", handlerUser.Register)

	log.Fatal(fiber.Listen(":8080"))

}
