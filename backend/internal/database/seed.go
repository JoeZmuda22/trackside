package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func Seed(db *sql.DB) error {
	// Check if demo user already exists
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM "User" WHERE email = ?`, "demo@trackside.com").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database...")
	now := time.Now().UTC().Format(time.RFC3339)

	// Create demo user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	if err != nil {
		return err
	}

	userID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "User" (id, name, email, passwordHash, experience, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, "Demo Driver", "demo@trackside.com", string(passwordHash), "INTERMEDIATE", now, now)
	if err != nil {
		return err
	}
	log.Println("Created demo user: demo@trackside.com")

	// Create car
	carID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "Car" (id, make, model, year, userId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		carID, "Nissan", "350Z", 2006, userID, now, now)
	if err != nil {
		return err
	}

	// Create car mods
	mods := []struct{ name, category string }{
		{"BC Racing BR Coilovers", "SUSPENSION"},
		{"Tomei Expreme Ti Exhaust", "EXHAUST"},
		{"Z1 Motorsports Cold Air Intake", "ENGINE"},
		{"Stoptech ST-40 Big Brake Kit", "BRAKES"},
		{"Enkei RPF1 18x9.5", "WHEELS_TIRES"},
	}
	for _, m := range mods {
		_, err = db.Exec(`INSERT INTO "CarMod" (id, name, category, carId) VALUES (?, ?, ?, ?)`,
			xid.New().String(), m.name, m.category, carID)
		if err != nil {
			return err
		}
	}
	log.Println("Created demo car: 2006 Nissan 350Z with 5 mods")

	// Create track 1: Laguna Seca
	track1ID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "Track" (id, name, location, state, latitude, longitude, description, uploadedById, status, isImported, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		track1ID, "Laguna Seca", "Monterey, CA", "CA", 36.5754, -121.7627,
		"Iconic road course featuring the famous Corkscrew turn. 2.238 miles, 11 turns, and significant elevation changes make this a challenging and rewarding track.",
		userID, "APPROVED", false, now, now)
	if err != nil {
		return err
	}

	// Track 1 events
	roadcourseEvent1ID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`, roadcourseEvent1ID, "ROADCOURSE", track1ID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`, xid.New().String(), "DRIFT", track1ID)
	if err != nil {
		return err
	}

	// Track 1 zones
	zones := []struct {
		name, desc string
		posX, posY float64
	}{
		{"The Corkscrew (T8-T8A)", "Famous downhill left-right combo with 5.5 stories of elevation change. Blind entry — use the tree as a braking marker.", 65, 25},
		{"Turn 2 (Andretti Hairpin)", "Tight left-hand hairpin. Late apex is key.", 30, 40},
		{"Turn 5", "High-speed left sweeper heading uphill. Carry momentum.", 45, 60},
	}

	zoneIDs := make([]string, len(zones))
	for i, z := range zones {
		zoneIDs[i] = xid.New().String()
		_, err = db.Exec(`INSERT INTO "TrackZone" (id, name, description, posX, posY, trackId, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			zoneIDs[i], z.name, z.desc, z.posX, z.posY, track1ID, now)
		if err != nil {
			return err
		}
	}
	log.Println("Created track: Laguna Seca")

	// Create track 2: Atlanta Motorsports Park
	track2ID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "Track" (id, name, location, state, latitude, longitude, description, uploadedById, status, isImported, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		track2ID, "Atlanta Motorsports Park", "Dawsonville, GA", "GA", 34.3705, -84.1643,
		"A 2-mile, 16-turn road course with 100ft of elevation change, designed by Hermann Tilke. Features a dedicated drift pad and drag strip.",
		userID, "APPROVED", false, now, now)
	if err != nil {
		return err
	}

	// Track 2 events
	roadcourseEvent2ID := xid.New().String()
	_, err = db.Exec(`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`, roadcourseEvent2ID, "ROADCOURSE", track2ID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`, xid.New().String(), "DRIFT", track2ID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO "TrackEvent" (id, eventType, trackId) VALUES (?, ?, ?)`, xid.New().String(), "DRAG", track2ID)
	if err != nil {
		return err
	}

	// Track 2 zones
	track2Zones := []struct {
		name, desc string
		posX, posY float64
	}{
		{"Turn 1", "Fast right-hander after the main straight. Heavy braking zone.", 80, 30},
		{"Turn 12 (Rollercoaster)", "Blind crest into a left-right combo. Commitment corner.", 35, 55},
	}
	for _, z := range track2Zones {
		_, err = db.Exec(`INSERT INTO "TrackZone" (id, name, description, posX, posY, trackId, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			xid.New().String(), z.name, z.desc, z.posX, z.posY, track2ID, now)
		if err != nil {
			return err
		}
	}
	log.Println("Created track: Atlanta Motorsports Park")

	// Zone tips
	tips := []struct {
		content string
		zoneIdx int
	}{
		{"Use the big tree on the left as your turn-in point. Trust the line and commit — hesitation here is dangerous.", 0},
		{"Brake deep and trail brake in. The car will rotate naturally. Get on power early for the uphill section.", 1},
		{"Stay wide and use all the road. The banking helps you carry more speed than you think.", 2},
	}
	for _, t := range tips {
		_, err = db.Exec(`INSERT INTO "ZoneTip" (id, content, conditions, zoneId, authorId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			xid.New().String(), t.content, "DRY", zoneIDs[t.zoneIdx], userID, now, now)
		if err != nil {
			return err
		}
	}
	log.Println("Created zone tips")

	// Review
	_, err = db.Exec(`INSERT INTO "TrackReview" (id, rating, content, conditions, trackId, trackEventId, authorId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		xid.New().String(), 5,
		"Absolutely world-class track. The Corkscrew is every bit as intense as it looks on TV. The facility is well-maintained and the tech inspection is thorough. Will definitely be back.",
		"DRY", track1ID, roadcourseEvent1ID, userID, now, now)
	if err != nil {
		return err
	}
	log.Println("Created reviews")

	// Lap records
	_, err = db.Exec(`INSERT INTO "LapRecord" (id, lapTime, conditions, notes, tirePressureFL, tirePressureFR, tirePressureRL, tirePressureRR, fuelLevel, camberFL, camberFR, camberRL, camberRR, casterFL, casterFR, toeFL, toeFR, toeRL, toeRR, trackId, trackEventId, carId, driverId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		xid.New().String(), "1:42.856", "DRY",
		"Best time of the day. Car felt great after adjusting front camber.",
		32.5, 32.5, 34.0, 34.0, 50.0, -2.5, -2.5, -1.8, -1.8, 5.2, 5.2, 0.1, 0.1, 0.15, 0.15,
		track2ID, roadcourseEvent2ID, carID, userID, now, now)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO "LapRecord" (id, lapTime, conditions, notes, tirePressureFL, tirePressureFR, tirePressureRL, tirePressureRR, fuelLevel, camberFL, camberFR, camberRL, camberRR, trackId, trackEventId, carId, driverId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		xid.New().String(), "1:45.112", "WET",
		"Started raining mid-session. Dropped tire pressure to help with wet grip.",
		30.0, 30.0, 32.0, 32.0, 40.0, -2.5, -2.5, -1.8, -1.8,
		track2ID, roadcourseEvent2ID, carID, userID, now, now)
	if err != nil {
		return err
	}
	log.Println("Created lap records")

	log.Println("")
	log.Println("Seed complete! Login with:")
	log.Println("  Email: demo@trackside.com")
	log.Println("  Password: password123")
	return nil
}
