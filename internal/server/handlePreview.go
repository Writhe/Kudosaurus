package server

import (
	"fmt"
	"strings"

	"github.com/writhe/kudosaurus"
	"github.com/writhe/kudosaurus/internal/templates"
)

func getKudoOptionBlock(kudo kudosaurus.Kudo) string {
	template := `
	{
		"text": {
			"type": "mrkdwn",
			"text": "%s"
		},
		"description": {
			"type": "mrkdwn",
			"text": "%s"
		},
		"value": "%s"
	}
	`

	return fmt.Sprintf(template, kudo.PersonFrom, kudo.Text, kudo.ID)
}

func getKudoListBlock(kudos []kudosaurus.Kudo) string {
	template := `
	{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "Select kudos to delete."
		},
		"accessory": {
			"type": "checkboxes",
			"options": [%s]
		}
	}
	`
	options := []string{}
	for _, kudo := range kudos {
		options = append(options, templates.Kudo(kudo))
	}

	return fmt.Sprintf(template, strings.Join(options, ","))
}

func getModalTemplate(kudos []kudosaurus.Kudo) string {
	template := `
	{
		"type": "modal",
		"submit": {
			"type": "plain_text",
			"text": "Delete selected",
			"emoji": true
		},
		"close": {
			"type": "plain_text",
			"text": "Close",
			"emoji": true
		},
		"title": {
			"type": "plain_text",
			"text": "This month's kudos",
			"emoji": false
		},
		"blocks": [%s]
	}
	`

	return fmt.Sprintf(template, getKudoListBlock(kudos))
}
