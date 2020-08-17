package sqlitestore

import (
	"github.com/jinzhu/gorm"
	"github.com/writhe/kudosaurus"
)

// Person - person model
type Person struct {
	gorm.Model
	SlackID     string
	Name        string
	TeamSlackID string
	Team        Team `gorm:"foreignkey:TeamSlackID;association_foreignkey:SlackID"`
	IsAdmin     bool
}

// GetKudos - gets kudos
func (p Person) GetKudos(db *gorm.DB) []Kudo {
	var kudos []Kudo
	db.Where(
		"owner_id = ? AND is_published = 0",
		p.ID,
	).Preload("Author").Preload("Owner").Find(&kudos)

	return kudos
}

// GetIssuedKudos - gets issued kudos
func (p Person) GetIssuedKudos(db *gorm.DB) []Kudo {
	var kudos []Kudo
	db.Where(
		"author_id = ? AND is_published = 0",
		p.ID,
	).Preload("Author").Preload("Owner").Find(&kudos)

	return kudos
}

// GetFullData - gets full kudosaurus.Person data from the wrapper
func (p Person) GetFullData(db *gorm.DB) kudosaurus.Person {
	var mappedKudos []kudosaurus.Kudo
	var mappedIssuedKudos []kudosaurus.Kudo

	kudos := p.GetKudos(db)
	issuedKudos := p.GetIssuedKudos(db)
	kudosLeft := 10 - len(issuedKudos)

	for _, k := range kudos {
		mappedKudos = append(mappedKudos, k.GetData())
	}

	for _, k := range issuedKudos {
		mappedIssuedKudos = append(mappedIssuedKudos, k.GetData())
	}

	return kudosaurus.Person{
		Name:       p.Name,
		ID:         p.SlackID,
		Kudos:      mappedKudos,
		KudosLeft:  kudosLeft,
		KudosGiven: mappedIssuedKudos,
	}
}

// GetData - gets short kudosaurus.Person data from the wrapper
func (p Person) GetData() kudosaurus.Person {
	return kudosaurus.Person{
		Name: p.Name,
		ID:   p.SlackID,
	}
}
