package pkg

import (
	"reflect"
	"testing"

	"github.com/dstotijn/go-notion"
)

func TestExtractRichText(t *testing.T) {
	richText := []notion.RichText{
		{PlainText: "Hello"},
		{PlainText: "World"},
	}

	blocks := []notion.Block{
		&notion.ParagraphBlock{
			RichText: []notion.RichText{
				{PlainText: "Hello"},
				{PlainText: "World"},
			},
		},
		&notion.Heading1Block{
			RichText: []notion.RichText{
				{PlainText: "Heading 1"},
			},
		},
		&notion.Heading2Block{
			RichText: []notion.RichText{
				{PlainText: "Heading 2"},
			},
		},
		&notion.Heading3Block{
			RichText: []notion.RichText{
				{PlainText: "Heading 3"},
			},
		},
		&notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{PlainText: "Item 1"},
			},
		},
	}

	pageWithBlocks := &PageWithBlocks{
		Blocks: blocks,
	}

	expected := "HelloWorldHeading 1Heading 2Heading 3Item 1"
	result := pageWithBlocks.NormalizeBody()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
	expected = "HelloWorld"
	result = extractRichText(richText)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestBlockToMarkdown(t *testing.T) {
	paragraphBlock := &notion.ParagraphBlock{
		RichText: []notion.RichText{
			{PlainText: "Hello"},
			{PlainText: "World"},
		},
	}
	expected := "HelloWorld"
	result := blockToMarkdown(paragraphBlock)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	heading1Block := &notion.Heading1Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 1"},
		},
	}
	expected = "Heading 1"
	result = blockToMarkdown(heading1Block)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	heading2Block := &notion.Heading2Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 2"},
		},
	}
	expected = "Heading 2"
	result = blockToMarkdown(heading2Block)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	heading3Block := &notion.Heading3Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 3"},
		},
	}
	expected = "Heading 3"
	result = blockToMarkdown(heading3Block)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	bulletedListItemBlock := &notion.BulletedListItemBlock{
		RichText: []notion.RichText{
			{PlainText: "Item 1"},
		},
	}
	expected = "Item 1"
	result = blockToMarkdown(bulletedListItemBlock)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestNormalizeBody(t *testing.T) {
	blocks := []notion.Block{
		&notion.ParagraphBlock{
			RichText: []notion.RichText{
				{PlainText: "Hello"},
				{PlainText: "World"},
			},
		},
		&notion.Heading1Block{
			RichText: []notion.RichText{
				{PlainText: "Heading 1"},
			},
		},
		&notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{PlainText: "Item 1"},
			},
		},
	}

	pageWithBlocks := &PageWithBlocks{
		Blocks: blocks,
	}

	expected := "HelloWorldHeading 1Item 1"
	result := pageWithBlocks.NormalizeBody()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestTagsToNotionProps(t *testing.T) {
	tags := []string{"Tag1", "Tag2", "Tag3"}
	expected := []notion.SelectOptions{
		{Name: "Tag1"},
		{Name: "Tag2"},
		{Name: "Tag3"},
	}
	result := tagsToNotionProps(tags)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, result)
	}
}
