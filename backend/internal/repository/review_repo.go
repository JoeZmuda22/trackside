package repository

import (
	"database/sql"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
)

type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) GetTrackReviews(trackID string) ([]models.ReviewWithAuthor, error) {
	rows, err := r.db.Query(
		`SELECT tr.id, tr.rating, tr.content, tr.conditions, tr.trackId, tr.trackEventId, tr.authorId, tr.createdAt, tr.updatedAt,
			u.id, u.name, u.experience
		FROM "TrackReview" tr
		JOIN "User" u ON tr.authorId = u.id
		WHERE tr.trackId = ?
		ORDER BY tr.createdAt DESC`,
		trackID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewWithAuthor
	for rows.Next() {
		var rv models.ReviewWithAuthor
		if err := rows.Scan(
			&rv.ID, &rv.Rating, &rv.Content, &rv.Conditions, &rv.TrackID,
			&rv.TrackEventID, &rv.AuthorID, &rv.CreatedAt, &rv.UpdatedAt,
			&rv.Author.ID, &rv.Author.Name, &rv.Author.Experience,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	// Fetch nested data AFTER closing rows to avoid SQLite deadlock
	for i := range reviews {
		// Get author's cars
		carRows, err := r.db.Query(
			`SELECT make, model, year FROM "Car" WHERE userId = ?`, reviews[i].AuthorID,
		)
		if err != nil {
			return nil, err
		}
		var cars []models.CarBrief
		for carRows.Next() {
			var cb models.CarBrief
			carRows.Scan(&cb.Make, &cb.Model, &cb.Year)
			cars = append(cars, cb)
		}
		carRows.Close()
		if cars == nil {
			cars = []models.CarBrief{}
		}
		reviews[i].Author.Cars = cars

		// Get track event if present
		if reviews[i].TrackEventID != nil {
			te := &models.TrackEvent{}
			err := r.db.QueryRow(
				`SELECT id, eventType, trackId FROM "TrackEvent" WHERE id = ?`, *reviews[i].TrackEventID,
			).Scan(&te.ID, &te.EventType, &te.TrackID)
			if err == nil {
				reviews[i].TrackEvent = te
			}
		}
	}

	if reviews == nil {
		reviews = []models.ReviewWithAuthor{}
	}
	return reviews, nil
}

func (r *ReviewRepo) Create(rating int, content *string, conditions, trackID string, trackEventID *string, authorID string) (*models.ReviewWithAuthor, error) {
	now := time.Now().UTC()
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "TrackReview" (id, rating, content, conditions, trackId, trackEventId, authorId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, rating, content, conditions, trackID, trackEventID, authorID, now, now,
	)
	if err != nil {
		return nil, err
	}

	rv := &models.ReviewWithAuthor{
		ID: id, Rating: rating, Content: content, Conditions: conditions,
		TrackID: trackID, TrackEventID: trackEventID, AuthorID: authorID,
		CreatedAt: now, UpdatedAt: now,
	}

	// Author with cars
	r.db.QueryRow(`SELECT id, name, experience FROM "User" WHERE id = ?`, authorID).Scan(
		&rv.Author.ID, &rv.Author.Name, &rv.Author.Experience,
	)

	carRows, _ := r.db.Query(`SELECT make, model, year FROM "Car" WHERE userId = ?`, authorID)
	if carRows != nil {
		defer carRows.Close()
		var cars []models.CarBrief
		for carRows.Next() {
			var cb models.CarBrief
			carRows.Scan(&cb.Make, &cb.Model, &cb.Year)
			cars = append(cars, cb)
		}
		if cars == nil {
			cars = []models.CarBrief{}
		}
		rv.Author.Cars = cars
	} else {
		rv.Author.Cars = []models.CarBrief{}
	}

	// Track event
	if trackEventID != nil {
		te := &models.TrackEvent{}
		err := r.db.QueryRow(
			`SELECT id, eventType, trackId FROM "TrackEvent" WHERE id = ?`, *trackEventID,
		).Scan(&te.ID, &te.EventType, &te.TrackID)
		if err == nil {
			rv.TrackEvent = te
		}
	}

	return rv, nil
}
