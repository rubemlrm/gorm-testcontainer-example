package main

import (
	"github.com/go-faker/faker/v4"
	"gorm-test/config"
	"gorm-test/internal/models"
	repositories "gorm-test/internal/repositories"
	"gorm-test/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := database.StartConnection(cfg.Database)

	repository := repositories.NewRepository(conn)

	// random book for testing
	var book models.Book
	err = faker.FakeData(&book)
	if err != nil {
		panic(err)
	}

	result := repository.Insert(&book)
	println(result)
}
