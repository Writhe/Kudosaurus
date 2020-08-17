package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/writhe/kudosaurus"
)

// CommandServerConfig - config for CommandServer
type CommandServerConfig struct {
	MaxKudos          int
	Token             string
	SigningSecret     string
	Port              int
	VerificationToken string
}

// CommandServer - Kudosaurus' Slack slash-command server
type CommandServer struct {
	store  kudosaurus.Store
	api    *slack.Client
	config CommandServerConfig
	logger *log.Logger
}

func (s *CommandServer) getUserName(id string) (string, bool) {
	user, err := s.api.GetUserInfo(id)

	if err != nil {
		return "", false
	}

	return user.Name, true
}

// Start - starts CommandServer
func (s *CommandServer) Start() {
	s.api = slack.New(
		s.config.Token,
		slack.OptionDebug(true),
		slack.OptionLog(s.logger),
	)

	http.HandleFunc("/slash", s.handleCommand)
	http.HandleFunc("/interactivity", s.handleAction)
	http.HandleFunc("/events", s.handleEvent)

	s.logger.Printf("Listening on port %d\n", s.config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), nil)
}

func (s *CommandServer) publishView(pageContent string) {
	r, _ := http.NewRequest("POST", "https://slack.com/api/views.publish", bytes.NewBuffer([]byte(pageContent)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.Token))

	client := &http.Client{}
	client.Do(r)
}

func (s *CommandServer) closeWithError(err error, statusCode int, res http.ResponseWriter) error {
	s.logger.Printf(err.Error())
	res.WriteHeader(statusCode)
	return err
}

func (s *CommandServer) verifySigningSecret(req *http.Request, res http.ResponseWriter) error {
	signingSecret := s.config.SigningSecret
	verifier, err := slack.NewSecretsVerifier(req.Header, signingSecret)
	if err != nil {
		return s.closeWithError(err, http.StatusUnauthorized, res)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return s.closeWithError(err, http.StatusInternalServerError, res)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		return s.closeWithError(err, http.StatusUnauthorized, res)
	}

	return nil
}

func (s *CommandServer) sendJSON(payload interface{}, res http.ResponseWriter) {
	b, err := json.Marshal(payload)
	if err != nil {
		s.closeWithError(err, http.StatusInternalServerError, res)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(b)
}

// Respond - sends a response in the form of a Slack message
func (s *CommandServer) Respond(body string, w http.ResponseWriter) {
	s.sendJSON(&slack.Msg{Text: body}, w)
}

// RespondBlocks - sends a response in the form of a blockwise Slack message
func (s *CommandServer) RespondBlocks(w http.ResponseWriter, blocks ...slack.Block) {
	s.sendJSON(&slack.Msg{Blocks: slack.Blocks{BlockSet: blocks}}, w)
}

// RespondMessage - sends a response in the form of a Slack message
func (s *CommandServer) RespondMessage(w http.ResponseWriter, msg slack.Message) {
	s.sendJSON(msg, w)
}

// New - returns new instance of CommandServer
func New(store kudosaurus.Store, config CommandServerConfig, logger *log.Logger) CommandServer {
	server := CommandServer{store: store, config: config, logger: logger}

	return server
}
