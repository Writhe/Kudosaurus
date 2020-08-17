package sqlitestore

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/writhe/kudosaurus"

	// NOTE: Needed for Sqlite3 support
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// SqliteStore - implements interface DataSource
type SqliteStore struct {
	db     *gorm.DB
	path   string
	logger *log.Logger
}

// Init initializes the database
func (s *SqliteStore) Init() {
	db, err := gorm.Open("sqlite3", s.path)
	if err != nil {
		s.logger.Panic(err)
	}
	s.db = db

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Kudo{})
	db.AutoMigrate(&Settings{})
	db.AutoMigrate(&Team{})
}

// PutPerson puts person
func (s *SqliteStore) PutPerson(id string, name string, teamID string, teamName string) (kudosaurus.Person, error) {
	_, found := s.GetPerson(teamID, id, false)

	if found {
		s.logger.Printf("ID '%s' is already taken.\n", id)
		return kudosaurus.Person{}, errors.New("Person already exists")
	}

	_, found = s.GetTeam(teamID)

	if !found {
		s.PutTeam(kudosaurus.Team{
			ID:   teamID,
			Name: teamName,
			Settings: kudosaurus.Settings{
				TeamID:          teamID,
				TargetChannelID: "placeholder",
			},
		})
	}

	p := Person{SlackID: id, Name: name, TeamSlackID: teamID}
	s.db.Create(&p)

	result, _ := s.GetPerson(teamID, id, true)

	return result, nil
}

// PutKudo puts kudo
func (s *SqliteStore) PutKudo(teamID string, idFrom string, idTo string, text string) {
	author, authorFound := s.getPerson(teamID, idFrom)
	owner, ownerFound := s.getPerson(teamID, idTo)

	if !(authorFound && ownerFound) {
		s.logger.Printf("No such id '%s'.\n", idTo)
		return
	}

	s.db.Create(&Kudo{
		Author:      author,
		Owner:       owner,
		Text:        text,
		TeamSlackID: teamID,
		IsPublished: false,
	})
}

func (s *SqliteStore) getKudo(id string) (Kudo, error) {
	var kudo Kudo
	s.db.Preload("Author").Preload("Owner").First(&kudo, id)

	if kudo.ID == 0 {
		return Kudo{}, fmt.Errorf("No such kudo: %s", id)
	}

	return kudo, nil
}

// GetKudo - removes a kudo
func (s *SqliteStore) GetKudo(id string) (kudosaurus.Kudo, error) {
	kudo, err := s.getKudo(id)

	if err != nil {
		return kudosaurus.Kudo{}, err
	}

	return kudo.GetData(), nil
}

func beginningOfMonth(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
}

func beginningOfCurrentMonth() int64 {
	return beginningOfMonth(time.Now())
}

func beginningOfPreviousMonth() int64 {
	return beginningOfMonth(time.Unix(beginningOfCurrentMonth()-1000, 0))
}

// GetKudos - gets all kudos from the current month
func (s *SqliteStore) GetKudos(teamID string) []kudosaurus.Kudo {
	var kudos []Kudo
	var mappedKudos []kudosaurus.Kudo

	s.db.Where("team_slack_id = ? AND is_published = 0", teamID).Preload("Author").Preload("Owner").Find(&kudos)

	for _, kudo := range kudos {
		mappedKudos = append(mappedKudos, kudo.GetData())
	}

	return mappedKudos
}

func (s *SqliteStore) getPerson(teamID string, id string) (Person, bool) {
	var person Person
	s.db.First(&person, "team_slack_id = ? AND slack_id = ?", teamID, id)

	return person, person.SlackID != ""
}

// RemoveKudo - removes a kudo
func (s *SqliteStore) RemoveKudo(id string) error {
	kudo, err := s.getKudo(id)

	if err == nil {
		s.db.Delete(kudo)
	}

	return err
}

// PublishKudos - marks kudos as published
func (s *SqliteStore) PublishKudos(teamID string) {
	s.db.Exec("UPDATE kudos SET is_published = 1 WHERE team_slack_id = ?", teamID)
}

// GetPerson - gets Person by id from DB
func (s *SqliteStore) GetPerson(teamID string, id string, includeKudos bool) (kudosaurus.Person, bool) {
	person, found := s.getPerson(teamID, id)

	if includeKudos {
		return person.GetFullData(s.db), found
	}
	return person.GetData(), found
}

func (s *SqliteStore) getTeam(slackID string) (Team, bool) {
	var team Team
	s.db.Where("slack_id = ?", slackID).Preload("People").Preload("Settings").First(&team)

	return team, team.Name != ""
}

// GetTeam - gets Team by SlackID
func (s *SqliteStore) GetTeam(slackID string) (kudosaurus.Team, bool) {
	team, found := s.getTeam(slackID)

	return team.getData(), found
}

// GetTeams returns a list of teams
func (s *SqliteStore) GetTeams() []kudosaurus.Team {
	var teams []Team
	var result []kudosaurus.Team

	s.db.Preload("People").Find(&teams)

	for _, team := range teams {
		result = append(result, team.getData())
	}

	return result
}

// PutTeam - upserts a Team
func (s *SqliteStore) PutTeam(newTeam kudosaurus.Team) {
	team, found := s.getTeam(newTeam.ID)
	if found {
		team.Settings.TargetChannelID = newTeam.Settings.TargetChannelID
	} else {
		team.Name = newTeam.Name
		team.SlackID = newTeam.ID
		team.Settings = makeSettings(newTeam.Settings)
	}

	s.db.Save(&team)
}

// SetAdmins - sets admins
func (s *SqliteStore) SetAdmins(teamID string, adminIDs []string) {
	s.db.Exec("UPDATE people SET is_admin = 0 WHERE team_slack_id = ?", teamID)
	s.db.Exec("UPDATE people SET is_admin = 1 WHERE team_slack_id = ? AND slack_id IN (?)", teamID, adminIDs)
}

// GetAdmins - gets admin IDs
func (s *SqliteStore) GetAdmins(teamID string) []string {
	var users []Person
	var result []string

	s.db.Where("team_slack_id = ? AND is_admin = 1", teamID).Find(&users)

	for _, user := range users {
		result = append(result, user.SlackID)
	}

	return result
}

// CheckAdmin - checks if userID belongs to an admin
func (s *SqliteStore) CheckAdmin(teamID string, userID string) (isAdmin bool, isFound bool) {
	adminIDs := s.GetAdmins(teamID)
	if len(adminIDs) == 0 { // NOTE: Admin party! If no one is an admin, everyone is!
		return true, false
	}
	user, foundUser := s.getPerson(teamID, userID)

	return user.IsAdmin, foundUser
}

// NewSource returns a new SqlSource
func NewSource(path string, logger *log.Logger) *SqliteStore {
	source := SqliteStore{path: path, logger: logger}
	source.Init()

	return &source
}
