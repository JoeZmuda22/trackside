package repository

import (
	"database/sql"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, name, email, passwordHash, image, experience, createdAt, updatedAt FROM "User" WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Image, &u.Experience, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindByID(id string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, name, email, passwordHash, image, experience, createdAt, updatedAt FROM "User" WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Image, &u.Experience, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) Create(name, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	id := xid.New().String()
	_, err = r.db.Exec(
		`INSERT INTO "User" (id, name, email, passwordHash, experience, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, name, email, string(hash), "BEGINNER", now, now,
	)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:        id,
		Name:      &name,
		Email:     email,
		Experience: "BEGINNER",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *UserRepo) UpdateProfile(id, name, experience string) (*models.User, error) {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE "User" SET name = ?, experience = ?, updatedAt = ? WHERE id = ?`,
		name, experience, now, id,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *UserRepo) GetProfile(id string) (*models.ProfileResponse, error) {
	u, err := r.FindByID(id)
	if err != nil || u == nil {
		return nil, err
	}

	profile := &models.ProfileResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Experience: u.Experience,
		Image:      u.Image,
		CreatedAt:  u.CreatedAt,
	}

	// Get cars with mods
	carRepo := NewCarRepo(r.db)
	cars, err := carRepo.FindByUserID(id)
	if err != nil {
		return nil, err
	}
	profile.Cars = cars

	// Get counts
	r.db.QueryRow(`SELECT COUNT(*) FROM "TrackReview" WHERE authorId = ?`, id).Scan(&profile.Count.TrackReviews)
	r.db.QueryRow(`SELECT COUNT(*) FROM "LapRecord" WHERE driverId = ?`, id).Scan(&profile.Count.LapRecords)
	r.db.QueryRow(`SELECT COUNT(*) FROM "Track" WHERE uploadedById = ?`, id).Scan(&profile.Count.Tracks)
	r.db.QueryRow(`SELECT COUNT(*) FROM "ZoneTip" WHERE authorId = ?`, id).Scan(&profile.Count.ZoneTips)

	return profile, nil
}

func (r *UserRepo) FindOrCreateSystem() (*models.User, error) {
	u, err := r.FindByEmail("system@trackside.local")
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}

	now := time.Now().UTC()
	id := xid.New().String()
	_, err = r.db.Exec(
		`INSERT INTO "User" (id, name, email, emailVerified, experience, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, "Trackside System", "system@trackside.local", now, "BEGINNER", now, now,
	)
	if err != nil {
		return nil, err
	}
	name := "Trackside System"
	return &models.User{
		ID:        id,
		Name:      &name,
		Email:     "system@trackside.local",
		Experience: "BEGINNER",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
