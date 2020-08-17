package templates

import (
	"fmt"
	"strings"

	"github.com/writhe/kudosaurus"
)

// Kudo - kudo view block
func Kudo(kudo kudosaurus.Kudo) string {
	template := `
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "· <@%s> - %s"
			},
			"accessory": {
				"type": "button",
				"action_id": "delete_kudo",
				"style": "danger",
				"text": {
					"type": "plain_text",
					"text": "Delete"
				},
				"value": "%s"
			}
		}	
	`

	return fmt.Sprintf(template, kudo.PersonTo.ID, kudo.Text, kudo.ID)
}

// KudoList - list of kudos
func KudoList(kudos []kudosaurus.Kudo) string {
	result := []string{}
	for _, kudo := range kudos {
		result = append(result, Kudo(kudo))
	}

	return strings.Join(result, ",")
}

// Admin - admin view template
func Admin(targetChannelID string, kudoCount int, adminIDs []string) string {
	template := `
	{
		"type": "section",
		"block_id": "channel_section",
		"text": {
			"type": "mrkdwn",
			"text": "*Select the channel that Kudosaurus will publish to*\n_#general? #random? Your call._"
		},
		"accessory": {
			"action_id": "choose_channel",
			"type": "channels_select",
			"initial_channel": "%s",
			"placeholder": {
				"type": "plain_text",
				"text": "Select channel",
				"emoji": true
			}
		}
	},
	{
		"type": "section",
		"block_id": "publish_section",
		"text": {
			"type": "mrkdwn",
			"text": "*Publish this month's kudos*\n_All *%d* of them will get published to <#%s> _"
		},
		"accessory": {
			"action_id": "publish_kudos",
			"type": "button",
			"style": "primary",
			"text": {
				"type": "plain_text",
				"text": "Publish!",
				"emoji": false
			},
			"value": "unused"
		}
	},
	{
		"type": "section",
		"block_id": "admin_section",
		"text": {
			"type": "mrkdwn",
			"text": "*Choose admins*\n_Admins can publish monthly lists of kudos_"
		},
		"accessory": {
			"action_id": "choose_admins",
			"type": "multi_users_select",
			"initial_users": [%s],
			"placeholder": {
				"type": "plain_text",
				"text": "Select admins"
			}
		}
	},
	{
		"type": "section",
		"text": {
			"type": "mrkdwn",
			"text": "Open modal (test)"
		},
		"accessory": {
			"type": "button",
			"text": {
				"type": "plain_text",
				"text": "Show all kudos",
				"emoji": true
			},
			"value": "show_all_kudos"
		}
	}
	`

	var encodedAdmins []string
	for _, adminID := range adminIDs {
		encodedAdmins = append(encodedAdmins, fmt.Sprintf("\"%s\"", adminID))
	}

	return fmt.Sprintf(
		template,
		targetChannelID,
		kudoCount,
		targetChannelID,
		strings.Join(encodedAdmins, ", "),
	)
}
