package model

type MotorBikeStatus string

const (
	BikeAvailable     MotorBikeStatus = "available"
	BikeInMaintenance MotorBikeStatus = "maintenance"
	BikeRented        MotorBikeStatus = "rented"
)

type LockStatus string

const (
	Locked   LockStatus = "locked"
	Unlocked LockStatus = "unlocked"
)

// Motorbike modeli
type Motorbike struct {
	BaseModel `bun:"table:motorbikes,alias:m"`

	Model             string           `json:"model" bun:"model"`
	LocationLatitude  float64          `json:"location_latitude" bun:"location_latitude"`
	LocationLongitude float64          `json:"location_longitude" bun:"location_longitude"`
	Photos            []MotorbikePhoto `json:"photos" bun:"rel:has-many,join:id=motorbike_id"`
	Status            MotorBikeStatus  `json:"status" bun:"status,type:motorbike_status"`
	LockStatus        LockStatus       `json:"lock_status" bun:"lock_status,type:lock_status"`
}

type MotorbikePhoto struct {
	BaseModel `bun:"table:motorbike_photos,alias:mp"`

	MotorbikeID int64  `json:"motorbike_id" bun:"motorbike_id,notnull"`
	PhotoURL    string `json:"photo_url" bun:"photo_url,notnull"`
}

func (Motorbike) TableName() string {
	return "motorbikes"
}

func (MotorbikePhoto) TableName() string {
	return "motorbike_photos"
}

func (r MotorBikeStatus) String() string {
	switch r {
	case BikeAvailable:
		return "available"
	case BikeInMaintenance:
		return "maintenance"
	case BikeRented:
		return "rented"
	default:
		return "unknown"
	}
}

func (r LockStatus) String() string {
	switch r {
	case Locked:
		return "locked"
	case Unlocked:
		return "unlocked"
	default:
		return "unknown"
	}
}
