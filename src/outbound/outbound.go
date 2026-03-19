package outbound

import (
	"github.com/zakiverse/zakiverse-api/config"
	"github.com/zakiverse/zakiverse-api/src/outbound/discord"
)

type Outbound struct {
	Discord *discord.Client
}

func New(conf config.ConfigConstant) *Outbound {
	return &Outbound{
		Discord: discord.New(conf),
	}
}

func (o *Outbound) Disconnect() {
	o.Discord.Disconnect()
}
