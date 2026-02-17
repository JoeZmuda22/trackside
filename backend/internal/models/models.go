package models

import "time"

// ─── Enums ──────────────────────────────────────────────────────────────────────

type ExperienceLevel string

const (
	ExperienceBeginner     ExperienceLevel = "BEGINNER"
	ExperienceIntermediate ExperienceLevel = "INTERMEDIATE"
	ExperienceAdvanced     ExperienceLevel = "ADVANCED"
	ExperiencePro          ExperienceLevel = "PRO"
)

func ValidExperienceLevel(s string) bool {
	switch ExperienceLevel(s) {
	case ExperienceBeginner, ExperienceIntermediate, ExperienceAdvanced, ExperiencePro:
		return true
	}
	return false
}

type ModCategory string

const (
	ModEngine      ModCategory = "ENGINE"
	ModSuspension  ModCategory = "SUSPENSION"
	ModAero        ModCategory = "AERO"
	ModBrakes      ModCategory = "BRAKES"
	ModWheelsTires ModCategory = "WHEELS_TIRES"
	ModDrivetrain  ModCategory = "DRIVETRAIN"
	ModExhaust     ModCategory = "EXHAUST"
	ModInterior    ModCategory = "INTERIOR"
	ModExterior    ModCategory = "EXTERIOR"
	ModElectronics ModCategory = "ELECTRONICS"
	ModOther       ModCategory = "OTHER"
)

func ValidModCategory(s string) bool {
	switch ModCategory(s) {
	case ModEngine, ModSuspension, ModAero, ModBrakes, ModWheelsTires,
		ModDrivetrain, ModExhaust, ModInterior, ModExterior, ModElectronics, ModOther:
		return true
	}
	return false
}

type TrackStatus string

const (
	TrackStatusPending  TrackStatus = "PENDING"
	TrackStatusApproved TrackStatus = "APPROVED"
	TrackStatusRejected TrackStatus = "REJECTED"
)

type EventType string

const (
	EventAutocross  EventType = "AUTOCROSS"
	EventRoadcourse EventType = "ROADCOURSE"
	EventDrift      EventType = "DRIFT"
	EventDrag       EventType = "DRAG"
)

func ValidEventType(s string) bool {
	switch EventType(s) {
	case EventAutocross, EventRoadcourse, EventDrift, EventDrag:
		return true
	}
	return false
}

type DrivingCondition string

const (
	ConditionDry DrivingCondition = "DRY"
	ConditionWet DrivingCondition = "WET"
)

func ValidDrivingCondition(s string) bool {
	switch DrivingCondition(s) {
	case ConditionDry, ConditionWet:
		return true
	}
	return false
}

// ─── Core Models ────────────────────────────────────────────────────────────────

type User struct {
	ID           string    `json:"id"`
	Name         *string   `json:"name"`
	Email        string    `json:"email"`
	PasswordHash *string   `json:"-"`
	Image        *string   `json:"image"`
	Experience   string    `json:"experience"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Car struct {
	ID        string    `json:"id"`
	Make      string    `json:"make"`
	Model     string    `json:"model"`
	Year      int       `json:"year"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Mods      []CarMod  `json:"mods"`
}

type CarMod struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Notes    *string `json:"notes"`
	CarID    string  `json:"carId"`
}

