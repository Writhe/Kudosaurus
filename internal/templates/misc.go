package templates

import "fmt"

// Header - header block
func Header(text string) string {
	template := `
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "%s",
				"emoji": true
			}
		}
	`

	return fmt.Sprintf(template, text)
}
