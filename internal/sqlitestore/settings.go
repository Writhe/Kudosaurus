package sqlitestore

import (
	"github.com/jinzhu/gorm"
	"github.com/writhe/kudosaurus"
)

// Settings - per-team config model
type Settings struct {
	gorm.Model
	TeamSlackID     string
	TargetChannelID string
}

func (s *Settings) getData() kudosaurus.Settings {
	return kudosaurus.Settings{
		TeamID:          s.TeamSlackID,
		TargetChannelID: s.TargetChannelID,
	}
}

func makeSettings(data kudosaurus.Settings) Settings {
	return Settings{
		TeamSlackID:     data.TeamID,
		TargetChannelID: data.TargetChannelID,
	}
}