type Track struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Location     string     `json:"location"`
	State        *string    `json:"state"`
	Description  *string    `json:"description"`
	ImageURL     *string    `json:"imageUrl"`
	Latitude     *float64   `json:"latitude"`
	Longitude    *float64   `json:"longitude"`
	Status       string     `json:"status"`
	IsImported   bool       `json:"isImported"`
	UploadedByID string     `json:"uploadedById"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type TrackImage struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Caption      *string   `json:"caption"`
	TrackID      string    `json:"trackId"`
	UploadedByID string    `json:"uploadedById"`
	CreatedAt    time.Time `json:"createdAt"`
}

type TrackEvent struct {
	ID        string `json:"id"`
	EventType string `json:"eventType"`
	TrackID   string `json:"trackId"`
}

type TrackZone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	PosX        float64   `json:"posX"`
	PosY        float64   `json:"posY"`
	TrackID     string    `json:"trackId"`
	EventType   *string   `json:"eventType"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ZoneTip struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	Conditions *string   `json:"conditions"`
	ZoneID     string    `json:"zoneId"`
	AuthorID   string    `json:"authorId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TrackReview struct {
	ID           string    `json:"id"`
	Rating       int       `json:"rating"`
	Content      *string   `json:"content"`
	Conditions   string    `json:"conditions"`
	TrackID      string    `json:"trackId"`
	TrackEventID *string   `json:"trackEventId"`
	AuthorID     string    `json:"authorId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LapRecord struct {
	ID              string    `json:"id"`
	LapTime         string    `json:"lapTime"`
	Conditions      string    `json:"conditions"`
	Notes           *string   `json:"notes"`
	TirePressureFL  *float64  `json:"tirePressureFL"`
	TirePressureFR  *float64  `json:"tirePressureFR"`
	TirePressureRL  *float64  `json:"tirePressureRL"`
	TirePressureRR  *float64  `json:"tirePressureRR"`
	FuelLevel       *float64  `json:"fuelLevel"`
	CamberFL        *float64  `json:"camberFL"`
	CamberFR        *float64  `json:"camberFR"`
	CamberRL        *float64  `json:"camberRL"`
	CamberRR        *float64  `json:"camberRR"`
	CasterFL        *float64  `json:"casterFL"`
	CasterFR        *float64  `json:"casterFR"`
	ToeFL           *float64  `json:"toeFL"`
	ToeFR           *float64  `json:"toeFR"`
	ToeRL           *float64  `json:"toeRL"`
	ToeRR           *float64  `json:"toeRR"`
	TrackID         string    `json:"trackId"`
	TrackEventID    *string   `json:"trackEventId"`
	CarID           string    `json:"carId"`
	DriverID        string    `json:"driverId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ─── API Response DTOs ──────────────────────────────────────────────────────────

type UserBrief struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type UserWithExperience struct {
	ID         string  `json:"id"`
	Name       *string `json:"name"`
	Experience string  `json:"experience"`
}

type UserWithCars struct {
	ID         string    `json:"id"`
	Name       *string   `json:"name"`
	Experience string    `json:"experience"`
	Cars       []CarBrief `json:"cars"`
}

type CarBrief struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type CarWithID struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type TrackBrief struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type TrackCounts struct {
	Reviews    int `json:"reviews"`
	Zones      int `json:"zones"`
	LapRecords int `json:"lapRecords"`
}

type TrackListItem struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Location     string      `json:"location"`
	State        *string     `json:"state"`
	Description  *string     `json:"description"`
	ImageURL     *string     `json:"imageUrl"`
	Latitude     *float64    `json:"latitude"`
	Longitude    *float64    `json:"longitude"`
	Status       string      `json:"status"`
	IsImported   bool        `json:"isImported"`
	UploadedByID string      `json:"uploadedById"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	Events       []TrackEvent `json:"events"`
	UploadedBy   UserBrief    `json:"uploadedBy"`
	Count        TrackCounts  `json:"_count"`
	AvgRating    float64      `json:"avgRating"`
}

type ZoneTipWithAuthor struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	Conditions *string   `json:"conditions"`
	ZoneID     string    `json:"zoneId"`
	AuthorID   string    `json:"authorId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Author     UserBrief `json:"author"`
}

type TrackZoneWithTips struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description *string             `json:"description"`
	PosX        float64             `json:"posX"`
	PosY        float64             `json:"posY"`
	TrackID     string              `json:"trackId"`
	EventType   *string             `json:"eventType"`
	CreatedAt   time.Time           `json:"createdAt"`
	Tips        []ZoneTipWithAuthor `json:"tips"`
}

type ReviewWithAuthor struct {
	ID           string        `json:"id"`
	Rating       int           `json:"rating"`
	Content      *string       `json:"content"`
	Conditions   string        `json:"conditions"`
	TrackID      string        `json:"trackId"`
	TrackEventID *string       `json:"trackEventId"`
	AuthorID     string        `json:"authorId"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Author       UserWithCars  `json:"author"`
	TrackEvent   *TrackEvent   `json:"trackEvent"`
}

