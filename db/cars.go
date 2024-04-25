package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/erazr/go_tz/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CarRepository interface {
	InsertCar(ctx context.Context, car models.CarRequest) error
	GetCars(ctx context.Context, filter models.CarFilter, search string, limit int, offset int) ([]models.CarResponse, error)
	GetCarByID(ctx context.Context, id string) (models.CarResponse, error)
	GetCarsByOwner(ctx context.Context, ownerID string) ([]models.CarResponse, error)
	UpdateCar(ctx context.Context, car models.CarResponse) error
	DeleteCar(ctx context.Context, id string) error
}

type carRepository struct {
	db *pgx.Conn
}

func NewCarRepository(db *pgx.Conn) *carRepository {
	return &carRepository{
		db: db,
	}
}

func (r *carRepository) InsertCar(ctx context.Context, car models.CarRequest) error {
	var personID string
	err := r.db.QueryRow(ctx, "SELECT id FROM person WHERE name = $1 AND surname = $2", car.Owner.Name, car.Owner.Surname).Scan(&personID)
	if err != nil {
		if err == pgx.ErrNoRows {
			personID = uuid.NewString()
			_, err = r.db.Exec(ctx, "INSERT INTO person (id, name, surname) VALUES ($1, $2, $3)", personID, car.Owner.Name, car.Owner.Surname)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	carID, _ := uuid.NewRandom()
	_, err = r.db.Exec(ctx, "INSERT INTO cars (id, reg_number, mark, model, year, owner) VALUES ($1, $2, $3, $4, $5, $6)", carID, car.RegNum, car.Mark, car.Model, car.Year, personID)
	if err != nil {
		return err
	}

	return nil
}

func (r *carRepository) GetCars(ctx context.Context, filter models.CarFilter, search string, limit int, offset int) ([]models.CarResponse, error) {
	var cars []models.CarResponse
	var args []interface{}
	query := "SELECT c.id, c.reg_number, c.mark, c.model, c.year, p.id, p.name, p.surname FROM cars c INNER JOIN person p ON c.owner = p.id"

	if filter != "" && search != "" {
		query += fmt.Sprintf(" WHERE c.%s = $1", filter)
		args = append(args, search)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cars not found")
	}
	defer rows.Close()

	for rows.Next() {
		var car models.CarResponse
		err := rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.ID, &car.Owner.Name, &car.Owner.Surname)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

func (r *carRepository) GetCarByID(ctx context.Context, id string) (models.CarResponse, error) {
	var car models.CarResponse
	err := r.db.QueryRow(ctx, `
		SELECT c.id, c.reg_number, c.mark, c.model, c.year, p.id, p.name, p.surname 
		FROM cars c
		INNER JOIN person p ON c.owner = p.id
		WHERE c.id = $1
	`, id).Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.ID, &car.Owner.Name, &car.Owner.Surname)
	if err != nil {
		return models.CarResponse{}, fmt.Errorf("car with id %s not found", id)
	}

	return car, nil
}

func (r *carRepository) GetCarsByOwner(ctx context.Context, ownerID string) ([]models.CarResponse, error) {
	var cars []models.CarResponse
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.reg_number, c.mark, c.model, c.year, p.id, p.name, p.surname 
		FROM cars c
		INNER JOIN person p ON c.owner = p.id
		WHERE p.id = $1
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.CarResponse
		err := rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.ID, &car.Owner.Name, &car.Owner.Surname)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *carRepository) UpdateCar(ctx context.Context, car models.CarResponse) error {
	var query strings.Builder
	var args []interface{}

	query.WriteString("UPDATE cars SET ")

	i := 1
	for k, v := range car.ValueToUpdate() {
		if i > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}

	query.WriteString(fmt.Sprintf(" WHERE id = $%d; ", i))
	args = append(args, car.ID)

	if car.Owner.Name != "" && car.Owner.Surname != "" {
		query.WriteString("UPDATE person SET name = $1, surname = $2 WHERE id = $3")
	}

	result, err := r.db.Exec(ctx, query.String(), args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("car with id %s not found", car.ID)
	}

	return nil
}

func (r *carRepository) DeleteCar(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("car with id %s not found", id)
	}
	return nil
}
