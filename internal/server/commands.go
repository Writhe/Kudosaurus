package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/writhe/kudosaurus/internal/misc"
)

func (s *CommandServer) handleKudo(c slack.SlashCommand, w http.ResponseWriter) {
	parts := strings.Fields(c.Text)

	userToken := parts[0]
	re := regexp.MustCompile(`<@(([A-Z0-9])+)(\|([a-z0-9\.]){0,21})?>`)
	targetUserMatches := re.FindStringSubmatch(userToken)
	if len(targetUserMatches) < 2 {
		fmt.Println("No user ID given")
		return
	}

	targetUserID := targetUserMatches[1]
	targetUserName, foundTargetUser := s.getUserName(targetUserID)

	if !foundTargetUser {
		fmt.Println("No target user")
		return
	}

	body := strings.Join(parts[1:], " ")

	personFrom, found := s.store.GetPerson(c.TeamID, c.UserID, true)
	if !found {
		personFrom, _ = s.store.PutPerson(c.UserID, c.UserName, c.TeamID, c.TeamDomain)
	}

	if personFrom.KudosLeft <= 0 {
		s.Respond("You have no kudos left. Bummer.\n", w)
		return
	}

	personTo, found := s.store.GetPerson(c.TeamID, targetUserID, false)
	if !found {
		personTo, _ = s.store.PutPerson(targetUserID, targetUserName, c.TeamID, c.TeamDomain)
	}

	s.store.PutKudo(c.TeamID, personFrom.ID, personTo.ID, body)

	var message string

	if targetUserID != personFrom.ID {
		message = fmt.Sprintf(
			"You give kudos to <@%s>. %s You have %d kudos left, by the way.",
			targetUserID,
			misc.RandomExclamation(),
			personFrom.KudosLeft-1,
		)
	} else {
		message = fmt.Sprintf(
			"You pat yourself on the back. Good job. Not sad at all. You have %d kudos left, by the way.",
			personFrom.KudosLeft-1,
		)
	}

	s.Respond(
		message,
		w,
	)
}

func (s *CommandServer) handleCommand(res http.ResponseWriter, req *http.Request) {
	if err := s.verifySigningSecret(req, res); err != nil {
		return
	}

	c, err := slack.SlashCommandParse(req)

	if err != nil {
		s.closeWithError(err, http.StatusInternalServerError, res)
		return
	}

	switch c.Command {
	case "/kudo":
		s.handleKudo(c, res)

	default:
		s.closeWithError(fmt.Errorf("No such command - %s", c.Command), http.StatusInternalServerError, res)
	}
}
