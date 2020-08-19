package sqlitestore

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/writhe/kudosaurus"
)

// Kudo - kudo model
type Kudo struct {
	gorm.Model
	Text        string
	TeamSlackID string
	Owner       Person `gorm:"foreignkey:OwnerID"`
	OwnerID     uint
	Author      Person `gorm:"foreignkey:AuthorID"`
	AuthorID    uint
	IsPublished bool
}

// GetData - gets datasource.Kudo data from the wrapper
func (k *Kudo) GetData() kudosaurus.Kudo {
	return kudosaurus.Kudo{
		ID:          fmt.Sprint(k.ID),
		TeamID:      k.TeamSlackID,
		Text:        k.Text,
		PersonFrom:  k.Author.GetData(),
		PersonTo:    k.Owner.GetData(),
		IsPublished: k.IsPublished,
	}
}
