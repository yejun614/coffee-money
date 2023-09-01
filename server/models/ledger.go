package models

import (
	"gorm.io/gorm"
)

type Ledger struct {
	gorm.Model
	StoreName   string
	Balance     int
	Description string
	Username    string
	IsDisabled  bool
}
