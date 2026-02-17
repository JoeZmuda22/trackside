package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
)

type TrackRepo struct {
	db *sql.DB
}

func NewTrackRepo(db *sql.DB) *TrackRepo {
	return &TrackRepo{db: db}
}

func (r *TrackRepo) List(search, eventType, state string) ([]models.TrackListItem, error) {
	query := `SELECT t.id, t.name, t.location, t.state, t.description, t.imageUrl, t.latitude, t.longitude,
		t.status, t.isImported, t.uploadedById, t.createdAt, t.updatedAt,
		u.id, u.name
		FROM "Track" t
		JOIN "User" u ON t.uploadedById = u.id
		WHERE t.status = 'APPROVED'`
	args := []interface{}{}

	if search != "" {
		query += ` AND (t.name LIKE ? OR t.location LIKE ?)`
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	if eventType != "" {
		query += ` AND EXISTS (SELECT 1 FROM "TrackEvent" te WHERE te.trackId = t.id AND te.eventType = ?)`
		args = append(args, eventType)
	}

	if state != "" {
		query += ` AND t.state = ?`
		args = append(args, strings.ToUpper(state))
	}

	query += ` ORDER BY t.createdAt DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.TrackListItem
	for rows.Next() {
		var t models.TrackListItem
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Location, &t.State, &t.Description, &t.ImageURL,
			&t.Latitude, &t.Longitude, &t.Status, &t.IsImported, &t.UploadedByID,
			&t.CreatedAt, &t.UpdatedAt,
			&t.UploadedBy.ID, &t.UploadedBy.Name,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	// Fetch nested data AFTER closing the rows cursor to avoid SQLite deadlock
	for i := range tracks {
		events, err := r.GetEvents(tracks[i].ID)
		if err != nil {
			return nil, err
		}
		tracks[i].Events = events

		r.db.QueryRow(`SELECT COUNT(*) FROM "TrackReview" WHERE trackId = ?`, tracks[i].ID).Scan(&tracks[i].Count.Reviews)
		r.db.QueryRow(`SELECT COUNT(*) FROM "TrackZone" WHERE trackId = ?`, tracks[i].ID).Scan(&tracks[i].Count.Zones)
		r.db.QueryRow(`SELECT COUNT(*) FROM "LapRecord" WHERE trackId = ?`, tracks[i].ID).Scan(&tracks[i].Count.LapRecords)
	}

	// Get avg ratings in bulk
	ratingRows, err := r.db.Query(`SELECT trackId, AVG(rating) FROM "TrackReview" GROUP BY trackId`)
	if err != nil {
		return nil, err
	}
	defer ratingRows.Close()

	ratingMap := make(map[string]float64)
	for ratingRows.Next() {
		var trackID string
		var avg float64
		ratingRows.Scan(&trackID, &avg)
		ratingMap[trackID] = avg
	}

	for i := range tracks {
		tracks[i].AvgRating = ratingMap[tracks[i].ID]
	}

	if tracks == nil {
		tracks = []models.TrackListItem{}
	}
	return tracks, nil
}

func (r *TrackRepo) FindByID(id string) (*models.Track, error) {
	t := &models.Track{}
	err := r.db.QueryRow(
		`SELECT id, name, location, state, description, imageUrl, latitude, longitude, status, isImported, uploadedById, createdAt, updatedAt FROM "Track" WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Name, &t.Location, &t.State, &t.Description, &t.ImageURL,
		&t.Latitude, &t.Longitude, &t.Status, &t.IsImported, &t.UploadedByID, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TrackRepo) GetDetail(id string, eventTypeFilter string) (*models.TrackDetail, error) {
	t, err := r.FindByID(id)
	if err != nil || t == nil {
		return nil, err
	}

	detail := &models.TrackDetail{
		ID: t.ID, Name: t.Name, Location: t.Location, State: t.State,
		Description: t.Description, ImageURL: t.ImageURL, Latitude: t.Latitude,
		Longitude: t.Longitude, Status: t.Status, IsImported: t.IsImported,
		UploadedByID: t.UploadedByID, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}

	// Uploader with experience
	r.db.QueryRow(
		`SELECT id, name, experience FROM "User" WHERE id = ?`, t.UploadedByID,
	).Scan(&detail.UploadedBy.ID, &detail.UploadedBy.Name, &detail.UploadedBy.Experience)

	// Events
	events, err := r.GetEvents(id)
	if err != nil {
		return nil, err
	}
	detail.Events = events

	// Zones (optionally filtered by eventType)
	zoneRepo := NewZoneRepo(r.db)
	zones, err := zoneRepo.GetZonesWithTips(id, eventTypeFilter)
	if err != nil {
		return nil, err
	}
	detail.Zones = zones

	// Reviews with author details
	reviewRepo := NewReviewRepo(r.db)
	reviews, err := reviewRepo.GetTrackReviews(id)
	if err != nil {
		return nil, err
	}
	detail.Reviews = reviews

	// Counts
	r.db.QueryRow(`SELECT COUNT(*) FROM "TrackReview" WHERE trackId = ?`, id).Scan(&detail.Count.Reviews)
	r.db.QueryRow(`SELECT COUNT(*) FROM "TrackZone" WHERE trackId = ?`, id).Scan(&detail.Count.Zones)
	r.db.QueryRow(`SELECT COUNT(*) FROM "LapRecord" WHERE trackId = ?`, id).Scan(&detail.Count.LapRecords)

	// Average rating
	var avg sql.NullFloat64
	r.db.QueryRow(`SELECT AVG(rating) FROM "TrackReview" WHERE trackId = ?`, id).Scan(&avg)
	if avg.Valid {
		detail.AvgRating = avg.Float64
	}

	return detail, nil
}

func (r *TrackRepo) Create(req models.TrackRequest, userID string) (*models.TrackListItem, error) {
	now := time.Now().UTC()
	id := xid.New().String()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO "Track" (id, name, location, description, imageUrl, uploadedById, status, isImported, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, 'APPROVED', false, ?, ?)`,
		id, req.Name, req.Location, req.Description, req.ImageURL, userID, now, now,
	)
	if err != nil {
		return nil, err
	}

	for _, et := range req.EventTypes {
		_, err = tx.Exec(
			`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`,
			xid.New().String(), et, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return created track with events and uploader
	result := &models.TrackListItem{
		ID: id, Name: req.Name, Location: req.Location,
		Description: req.Description, ImageURL: req.ImageURL,
		Status: "APPROVED", IsImported: false, UploadedByID: userID,
		CreatedAt: now, UpdatedAt: now,
	}

	r.db.QueryRow(`SELECT id, name FROM "User" WHERE id = ?`, userID).Scan(&result.UploadedBy.ID, &result.UploadedBy.Name)

	events, _ := r.GetEvents(id)
	result.Events = events
	result.Count = models.TrackCounts{}
	return result, nil
}

func (r *TrackRepo) Update(id string, req models.TrackPatchRequest) (*models.Track, error) {
	now := time.Now().UTC()
	sets := []string{"updatedAt = ?"}
	args := []interface{}{now}

	if req.ImageURL != nil {
		sets = append(sets, `imageUrl = ?`)
		args = append(args, *req.ImageURL)
	}
	if req.Name != nil {
		sets = append(sets, `name = ?`)
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		sets = append(sets, `description = ?`)
		args = append(args, *req.Description)
	}
	if req.Location != nil {
		sets = append(sets, `location = ?`)
		args = append(args, *req.Location)
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE "Track" SET %s WHERE id = ?`, strings.Join(sets, ", "))
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *TrackRepo) GetEvents(trackID string) ([]models.TrackEvent, error) {
	rows, err := r.db.Query(`SELECT id, eventType, trackId FROM "TrackEvent" WHERE trackId = ?`, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.TrackEvent
	for rows.Next() {
		var e models.TrackEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.TrackID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if events == nil {
		events = []models.TrackEvent{}
	}
	return events, nil
}

func (r *TrackRepo) FindByNameAndLocation(name, location string) (*models.Track, error) {
	t := &models.Track{}
	err := r.db.QueryRow(
		`SELECT id, name, location, state, description, imageUrl, latitude, longitude, status, isImported, uploadedById, createdAt, updatedAt FROM "Track" WHERE name = ? AND location = ?`,
		name, location,
	).Scan(&t.ID, &t.Name, &t.Location, &t.State, &t.Description, &t.ImageURL,
		&t.Latitude, &t.Longitude, &t.Status, &t.IsImported, &t.UploadedByID, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TrackRepo) UpsertImported(data models.ImportedTrack, systemUserID string) error {
	existing, err := r.FindByNameAndLocation(data.Name, data.Location)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	if existing != nil {
		_, err = r.db.Exec(
			`UPDATE "Track" SET description = ?, latitude = ?, longitude = ?, state = ?, isImported = true, updatedAt = ? WHERE id = ?`,
			data.Description, data.Latitude, data.Longitude, data.State, now, existing.ID,
		)
		return err
	}

	id := xid.New().String()
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO "Track" (id, name, location, state, description, latitude, longitude, isImported, status, uploadedById, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, true, 'APPROVED', ?, ?, ?)`,
		id, data.Name, data.Location, data.State, data.Description, data.Latitude, data.Longitude, systemUserID, now, now,
	)
	if err != nil {
		return err
	}

	for _, t := range data.Types {
		_, err = tx.Exec(
			`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`,
			xid.New().String(), strings.ToUpper(t), id,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ─── Track Images ───────────────────────────────────────────────────────────────

func (r *TrackRepo) GetImages(trackID string) ([]models.TrackImageWithUploader, error) {
	rows, err := r.db.Query(
		`SELECT ti.id, ti.url, ti.caption, ti.createdAt, u.id, u.name
		FROM "TrackImage" ti
		JOIN "User" u ON ti.uploadedById = u.id
		WHERE ti.trackId = ?
		ORDER BY ti.createdAt DESC`,
		trackID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.TrackImageWithUploader
	for rows.Next() {
		var img models.TrackImageWithUploader
		if err := rows.Scan(&img.ID, &img.URL, &img.Caption, &img.CreatedAt, &img.UploadedBy.ID, &img.UploadedBy.Name); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	if images == nil {
		images = []models.TrackImageWithUploader{}
	}
	return images, nil
}

func (r *TrackRepo) CreateImage(url string, caption *string, trackID, userID string) (*models.TrackImageWithUploader, error) {
	id := xid.New().String()
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`INSERT INTO "TrackImage" (id, url, caption, trackId, uploadedById, createdAt) VALUES (?, ?, ?, ?, ?, ?)`,
		id, url, caption, trackID, userID, now,
	)
	if err != nil {
		return nil, err
	}

	img := &models.TrackImageWithUploader{
		ID: id, URL: url, Caption: caption, CreatedAt: now,
	}
	r.db.QueryRow(`SELECT id, name FROM "User" WHERE id = ?`, userID).Scan(&img.UploadedBy.ID, &img.UploadedBy.Name)
	return img, nil
}

func (r *TrackRepo) FindImage(imageID string) (*models.TrackImage, error) {
	img := &models.TrackImage{}
	err := r.db.QueryRow(
		`SELECT id, url, caption, trackId, uploadedById, createdAt FROM "TrackImage" WHERE id = ?`, imageID,
	).Scan(&img.ID, &img.URL, &img.Caption, &img.TrackID, &img.UploadedByID, &img.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (r *TrackRepo) DeleteImage(imageID string) error {
	_, err := r.db.Exec(`DELETE FROM "TrackImage" WHERE id = ?`, imageID)
	return err
}

func (r *TrackRepo) Exists(id string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM "Track" WHERE id = ?`, id).Scan(&count)
	return count > 0, err
}
