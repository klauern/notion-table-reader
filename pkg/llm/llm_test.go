package llm_test

import (
	"testing"

	"github.com/klauern/notion-table-reader/pkg/llm"
	. "github.com/onsi/gomega"
)

func TestGenerateSystemPrompt(t *testing.T) {
	RegisterTestingT(t)
	tags := []string{"tag1", "tag2", "tag3"}
	expected := `
		You are a command-line app that responds with only a list of tags that categorize the content of the messages being sent to you.
		You can only provide AT MOST 3 tags, and at least 1 TAG.  Less is preferable.  Do not infer tags. The list of tags you output are:
		- tag1
		- tag2
		- tag3
	`
	result := llm.GenerateSystemPrompt(tags)
	Expect(result).To(Equal(expected))
}

func TestGenerateTagInputMessage(t *testing.T) {
	input := &llm.TagInput{
		Title: "Test Title",
		URL:   "http://example.com",
		Raw:   "Test content",
	}
	expected := `
		Title: Test Title
		URL: http://example.com

		Content Raw: Test content
	`
	result := llm.GenerateTagInputMessage(input, 1000)
	Expect(result).To(Equal(expected))
}

func TestGenerateTagInputMessage_Truncate(t *testing.T) {
	input := &llm.TagInput{
		Title: "Test Title",
		URL:   "http://example.com",
		Raw:   "Test content",
	}
	expected := `
		Title: Test Title
		URL: http://example.com

		Content Raw: Test content
	`
	result := llm.GenerateTagInputMessage(input, 10)
	Expect(result).To(Equal(expected[:10]))
}

func TestSplitResponse(t *testing.T) {
	response := "line1\nline2\nline3"
	expected := []string{"line1", "line2", "line3"}
	result := llm.SplitResponse(response)
	Expect(result).To(Equal(expected))
}
