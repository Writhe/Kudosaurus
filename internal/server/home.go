package server

import (
	"fmt"
	"strings"

	"github.com/writhe/kudosaurus/internal/templates"
)

const separatorBlock = `{"type": "divider"}`

func (s *CommandServer) getUserBlocks(teamID string, userID string) string {
	user, _ := s.store.GetPerson(teamID, userID, true)

	kudoBlocks := templates.KudoList(user.KudosGiven)

	blocks := []string{
		templates.Header(fmt.Sprintf("Kudos you gave this month (%d left)", user.KudosLeft)),
		separatorBlock,
		kudoBlocks,
	}

	return strings.Join(blocks, ",")
}

func (s *CommandServer) getAdminBlocks(teamID string, userID string) string {
	team, _ := s.store.GetTeam(teamID)
	kudos := s.store.GetKudos(teamID)
	adminIDs := s.store.GetAdmins(teamID)
	body := templates.Admin(team.Settings.TargetChannelID, len(kudos), adminIDs)

	blocks := []string{
		templates.Header("Admin"),
		separatorBlock,
		body,
	}

	return strings.Join(blocks, ",")
}

func (s *CommandServer) getHomeView(userID string, teamID string) string {
	isAdmin, _ := s.store.CheckAdmin(teamID, userID)
	template := `
	{
		"user_id": "%s",
		"view": {
			"type": "home",
			"blocks": [
				%s
			]
		}
	}
	`

	blocks := []string{s.getUserBlocks(teamID, userID)}

	if isAdmin {
		blocks = append(blocks, s.getAdminBlocks(teamID, userID))
	}

	return fmt.Sprintf(
		template,
		userID,
		strings.Join(blocks, ","),
	)
}
