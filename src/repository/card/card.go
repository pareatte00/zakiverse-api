package card

import (
	"github.com/go-jet/jet/v2/qrm"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
)

type Repository struct {
	db qrm.DB
}

type CardWithAnime struct {
	model.Card

	Anime   model.Anime
	CardTag *model.CardTag
}

func New(db qrm.DB) *Repository {
	return &Repository{db: db}
}
