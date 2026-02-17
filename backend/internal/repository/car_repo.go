package repository

import (
	"database/sql"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
)

type CarRepo struct {
	db *sql.DB
}

func NewCarRepo(db *sql.DB) *CarRepo {
	return &CarRepo{db: db}
}

func (r *CarRepo) FindByUserID(userID string) ([]models.Car, error) {
	rows, err := r.db.Query(
		`SELECT id, make, model, year, userId, createdAt, updatedAt FROM "Car" WHERE userId = ? ORDER BY createdAt DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var c models.Car
		if err := rows.Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.UserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Fetch mods after closing the rows cursor to avoid nested query deadlock
	for i := range cars {
		mods, err := r.GetMods(cars[i].ID)
		if err != nil {
			return nil, err
		}
		cars[i].Mods = mods
	}

	if cars == nil {
		cars = []models.Car{}
	}
	return cars, nil
}

func (r *CarRepo) FindByIDAndUser(id, userID string) (*models.Car, error) {
	c := &models.Car{}
	err := r.db.QueryRow(
		`SELECT id, make, model, year, userId, createdAt, updatedAt FROM "Car" WHERE id = ? AND userId = ?`,
		id, userID,
	).Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.UserID, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	mods, err := r.GetMods(c.ID)
	if err != nil {
		return nil, err
	}
	c.Mods = mods
	return c, nil
}

func (r *CarRepo) Create(make, model string, year int, userID string) (*models.Car, error) {
	now := time.Now().UTC()
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "Car" (id, make, model, year, userId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, make, model, year, userID, now, now,
	)
	if err != nil {
		return nil, err
	}
	return &models.Car{
		ID:        id,
		Make:      make,
		Model:     model,
		Year:      year,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		Mods:      []models.CarMod{},
	}, nil
}

func (r *CarRepo) Update(id, make, model string, year int) (*models.Car, error) {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE "Car" SET make = ?, model = ?, year = ?, updatedAt = ? WHERE id = ?`,
		make, model, year, now, id,
	)
	if err != nil {
		return nil, err
	}
	c := &models.Car{}
	err = r.db.QueryRow(
		`SELECT id, make, model, year, userId, createdAt, updatedAt FROM "Car" WHERE id = ?`, id,
	).Scan(&c.ID, &c.Make, &c.Model, &c.Year, &c.UserID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mods, err := r.GetMods(id)
	if err != nil {
		return nil, err
	}
	c.Mods = mods
	return c, nil
}

func (r *CarRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM "Car" WHERE id = ?`, id)
	return err
}

func (r *CarRepo) GetMods(carID string) ([]models.CarMod, error) {
	rows, err := r.db.Query(
		`SELECT id, name, category, notes, carId FROM "CarMod" WHERE carId = ?`, carID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mods []models.CarMod
	for rows.Next() {
		var m models.CarMod
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.Notes, &m.CarID); err != nil {
			return nil, err
		}
		mods = append(mods, m)
	}
	if mods == nil {
		mods = []models.CarMod{}
	}
	return mods, nil
}

func (r *CarRepo) CreateMod(name, category string, notes *string, carID string) (*models.CarMod, error) {
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "CarMod" (id, name, category, notes, carId) VALUES (?, ?, ?, ?, ?)`,
		id, name, category, notes, carID,
	)
	if err != nil {
		return nil, err
	}
	return &models.CarMod{
		ID:       id,
		Name:     name,
		Category: category,
		Notes:    notes,
		CarID:    carID,
	}, nil
}

func (r *CarRepo) FindMod(modID, carID string) (*models.CarMod, error) {
	m := &models.CarMod{}
	err := r.db.QueryRow(
		`SELECT id, name, category, notes, carId FROM "CarMod" WHERE id = ? AND carId = ?`,
		modID, carID,
	).Scan(&m.ID, &m.Name, &m.Category, &m.Notes, &m.CarID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *CarRepo) DeleteMod(modID string) error {
	_, err := r.db.Exec(`DELETE FROM "CarMod" WHERE id = ?`, modID)
	return err
}

func (r *CarRepo) ExistsForUser(carID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM "Car" WHERE id = ? AND userId = ?`, carID, userID).Scan(&count)
	return count > 0, err
}