type TrackDetail struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Location     string              `json:"location"`
	State        *string             `json:"state"`
	Description  *string             `json:"description"`
	ImageURL     *string             `json:"imageUrl"`
	Latitude     *float64            `json:"latitude"`
	Longitude    *float64            `json:"longitude"`
	Status       string              `json:"status"`
	IsImported   bool                `json:"isImported"`
	UploadedByID string              `json:"uploadedById"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	Events       []TrackEvent        `json:"events"`
	UploadedBy   UserWithExperience  `json:"uploadedBy"`
	Zones        []TrackZoneWithTips `json:"zones"`
	Reviews      []ReviewWithAuthor  `json:"reviews"`
	Count        TrackCounts         `json:"_count"`
	AvgRating    float64             `json:"avgRating"`
}

type TrackImageWithUploader struct {
	ID         string    `json:"id"`
	URL        string    `json:"url"`
	Caption    *string   `json:"caption"`
	CreatedAt  time.Time `json:"createdAt"`
	UploadedBy UserBrief `json:"uploadedBy"`
}

type LapRecordWithDetails struct {
	ID              string      `json:"id"`
	LapTime         string      `json:"lapTime"`
	Conditions      string      `json:"conditions"`
	Notes           *string     `json:"notes"`
	TirePressureFL  *float64    `json:"tirePressureFL"`
	TirePressureFR  *float64    `json:"tirePressureFR"`
	TirePressureRL  *float64    `json:"tirePressureRL"`
	TirePressureRR  *float64    `json:"tirePressureRR"`
	FuelLevel       *float64    `json:"fuelLevel"`
	CamberFL        *float64    `json:"camberFL"`
	CamberFR        *float64    `json:"camberFR"`
	CamberRL        *float64    `json:"camberRL"`
	CamberRR        *float64    `json:"camberRR"`
	CasterFL        *float64    `json:"casterFL"`
	CasterFR        *float64    `json:"casterFR"`
	ToeFL           *float64    `json:"toeFL"`
	ToeFR           *float64    `json:"toeFR"`
	ToeRL           *float64    `json:"toeRL"`
	ToeRR           *float64    `json:"toeRR"`
	TrackID         string      `json:"trackId"`
	TrackEventID    *string     `json:"trackEventId"`
	CarID           string      `json:"carId"`
	DriverID        string      `json:"driverId"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	Track           TrackBrief  `json:"track"`
	TrackEvent      *TrackEvent `json:"trackEvent"`
	Car             CarWithID   `json:"car"`
}

type ProfileResponse struct {
	ID         string    `json:"id"`
	Name       *string   `json:"name"`
	Email      string    `json:"email"`
	Experience string    `json:"experience"`
	Image      *string   `json:"image"`
	CreatedAt  time.Time `json:"createdAt"`
	Cars       []Car     `json:"cars"`
	Count      struct {
		TrackReviews int `json:"trackReviews"`
		LapRecords   int `json:"lapRecords"`
		Tracks       int `json:"tracks"`
		ZoneTips     int `json:"zoneTips"`
	} `json:"_count"`
}

// ─── Request DTOs ───────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string  `json:"id"`
		Name  *string `json:"name"`
		Email string  `json:"email"`
	} `json:"user"`
}

type ProfileUpdateRequest struct {
	Name       string `json:"name"`
	Experience string `json:"experience"`
}

type CarRequest struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type CarModRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Notes    *string `json:"notes"`
}

type TrackRequest struct {
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	Description *string  `json:"description"`
	ImageURL    *string  `json:"imageUrl"`
	EventTypes  []string `json:"eventTypes"`
}

type TrackPatchRequest struct {
	ImageURL    *string `json:"imageUrl"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type TrackImageRequest struct {
	URL     string  `json:"url"`
	Caption *string `json:"caption"`
}

type TrackZoneRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	PosX        float64 `json:"posX"`
	PosY        float64 `json:"posY"`
	EventType   *string `json:"eventType"`
}

type ZoneUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ZoneTipRequest struct {
	Content    string  `json:"content"`
	Conditions *string `json:"conditions"`
}

type TrackReviewRequest struct {
	Rating       int     `json:"rating"`
	Content      *string `json:"content"`
	Conditions   string  `json:"conditions"`
	TrackEventID *string `json:"trackEventId"`
}

type LapRecordRequest struct {
	LapTime        string   `json:"lapTime"`
	Conditions     string   `json:"conditions"`
	Notes          *string  `json:"notes"`
	TirePressureFL *float64 `json:"tirePressureFL"`
	TirePressureFR *float64 `json:"tirePressureFR"`
	TirePressureRL *float64 `json:"tirePressureRL"`
	TirePressureRR *float64 `json:"tirePressureRR"`
	FuelLevel      *float64 `json:"fuelLevel"`
	CamberFL       *float64 `json:"camberFL"`
	CamberFR       *float64 `json:"camberFR"`
	CamberRL       *float64 `json:"camberRL"`
	CamberRR       *float64 `json:"camberRR"`
	CasterFL       *float64 `json:"casterFL"`
	CasterFR       *float64 `json:"casterFR"`
	ToeFL          *float64 `json:"toeFL"`
	ToeFR          *float64 `json:"toeFR"`
	ToeRL          *float64 `json:"toeRL"`
	ToeRR          *float64 `json:"toeRR"`
	TrackID        string   `json:"trackId"`
	TrackEventID   *string  `json:"trackEventId"`
	CarID          string   `json:"carId"`
}

// ─── Admin ──────────────────────────────────────────────────────────────────────

type ImportedTrack struct {
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	State       string   `json:"state"`
	Types       []string `json:"types"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Description string   `json:"description"`
}

type SyncTracksResponse struct {
	Status  string `json:"status"`
	Summary struct {
		Total   int `json:"total"`
		Created int `json:"created"`
		Updated int `json:"updated"`
		Failed  int `json:"failed"`
	} `json:"summary"`
	Errors []string `json:"errors,omitempty"`
}
