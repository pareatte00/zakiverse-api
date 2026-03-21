package card

import (
	"github.com/go-jet/jet/v2/qrm"
)

type Repository struct {
	db qrm.DB
}

func New(db qrm.DB) *Repository {
	return &Repository{db: db}
}
