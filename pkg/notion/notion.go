package notion

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
	"github.com/klauern/notion-table-reader/pkg/llm"
)

//go:generate mockgen -destination=../mocks/mock_notion.go -package=mocks . NotionClient,NotionTableReader

type NotionTableReader interface {
	FetchPages(databaseID string, untagged bool) ([]PageDetail, error)
	TagPage(id string, availableTags []string) error
	NotionClient
}
type NotionClient interface {
	FindDatabaseByID(ctx context.Context, databaseId string) (notion.Database, error)
	Search(ctx context.Context, opts *notion.SearchOpts) (notion.SearchResponse, error)
	QueryDatabase(ctx context.Context, databaseId string, query *notion.DatabaseQuery) (notion.DatabaseQueryResponse, error)
	FindPageByID(ctx context.Context, pageId string) (notion.Page, error)
	FindBlockChildrenByID(ctx context.Context, blockId string, pagination *notion.PaginationQuery) (notion.BlockChildrenResponse, error)
	UpdatePage(ctx context.Context, pageId string, params notion.UpdatePageParams) (notion.Page, error)
}

type PageDetail struct {
	ID   string
	Name string
}

type PageWithBlocks struct {
	Page   *notion.Page
	Blocks []notion.Block
}

func (p PageWithBlocks) NormalizeBody() string {
	var buf bytes.Buffer
	for _, block := range p.Blocks {
		buf.WriteString(BlockToMarkdown(block))
	}
	return buf.String()
}

func BlockToMarkdown(block notion.Block) string {
	switch b := block.(type) {
	case *notion.ParagraphBlock:
		return ExtractRichText(b.RichText)
	case *notion.Heading1Block:
		return ExtractRichText(b.RichText)
	case *notion.Heading2Block:
		return ExtractRichText(b.RichText)
	case *notion.Heading3Block:
		return ExtractRichText(b.RichText)
	case *notion.BulletedListItemBlock:
		return ExtractRichText(b.RichText)
	case *notion.NumberedListItemBlock:
		return ExtractRichText(b.RichText)
	case *notion.ToDoBlock:
		return ExtractRichText(b.RichText)
	case *notion.CalloutBlock:
		return ExtractRichText(b.RichText)
	default:
		return ""
	}
}

func ExtractRichText(richText []notion.RichText) string {
	var buf bytes.Buffer
	for _, t := range richText {
		buf.WriteString(t.PlainText)
	}
	return buf.String()
}

func NewTagInput(page *PageWithBlocks) *llm.TagInput {
	tag := &llm.TagInput{
		Title: page.Page.Properties.(notion.DatabasePageProperties)["Name"].Title[0].PlainText,
		URL:   page.Page.URL,
		Raw:   page.NormalizeBody(),
	}
	return tag
}

// PrintPageDetails prints the details of the provided pages.
func PrintPageDetails(pages []PageDetail) {
	for _, page := range pages {
		fmt.Printf("Page(%s): %s\n", page.ID, page.Name)
	}
}
