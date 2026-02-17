package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
)

type LapbookRepo struct {
	db *sql.DB
}

func NewLapbookRepo(db *sql.DB) *LapbookRepo {
	return &LapbookRepo{db: db}
}

func (r *LapbookRepo) List(driverID, trackID, eventType, carID string) ([]models.LapRecordWithDetails, error) {
	query := `SELECT lr.id, lr.lapTime, lr.conditions, lr.notes,
		lr.tirePressureFL, lr.tirePressureFR, lr.tirePressureRL, lr.tirePressureRR,
		lr.fuelLevel, lr.camberFL, lr.camberFR, lr.camberRL, lr.camberRR,
		lr.casterFL, lr.casterFR, lr.toeFL, lr.toeFR, lr.toeRL, lr.toeRR,
		lr.trackId, lr.trackEventId, lr.carId, lr.driverId, lr.createdAt, lr.updatedAt,
		t.id, t.name, t.location,
		c.id, c.make, c.model, c.year
		FROM "LapRecord" lr
		JOIN "Track" t ON lr.trackId = t.id
		JOIN "Car" c ON lr.carId = c.id
		WHERE lr.driverId = ?`

	args := []interface{}{driverID}
	conditions := []string{}

	if trackID != "" {
		conditions = append(conditions, `lr.trackId = ?`)
		args = append(args, trackID)
	}
	if carID != "" {
		conditions = append(conditions, `lr.carId = ?`)
		args = append(args, carID)
	}
	if eventType != "" {
		conditions = append(conditions, `EXISTS (SELECT 1 FROM "TrackEvent" te WHERE te.id = lr.trackEventId AND te.eventType = ?)`)
		args = append(args, eventType)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY lr.createdAt DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("lapbook list query: %w", err)
	}
	defer rows.Close()

	var records []models.LapRecordWithDetails
	for rows.Next() {
		var lr models.LapRecordWithDetails
		if err := rows.Scan(
			&lr.ID, &lr.LapTime, &lr.Conditions, &lr.Notes,
			&lr.TirePressureFL, &lr.TirePressureFR, &lr.TirePressureRL, &lr.TirePressureRR,
			&lr.FuelLevel, &lr.CamberFL, &lr.CamberFR, &lr.CamberRL, &lr.CamberRR,
			&lr.CasterFL, &lr.CasterFR, &lr.ToeFL, &lr.ToeFR, &lr.ToeRL, &lr.ToeRR,
			&lr.TrackID, &lr.TrackEventID, &lr.CarID, &lr.DriverID, &lr.CreatedAt, &lr.UpdatedAt,
			&lr.Track.ID, &lr.Track.Name, &lr.Track.Location,
			&lr.Car.ID, &lr.Car.Make, &lr.Car.Model, &lr.Car.Year,
		); err != nil {
			return nil, err
		}
		records = append(records, lr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	// Fetch track events AFTER closing rows to avoid SQLite deadlock
	for i := range records {
		if records[i].TrackEventID != nil {
			te := &models.TrackEvent{}
			err := r.db.QueryRow(
				`SELECT id, eventType, trackId FROM "TrackEvent" WHERE id = ?`, *records[i].TrackEventID,
			).Scan(&te.ID, &te.EventType, &te.TrackID)
			if err == nil {
				records[i].TrackEvent = te
			}
		}
	}

	if records == nil {
		records = []models.LapRecordWithDetails{}
	}
	return records, nil
}

func (r *LapbookRepo) Create(req models.LapRecordRequest, driverID string) (*models.LapRecordWithDetails, error) {
	now := time.Now().UTC()
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "LapRecord" (id, lapTime, conditions, notes,
			tirePressureFL, tirePressureFR, tirePressureRL, tirePressureRR,
			fuelLevel, camberFL, camberFR, camberRL, camberRR,
			casterFL, casterFR, toeFL, toeFR, toeRL, toeRR,
			trackId, trackEventId, carId, driverId, createdAt, updatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, req.LapTime, req.Conditions, req.Notes,
		req.TirePressureFL, req.TirePressureFR, req.TirePressureRL, req.TirePressureRR,
		req.FuelLevel, req.CamberFL, req.CamberFR, req.CamberRL, req.CamberRR,
		req.CasterFL, req.CasterFR, req.ToeFL, req.ToeFR, req.ToeRL, req.ToeRR,
		req.TrackID, req.TrackEventID, req.CarID, driverID, now, now,
	)
	if err != nil {
		return nil, err
	}

	lr := &models.LapRecordWithDetails{
		ID: id, LapTime: req.LapTime, Conditions: req.Conditions, Notes: req.Notes,
		TirePressureFL: req.TirePressureFL, TirePressureFR: req.TirePressureFR,
		TirePressureRL: req.TirePressureRL, TirePressureRR: req.TirePressureRR,
		FuelLevel: req.FuelLevel,
		CamberFL: req.CamberFL, CamberFR: req.CamberFR, CamberRL: req.CamberRL, CamberRR: req.CamberRR,
		CasterFL: req.CasterFL, CasterFR: req.CasterFR,
		ToeFL: req.ToeFL, ToeFR: req.ToeFR, ToeRL: req.ToeRL, ToeRR: req.ToeRR,
		TrackID: req.TrackID, TrackEventID: req.TrackEventID, CarID: req.CarID,
		DriverID: driverID, CreatedAt: now, UpdatedAt: now,
	}

	r.db.QueryRow(`SELECT id, name, location FROM "Track" WHERE id = ?`, req.TrackID).Scan(
		&lr.Track.ID, &lr.Track.Name, &lr.Track.Location,
	)
	r.db.QueryRow(`SELECT id, make, model, year FROM "Car" WHERE id = ?`, req.CarID).Scan(
		&lr.Car.ID, &lr.Car.Make, &lr.Car.Model, &lr.Car.Year,
	)

	if req.TrackEventID != nil {
		te := &models.TrackEvent{}
		err := r.db.QueryRow(
			`SELECT id, eventType, trackId FROM "TrackEvent" WHERE id = ?`, *req.TrackEventID,
		).Scan(&te.ID, &te.EventType, &te.TrackID)
		if err == nil {
			lr.TrackEvent = te
		}
	}

	return lr, nil
}

func (r *LapbookRepo) FindByIDAndDriver(id, driverID string) (*models.LapRecord, error) {
	lr := &models.LapRecord{}
	err := r.db.QueryRow(
		`SELECT id, driverId FROM "LapRecord" WHERE id = ? AND driverId = ?`, id, driverID,
	).Scan(&lr.ID, &lr.DriverID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return lr, nil
}

func (r *LapbookRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM "LapRecord" WHERE id = ?`, id)
	return err
}
