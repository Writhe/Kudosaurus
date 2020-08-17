package kudosaurus

// Kudo - kudo entity
type Kudo struct {
	ID          string
	TeamID      string
	Text        string
	PersonFrom  Person
	PersonTo    Person
	IsPublished bool
}

// Person - person entity
type Person struct {
	ID         string
	Name       string
	Kudos      []Kudo
	KudosGiven []Kudo
	KudosLeft  int
	IsAdmin    bool
}

// Team - team entity
type Team struct {
	ID       string
	Name     string
	People   []Person
	Settings Settings
}

// Settings - runtime settings for a team
type Settings struct {
	TeamID          string
	TargetChannelID string
}

// Store - stores data
type Store interface {
	PutPerson(id string, name string, teamID string, teamName string) (Person, error)
	PutKudo(teamID string, idFrom string, idTo string, text string)
	GetKudo(id string) (Kudo, error)
	GetPerson(teamID string, id string, includeKudos bool) (Person, bool)
	RemoveKudo(id string) error
	GetKudos(teamID string) []Kudo
	GetTeams() []Team
	PutTeam(team Team)
	GetTeam(id string) (Team, bool)
	SetAdmins(teamID string, adminIDs []string)
	GetAdmins(teamID string) []string
	CheckAdmin(teamID string, userID string) (isAdmin bool, isFound bool)
	PublishKudos(teamID string)
}
