package probe

import (
	"github.com/affanshahid/oversight/util"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Config represents a probe's config
type Config struct {
	util.BaseModel
	Interval      int64 //milliseconds
	Options       postgres.Jsonb
	Descriminator string
}
