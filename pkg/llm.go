package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

const (
	SystemPromptTemplate = `
		You are a command-line app that responds with only a list of tags that categorize the content of the messages being sent to you.
		You can only provide AT MOST 3 tags, and at least 1 TAG.  Less is preferable.  Do not infer tags. The list of tags you output are:

		{{- range .}}
		- {{.}}
		{{- end}}
	`

	TagInputTemplate = `
		Title: {{.Title}}
		URL: {{.URL}}

		Content Raw: {{.Raw}}
	`
)

type TagInput struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Raw   string `json:"raw"`
}

func GenerateSystemPrompt(tags []string) string {
	tmpl, err := template.New("system-prompt").Parse(SystemPromptTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tags)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func GenerateTagInputMessage(input *TagInput, tokenLimit int) string {
	tmpl, err := template.New("tag-input").Parse(TagInputTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		panic(err)
	}

	// Truncate the message to the token limit
	message := buf.String()
	if len(message) > tokenLimit {
		message = message[:tokenLimit]
	}

	return message
}

func (l Client) IdentifyTags(messageContent *TagInput, tagOptions []string) ([]string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: GenerateSystemPrompt(tagOptions),
		},
		{
			Role:    "user",
			Content: GenerateTagInputMessage(messageContent, l.MaxTokens),
		},
	}

	response, err := l.RequestChatCompletion(messages)
	if err != nil {
		return nil, err
	}

	return splitResponse(response), nil
}

func splitResponse(response string) []string {
	return strings.Split(response, "\n")
}

func (client Client) NewTagInput(pageID string) (*TagInput, error) {
	page, err := client.GetPage(pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page %s: %w", pageID, err)
	}
	tag := &TagInput{
		Title: page.Page.Properties.(notion.DatabasePageProperties)["Name"].Title[0].PlainText,
		URL:   page.Page.URL,
		Raw:   page.NormalizeBody(),
	}
	return tag, nil
}

func (client Client) TagPage(id string, availableTags []string) error {
	tag, err := client.NewTagInput(id)
	if err != nil {
		return fmt.Errorf("failed to retrive Notion Page: %w", err)
	}
	tagList, err := client.IdentifyTags(tag, availableTags)
	if err != nil {
		return fmt.Errorf("failed to identify tags for page %s: %w", id, err)
	}
	fmt.Printf("Tagging page %s with tags: %s\n", id, strings.Join(tagList, ", "))
	if err := client.TagDatabasePage(id, tagList); err != nil {
		return fmt.Errorf("failed to tag page %s: %w", id, err)
	}
	return nil
}
