package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// DiscordWebhookResponse is the API message for Discord webhook response.
type DiscordWebhookResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DiscordWebhookEmbedField is the API message for Discord webhook embed field.
type DiscordWebhookEmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// DiscordWebhookEmbedAuthor is the API message for Discord webhook embed Author.
type DiscordWebhookEmbedAuthor struct {
	Name string `json:"name"`
}

// DiscordWebhookEmbed is the API message for Discord webhook embed.
type DiscordWebhookEmbed struct {
	Title       string                     `json:"title"`
	Type        string                     `json:"type"`
	Description string                     `json:"description,omitempty"`
	URL         string                     `json:"url,omitempty"`
	Timestamp   string                     `json:"timestamp"`
	Author      DiscordWebhookEmbedAuthor  `json:"author"`
	FieldList   []DiscordWebhookEmbedField `json:"fields,omitempty"`
}

// DiscordWebhook is the API message for Discord webhook.
type DiscordWebhook struct {
	EmbedList []DiscordWebhookEmbed `json:"embeds"`
}

func init() {
	Register("bb.plugin.webhook.discord", &DiscordReceiver{})
}

// DiscordReceiver is the receiver for Discord.
type DiscordReceiver struct {
}

func (*DiscordReceiver) Post(context Context) error {
	embedList := []DiscordWebhookEmbed{}

	fieldList := []DiscordWebhookEmbedField{}
	for _, meta := range context.GetMetaList() {
		fieldList = append(fieldList, DiscordWebhookEmbedField(meta))
	}

	status := ""
	switch context.Level {
	case WebhookSuccess:
		status = ":white_check_mark: "
	case WebhookWarn:
		status = ":warning: "
	case WebhookError:
		status = ":exclamation: "
	default:
		// No status icon for other levels
		status = ""
	}
	embedList = append(embedList, DiscordWebhookEmbed{
		Title:       fmt.Sprintf("%s%s", status, context.Title),
		Type:        "rich",
		Description: context.Description,
		URL:         context.Link,
		Author: DiscordWebhookEmbedAuthor{
			Name: fmt.Sprintf("%s (%s)", context.ActorName, context.ActorEmail),
		},
		FieldList: fieldList,
	})

	post := DiscordWebhook{
		EmbedList: embedList,
	}
	body, err := json.Marshal(post)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal webhook POST request to %s", context.URL)
	}
	req, err := http.NewRequest("POST",
		context.URL, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrapf(err, "failed to construct webhook POST request to %s", context.URL)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to POST webhook to %s", context.URL)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read POST webhook response from %s", context.URL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.Errorf("failed to POST webhook %s, status code: %d, response body: %s", context.URL, resp.StatusCode, b)
	}

	webhookResponse := &DiscordWebhookResponse{}
	if err := json.Unmarshal(b, webhookResponse); err != nil {
		return errors.Wrapf(err, "malformed webhook response from %s", context.URL)
	}

	if webhookResponse.Code != 0 {
		return errors.Errorf("%s", webhookResponse.Message)
	}

	return nil
}
