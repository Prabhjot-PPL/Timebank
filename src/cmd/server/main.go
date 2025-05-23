package main

import (
	"log"
	"net/http"
	"timebank/src/internal/adaptors/persistance"
	userhandler "timebank/src/internal/interfaces/input/api/rest/handlers"
	"timebank/src/internal/interfaces/input/api/rest/routes"
	"timebank/src/internal/usecase"
)

func main() {
	database, err := persistance.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	// fmt.Println("Database : ", database)

	UserRepo := persistance.NewUserRepo(database)
	UserService := usecase.NewUserService(UserRepo)
	UserHandler := userhandler.NewUserHandler(UserService)

	router := routes.InitRoutes(UserHandler)

	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Println("Server running on http://localhost:8080")
}
