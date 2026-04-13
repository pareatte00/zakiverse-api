package pack

import (
	"github.com/go-jet/jet/v2/qrm"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
)

type Repository struct {
	db qrm.DB
}

type CardWithAnime struct {
	model.Card

	Anime  model.Anime
	TagName *string `alias:"card_tag.name"`
}

type PackCardWithCard struct {
	model.PackCard

	Card CardWithAnime
}

type PackWithCards struct {
	model.Pack

	Cards []PackCardWithCard
}

func New(db qrm.DB) *Repository {
	return &Repository{db: db}
}
