package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/writhe/kudosaurus/internal/misc"
)

// InteractivityPayload - workaround for missing nlopes/slack block types
type InteractivityPayload struct {
	Type    string `json:"block_actions"`
	Actions []struct {
		ActionID      string   `json:"action_id"`
		SelectedUsers []string `json:"selected_users"`
	} `json:"actions"`
}

func (s *CommandServer) handleDeleteKudo(
	id string,
	slackUser slack.User,
	res http.ResponseWriter,
) {
	k, err := s.store.GetKudo(id)

	if err != nil {
		s.closeWithError(err, http.StatusNotFound, res)
		return
	}

	if k.PersonFrom.ID != slackUser.ID {
		s.closeWithError(errors.New("this is not your kudo to delete"), http.StatusUnauthorized, res)
		return
	}

	if err := s.store.RemoveKudo(id); err != nil {
		s.closeWithError(err, http.StatusInternalServerError, res)
		return
	}

	res.WriteHeader(http.StatusOK)
	s.publishView(s.getHomeView(slackUser.ID, slackUser.TeamID))
}

func (s *CommandServer) handleChooseAdmins(
	adminIDs []string,
	slackUser slack.User,
	res http.ResponseWriter,
) {
	_, foundTeam := s.store.GetTeam(slackUser.TeamID)

	if !foundTeam {
		s.closeWithError(fmt.Errorf("no such team - %s", slackUser.TeamID), http.StatusBadRequest, res)
		return
	}

	if isAdmin, _ := s.store.CheckAdmin(slackUser.TeamID, slackUser.ID); !isAdmin {
		s.closeWithError(errors.New("user is not an admin"), http.StatusUnauthorized, res)
		return
	}

	s.store.SetAdmins(slackUser.TeamID, adminIDs)
	res.WriteHeader(http.StatusOK)
}

func (s *CommandServer) handleChooseChannel(
	channelID string,
	slackUser slack.User,
	res http.ResponseWriter,
) {
	team, foundTeam := s.store.GetTeam(slackUser.TeamID)

	if !foundTeam {
		s.closeWithError(fmt.Errorf("no such team - %s", slackUser.TeamID), http.StatusBadRequest, res)
		return
	}

	if isAdmin, _ := s.store.CheckAdmin(slackUser.TeamID, slackUser.ID); !isAdmin {
		s.closeWithError(errors.New("user is not an admin"), http.StatusUnauthorized, res)
		return
	}

	team.Settings.TargetChannelID = channelID
	s.store.PutTeam(team)
	res.WriteHeader(http.StatusOK)
	s.publishView(s.getHomeView(slackUser.ID, slackUser.TeamID))
}

func (s *CommandServer) handlePublishKudos(
	channelID string,
	slackUser slack.User,
	res http.ResponseWriter,
) {
	if isAdmin, _ := s.store.CheckAdmin(slackUser.TeamID, slackUser.ID); !isAdmin {
		s.closeWithError(errors.New("user is not an admin"), http.StatusUnauthorized, res)
		return
	}

	kudos := s.store.GetKudos(slackUser.TeamID)
	team, teamFound := s.store.GetTeam(slackUser.TeamID)

	if !teamFound {
		s.closeWithError(fmt.Errorf("no such team - %s", slackUser.TeamID), http.StatusBadRequest, res)
		return
	}

	msg := "*This month's kudos:*"

	for _, k := range kudos {
		msg += fmt.Sprintf("\n%s <@%s|%s> - %s", misc.RandomBullet(), k.PersonTo.ID, k.PersonTo.Name, k.Text)
	}

	s.api.PostMessage(
		team.Settings.TargetChannelID,
		slack.MsgOptionText(msg, false),
	)

	res.WriteHeader(http.StatusOK)
}

func (s *CommandServer) handleAction(res http.ResponseWriter, req *http.Request) {
	if err := s.verifySigningSecret(req, res); err != nil {
		return
	}

	var payload slack.InteractionCallback

	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)

	if err != nil {
		s.logger.Printf("Could not parse action response JSON: %v", err)
	}

	action := payload.ActionCallback.BlockActions[0]

	if action.ActionID == "delete_kudo" {
		s.handleDeleteKudo(action.Value, payload.User, res)
		return
	}

	if action.ActionID == "publish_kudos" {
		s.handlePublishKudos(payload.Channel.Name, payload.User, res)
		return
	}

	if action.ActionID == "choose_admins" {
		var parsedJSON InteractivityPayload
		body := req.FormValue("payload")
		err := json.Unmarshal([]byte(body), &parsedJSON)

		if err != nil {
			s.logger.Println(err)
			return
		}

		s.handleChooseAdmins(parsedJSON.Actions[0].SelectedUsers, payload.User, res)
		return
	}

	if action.ActionID == "choose_channel" {
		s.handleChooseChannel(action.SelectedChannel, payload.User, res)
		return
	}

	s.logger.Printf("Unknown action - '%s'\n", action.ActionID)
}
