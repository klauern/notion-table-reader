package pkg_test

import (
	"context"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/klauern/notion-table-reader/pkg"
	"github.com/klauern/notion-table-reader/pkg/mocks"
	myNotion "github.com/klauern/notion-table-reader/pkg/notion"
	"go.uber.org/mock/gomock"
)

func TestExtractRichText(t *testing.T) {
	RegisterTestingT(t)
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

	pageWithBlocks := &myNotion.PageWithBlocks{
		Blocks: blocks,
	}

	expected := "HelloWorldHeading 1Heading 2Heading 3Item 1"
	result := pageWithBlocks.NormalizeBody()

	Expect(result).To(Equal(expected))
	Expect(result).To(Equal(expected))
	expected = "HelloWorld"
	result = myNotion.ExtractRichText(richText)
	Expect(result).To(Equal(expected))
}

func TestListMultiSelectProps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotionClient := mocks.NewMockNotionClient(ctrl)
	client := pkg.NewClient(context.Background(), "", "")
	client.NotionClient = mockNotionClient

	databaseId := "test-database-id"
	columnName := "Tags"
	expectedProps := []string{"Tag1", "Tag2"}

	mockNotionClient.EXPECT().FindDatabaseByID(gomock.Any(), databaseId).Return(notion.Database{
		Properties: map[string]notion.DatabaseProperty{
			columnName: {
				Type: notion.DBPropTypeMultiSelect,
				Name: columnName,
				Select: &notion.SelectMetadata{
					Options: []notion.SelectOptions{
						{Name: "Tag1"},
						{Name: "Tag2"},
					},
				},
			},
		},
	}, nil)

	props, err := client.ListMultiSelectProps(databaseId, columnName)
	Expect(err).To(BeNil())
	Expect(props).To(Equal(expectedProps))
}

func TestListDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotionClient := mocks.NewMockNotionClient(ctrl)
	client := pkg.NewClient(context.Background(), "", "")
	client.NotionClient = mockNotionClient

	query := "test-query"
	expectedDatabases := []notion.Database{
		{Title: []notion.RichText{{PlainText: "Database 1"}}},
		{Title: []notion.RichText{{PlainText: "Database 2"}}},
	}

	mockNotionClient.EXPECT().Search(gomock.Any(), &notion.SearchOpts{
		Query: query,
		Filter: &notion.SearchFilter{
			Value:    "database",
			Property: "object",
		},
	}).Return(notion.SearchResponse{
		Results: []interface{}{
			notion.Database{Title: []notion.RichText{{PlainText: "Database 1"}}},
			notion.Database{Title: []notion.RichText{{PlainText: "Database 2"}}},
		},
	}, nil)

	databases, err := client.ListDatabases(query)
	Expect(err).To(BeNil())
	Expect(databases).To(Equal(expectedDatabases))
}

func TestListPages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotionClient := mocks.NewMockNotionClient(ctrl)
	client := pkg.NewClient(context.Background(), "", "")
	client.NotionClient = mockNotionClient

	databaseId := "test-database-id"
	expectedPages := []notion.Page{
		{ID: "page-1"},
		{ID: "page-2"},
	}

	mockNotionClient.EXPECT().QueryDatabase(gomock.Any(), databaseId, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "Tags",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				MultiSelect: &notion.MultiSelectDatabaseQueryFilter{
					IsEmpty: true,
				},
			},
		},
	}).Return(notion.DatabaseQueryResponse{
		Results: expectedPages,
	}, nil)

	pages, err := client.ListPages(databaseId, true)
	Expect(err).To(BeNil())
	Expect(pages).To(Equal(expectedPages))
}

func TestGetPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotionClient := mocks.NewMockNotionClient(ctrl)
	client := pkg.NewClient(context.Background(), "", "")
	client.NotionClient = mockNotionClient

	pageId := "test-page-id"
	expectedPage := &myNotion.PageWithBlocks{
		Page: &notion.Page{ID: pageId},
		Blocks: []notion.Block{
			&notion.ParagraphBlock{RichText: []notion.RichText{{PlainText: "Hello"}}},
		},
	}

	mockNotionClient.EXPECT().FindPageByID(gomock.Any(), pageId).Return(notion.Page{ID: pageId}, nil)
	mockNotionClient.EXPECT().FindBlockChildrenByID(gomock.Any(), pageId, &notion.PaginationQuery{}).Return(notion.BlockChildrenResponse{
		Results: []notion.Block{
			&notion.ParagraphBlock{RichText: []notion.RichText{{PlainText: "Hello"}}},
		},
	}, nil)

	page, err := client.GetPage(pageId)
	Expect(err).To(BeNil())
	Expect(page).To(Equal(expectedPage))
}

func TestTagDatabasePage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotionClient := mocks.NewMockNotionClient(ctrl)
	client := pkg.NewClient(context.Background(), "", "")
	client.NotionClient = mockNotionClient

	pageId := "test-page-id"
	tags := []string{"Tag1", "Tag2"}

	mockNotionClient.EXPECT().UpdatePage(gomock.Any(), pageId, notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"Tags": notion.DatabasePageProperty{
				MultiSelect: pkg.TagsToNotionProps(tags),
			},
		},
	}).Return(notion.Page{ID: pageId}, nil)

	err := client.TagDatabasePage(pageId, tags)
	Expect(err).To(BeNil())
}

func TestBlockToMarkdown(t *testing.T) {
	paragraphBlock := &notion.ParagraphBlock{
		RichText: []notion.RichText{
			{PlainText: "Hello"},
			{PlainText: "World"},
		},
	}
	expected := "HelloWorld"
	result := myNotion.BlockToMarkdown(paragraphBlock)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	heading1Block := &notion.Heading1Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 1"},
		},
	}
	expected = "Heading 1"
	result = myNotion.BlockToMarkdown(heading1Block)
	Expect(result).To(Equal(expected))

	heading2Block := &notion.Heading2Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 2"},
		},
	}
	expected = "Heading 2"
	result = myNotion.BlockToMarkdown(heading2Block)
	Expect(result).To(Equal(expected))

	heading3Block := &notion.Heading3Block{
		RichText: []notion.RichText{
			{PlainText: "Heading 3"},
		},
	}
	expected = "Heading 3"
	result = myNotion.BlockToMarkdown(heading3Block)
	Expect(result).To(Equal(expected))

	bulletedListItemBlock := &notion.BulletedListItemBlock{
		RichText: []notion.RichText{
			{PlainText: "Item 1"},
		},
	}
	expected = "Item 1"
	result = myNotion.BlockToMarkdown(bulletedListItemBlock)
	Expect(result).To(Equal(expected))
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

	pageWithBlocks := &myNotion.PageWithBlocks{
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
	result := pkg.TagsToNotionProps(tags)
	Expect(result).To(Equal(expected))
}
