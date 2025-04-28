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
	BaseModel `bun:"table:motorbikes,alias:motorbike"`

	Model             string           `json:"model"`
	LocationLatitude  float64          `json:"location_latitude"`
	LocationLongitude float64          `json:"location_longitude"`
	Photos            []MotorbikePhoto `json:"photos" bun:"rel:has-many,join:ID=MotorbikeID"`
	Status            MotorBikeStatus  `json:"status" bun:"type:varchar(20)"`
	LockStatus        LockStatus       `json:"lock_status" bun:"type:varchar(10)"`
}

type MotorbikePhoto struct {
	BaseModel `bun:"table:motorbike_photo,alias:motorbike_photo"`

	MotorbikeID int    `gorm:"not null"`
	PhotoURL    string `gorm:"type:varchar(255);not null"`
}

func (Motorbike) TableName() string {
	return "motorbike"
}

func (MotorbikePhoto) TableName() string { return "motorbike_photo" }

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
