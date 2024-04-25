package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"

	db "github.com/erazr/go_tz/db"
	_ "github.com/erazr/go_tz/docs"
	"github.com/erazr/go_tz/models"
)

type CarService struct {
	carRepository db.CarRepository
}

// @title Car API
// @version 1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func RunHttp(ctx context.Context, database *db.Database) {
	carRepository := db.NewCarRepository(database.Conn)
	CarService := CarService{carRepository: carRepository}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	swaggerPath := os.Getenv("SWAGGER_PATH")
	mux.Handle(swaggerPath, httpSwagger.WrapHandler)

	basePath := os.Getenv("BASE_PATH")
	mux.HandleFunc(basePath+"/car/insert", CarService.InsertCar)
	mux.HandleFunc(basePath+"/car/info", CarService.GetCars)
	mux.HandleFunc(basePath+"/car/update", CarService.UpdateCar)
	mux.HandleFunc(basePath+"/car", CarService.GetCarByID)
	mux.HandleFunc(basePath+"/car/owner", CarService.GetCarsByOwner)
	mux.HandleFunc(basePath+"/car/delete", CarService.DeleteCar)

	log.Println("Server started")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()
}

// @Summary Insert new car
// @ID insert-car
// @Tags car
// @Produce json
// @Param car body models.CarRequest true "Car object to be added"
// @Success 200 {object} models.CarResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car/insert [post]
func (c *CarService) InsertCar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var car models.CarRequest
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error(), "Invalid JSON")
		return
	}

	err = c.carRepository.InsertCar(r.Context(), car)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error inserting car")
		return
	}

	json.NewEncoder(w).Encode(car)
}

// @Summary Get cars by filter with pagination
// @ID get-cars-by-filter-with-pagination
// @Tags car
// @Produce json
// @Param filter query string false "Filter by"
// @Param search query string false "Search by"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} models.CarResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car/info [get]
func (c *CarService) GetCars(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	filter := params.Get("filter")
	search := params.Get("search")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	cars, err := c.carRepository.GetCars(r.Context(), models.CarFilter(filter), search, limit, offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error getting car info")
		return
	}
	json.NewEncoder(w).Encode(cars)
}

// @Summary Get car by ID
// @ID get-car-by-id
// @Tags car
// @Produce json
// @Param id query string true "Car ID"
// @Success 200 {object} models.CarResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car [get]
func (c *CarService) GetCarByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	id := params.Get("id")

	car, err := c.carRepository.GetCarByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error getting car by ID")
		return
	}
	json.NewEncoder(w).Encode(car)
}

// @Summary Get cars by owner
// @ID get-cars-by-owner
// @Tags car
// @Produce json
// @Param id query string true "Owner ID"
// @Success 200 {object} models.CarResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car/owner [get]
func (c *CarService) GetCarsByOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	ownerID := params.Get("id")

	cars, err := c.carRepository.GetCarsByOwner(r.Context(), ownerID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error getting cars by owner")
		return
	}
	json.NewEncoder(w).Encode(cars)
}

// @Summary Update car info
// @ID update-car
// @Tags car
// @Produce json
// @Param car body models.CarRequest true "Car object to be updated"
// @Param id query string true "Car ID"
// @Success 200 {object} models.CarResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car/update [put]
func (c *CarService) UpdateCar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	id := params.Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var car models.CarResponse
	car.ID = id

	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error(), "Invalid JSON")
		return
	}

	err = c.carRepository.UpdateCar(r.Context(), car)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error updating car")
		return
	}

	json.NewEncoder(w).Encode(car)
}

// @Summary Delete car
// @ID delete-car
// @Tags car
// @Produce json
// @Param id query string true "Car ID"
// @Success 200 {string} string "Car deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /car/delete [delete]
func (c *CarService) DeleteCar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	id := params.Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := c.carRepository.DeleteCar(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err.Error(), "Error deleting car")
		return
	}

	json.NewEncoder(w).Encode("Car deleted")
}
