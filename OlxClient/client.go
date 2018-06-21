package olxclient

import (
	"github.com/artemzi/olx-parser/entities"
)

func Run() []*entities.Adverts {
	adverts := parse()

	return adverts
}
