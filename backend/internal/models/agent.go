package models

import (
    "github.com/google/uuid"
    "time"
)

type Agent struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
    MacAddress      string    `gorm:"uniqueIndex"`
    Status          string    // unassociated, associated, deactivated
    WireGuardPeerID uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
