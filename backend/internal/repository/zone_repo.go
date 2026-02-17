package repository

import (
	"database/sql"
	"time"

	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/rs/xid"
)

type ZoneRepo struct {
	db *sql.DB
}

func NewZoneRepo(db *sql.DB) *ZoneRepo {
	return &ZoneRepo{db: db}
}

func (r *ZoneRepo) GetZonesWithTips(trackID string, eventTypeFilter string) ([]models.TrackZoneWithTips, error) {
	query := `SELECT id, name, description, posX, posY, trackId, eventType, createdAt FROM "TrackZone" WHERE trackId = ?`
	args := []interface{}{trackID}

	if eventTypeFilter != "" {
		query += ` AND eventType = ?`
		args = append(args, eventTypeFilter)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []models.TrackZoneWithTips
	for rows.Next() {
		var z models.TrackZoneWithTips
		if err := rows.Scan(&z.ID, &z.Name, &z.Description, &z.PosX, &z.PosY, &z.TrackID, &z.EventType, &z.CreatedAt); err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	// Fetch tips AFTER closing rows to avoid SQLite deadlock
	for i := range zones {
		tips, err := r.GetTipsForZone(zones[i].ID)
		if err != nil {
			return nil, err
		}
		zones[i].Tips = tips
	}

	if zones == nil {
		zones = []models.TrackZoneWithTips{}
	}
	return zones, nil
}

func (r *ZoneRepo) GetTipsForZone(zoneID string) ([]models.ZoneTipWithAuthor, error) {
	rows, err := r.db.Query(
		`SELECT zt.id, zt.content, zt.conditions, zt.zoneId, zt.authorId, zt.createdAt, zt.updatedAt,
			u.id, u.name
		FROM "ZoneTip" zt
		JOIN "User" u ON zt.authorId = u.id
		WHERE zt.zoneId = ?
		ORDER BY zt.createdAt DESC`,
		zoneID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tips []models.ZoneTipWithAuthor
	for rows.Next() {
		var t models.ZoneTipWithAuthor
		if err := rows.Scan(&t.ID, &t.Content, &t.Conditions, &t.ZoneID, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt, &t.Author.ID, &t.Author.Name); err != nil {
			return nil, err
		}
		tips = append(tips, t)
	}
	if tips == nil {
		tips = []models.ZoneTipWithAuthor{}
	}
	return tips, nil
}

func (r *ZoneRepo) Create(name string, description *string, posX, posY float64, trackID string, eventType *string) (*models.TrackZoneWithTips, error) {
	now := time.Now().UTC()
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "TrackZone" (id, name, description, posX, posY, trackId, eventType, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, name, description, posX, posY, trackID, eventType, now,
	)
	if err != nil {
		return nil, err
	}
	return &models.TrackZoneWithTips{
		ID: id, Name: name, Description: description, PosX: posX, PosY: posY,
		TrackID: trackID, EventType: eventType, CreatedAt: now,
		Tips: []models.ZoneTipWithAuthor{},
	}, nil
}

func (r *ZoneRepo) FindByID(zoneID string) (*models.TrackZone, error) {
	z := &models.TrackZone{}
	err := r.db.QueryRow(
		`SELECT id, name, description, posX, posY, trackId, eventType, createdAt FROM "TrackZone" WHERE id = ?`,
		zoneID,
	).Scan(&z.ID, &z.Name, &z.Description, &z.PosX, &z.PosY, &z.TrackID, &z.EventType, &z.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return z, nil
}

func (r *ZoneRepo) Update(zoneID string, name *string, description *string) (*models.TrackZoneWithTips, error) {
	if name != nil {
		r.db.Exec(`UPDATE "TrackZone" SET name = ? WHERE id = ?`, *name, zoneID)
	}
	if description != nil {
		r.db.Exec(`UPDATE "TrackZone" SET description = ? WHERE id = ?`, *description, zoneID)
	}

	z, err := r.FindByID(zoneID)
	if err != nil || z == nil {
		return nil, err
	}

	tips, err := r.GetTipsForZone(zoneID)
	if err != nil {
		return nil, err
	}

	return &models.TrackZoneWithTips{
		ID: z.ID, Name: z.Name, Description: z.Description, PosX: z.PosX, PosY: z.PosY,
		TrackID: z.TrackID, EventType: z.EventType, CreatedAt: z.CreatedAt,
		Tips: tips,
	}, nil
}

func (r *ZoneRepo) Delete(zoneID string) error {
	_, err := r.db.Exec(`DELETE FROM "TrackZone" WHERE id = ?`, zoneID)
	return err
}

func (r *ZoneRepo) CreateTip(content string, conditions *string, zoneID, authorID string) (*models.ZoneTipWithAuthor, error) {
	now := time.Now().UTC()
	id := xid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "ZoneTip" (id, content, conditions, zoneId, authorId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, content, conditions, zoneID, authorID, now, now,
	)
	if err != nil {
		return nil, err
	}

	tip := &models.ZoneTipWithAuthor{
		ID: id, Content: content, Conditions: conditions, ZoneID: zoneID,
		AuthorID: authorID, CreatedAt: now, UpdatedAt: now,
	}
	r.db.QueryRow(`SELECT id, name FROM "User" WHERE id = ?`, authorID).Scan(&tip.Author.ID, &tip.Author.Name)
	return tip, nil
}

func (r *ZoneRepo) FindZoneForTrack(zoneID, trackID string) (*models.TrackZone, error) {
	z := &models.TrackZone{}
	err := r.db.QueryRow(
		`SELECT id, name, description, posX, posY, trackId, eventType, createdAt FROM "TrackZone" WHERE id = ? AND trackId = ?`,
		zoneID, trackID,
	).Scan(&z.ID, &z.Name, &z.Description, &z.PosX, &z.PosY, &z.TrackID, &z.EventType, &z.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return z, nil
}
