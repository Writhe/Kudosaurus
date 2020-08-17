package sqlitestore

import (
	"github.com/jinzhu/gorm"
	"github.com/writhe/kudosaurus"
)

// Team - person model
type Team struct {
	gorm.Model
	SlackID  string
	Name     string
	People   []Person `gorm:"foreignkey:TeamSlackID;association_foreignkey:SlackID"`
	Settings Settings `gorm:"foreignkey:TeamSlackID;association_foreignkey:SlackID"`
}

func (t *Team) getData() kudosaurus.Team {
	people := []kudosaurus.Person{}

	for _, person := range t.People {
		people = append(people, person.GetData())
	}

	return kudosaurus.Team{
		ID:       t.SlackID,
		Name:     t.Name,
		People:   people,
		Settings: t.Settings.getData(),
	}
}
