package pkg

import (
	"context"
	"reflect"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/stretchr/testify/mock"
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

type MockNotionClient struct {
	mock.Mock
}

func (m *MockNotionClient) FindDatabaseByID(ctx context.Context, databaseId string) (notion.Database, error) {
	args := m.Called(ctx, databaseId)
	return args.Get(0).(notion.Database), args.Error(1)
}

func (m *MockNotionClient) Search(ctx context.Context, opts *notion.SearchOpts) (notion.SearchResponse, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(notion.SearchResponse), args.Error(1)
}

func (m *MockNotionClient) FindPageByID(ctx context.Context, pageId string) (notion.Page, error) {
	args := m.Called(ctx, pageId)
	return args.Get(0).(notion.Page), args.Error(1)
}

func (m *MockNotionClient) FindBlockChildrenByID(ctx context.Context, blockId string, pagination *notion.PaginationQuery) (notion.BlockChildrenResponse, error) {
	args := m.Called(ctx, blockId, pagination)
	return args.Get(0).(notion.BlockChildrenResponse), args.Error(1)
}

func (m *MockNotionClient) UpdatePage(ctx context.Context, pageId string, params notion.UpdatePageParams) (notion.Page, error) {
	args := m.Called(ctx, pageId, params)
	return args.Get(0).(notion.Page), args.Error(1)
}

func TestListMultiSelectProps(t *testing.T) {
	mockNotionClient := new(MockNotionClient)
	client := &Client{notionClient: mockNotionClient, context: context.TODO()}

	databaseId := "test-database-id"
	columnName := "Tags"
	expectedProps := []string{"Tag1", "Tag2"}

	mockNotionClient.On("FindDatabaseByID", mock.Anything, databaseId).Return(notion.Database{
		Properties: map[string]notion.DatabaseProperty{
			columnName: {
				Type: notion.DBPropTypeMultiSelect,
				Name: columnName,
				Select: &notion.SelectMetadata{
					Options: []notion.Option{
						{Name: "Tag1"},
						{Name: "Tag2"},
					},
				},
			},
		},
	}, nil)

	props, err := client.ListMultiSelectProps(databaseId, columnName)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(props, expectedProps) {
		t.Errorf("Expected %v, but got %v", expectedProps, props)
	}
}

func TestListDatabases(t *testing.T) {
	mockNotionClient := new(MockNotionClient)
	client := &Client{notionClient: mockNotionClient, context: context.TODO()}

	query := "test-query"
	expectedDatabases := []notion.Database{
		{Title: []notion.RichText{{PlainText: "Database 1"}}},
		{Title: []notion.RichText{{PlainText: "Database 2"}}},
	}

	mockNotionClient.On("Search", mock.Anything, &notion.SearchOpts{
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
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(databases, expectedDatabases) {
		t.Errorf("Expected %v, but got %v", expectedDatabases, databases)
	}
}

func TestListPages(t *testing.T) {
	mockNotionClient := new(MockNotionClient)
	client := &Client{notionClient: mockNotionClient, context: context.TODO()}

	databaseId := "test-database-id"
	expectedPages := []notion.Page{
		{ID: "page-1"},
		{ID: "page-2"},
	}

	mockNotionClient.On("QueryDatabase", mock.Anything, databaseId, &notion.DatabaseQuery{
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
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(pages, expectedPages) {
		t.Errorf("Expected %v, but got %v", expectedPages, pages)
	}
}

func TestGetPage(t *testing.T) {
	mockNotionClient := new(MockNotionClient)
	client := &Client{notionClient: mockNotionClient, context: context.TODO()}

	pageId := "test-page-id"
	expectedPage := &PageWithBlocks{
		Page: &notion.Page{ID: pageId},
		Blocks: []notion.Block{
			&notion.ParagraphBlock{RichText: []notion.RichText{{PlainText: "Hello"}}},
		},
	}

	mockNotionClient.On("FindPageByID", mock.Anything, pageId).Return(notion.Page{ID: pageId}, nil)
	mockNotionClient.On("FindBlockChildrenByID", mock.Anything, pageId, &notion.PaginationQuery{}).Return(notion.BlockChildrenResponse{
		Results: []notion.Block{
			&notion.ParagraphBlock{RichText: []notion.RichText{{PlainText: "Hello"}}},
		},
	}, nil)

	page, err := client.GetPage(pageId)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if !reflect.DeepEqual(page, expectedPage) {
		t.Errorf("Expected %v, but got %v", expectedPage, page)
	}
}

func TestTagDatabasePage(t *testing.T) {
	mockNotionClient := new(MockNotionClient)
	client := &Client{notionClient: mockNotionClient, context: context.TODO()}

	pageId := "test-page-id"
	tags := []string{"Tag1", "Tag2"}

	mockNotionClient.On("UpdatePage", mock.Anything, pageId, notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"Tags": notion.DatabasePageProperty{
				MultiSelect: tagsToNotionProps(tags),
			},
		},
	}).Return(notion.Page{ID: pageId}, nil)

	err := client.TagDatabasePage(pageId, tags)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
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
	expected := []notion.Option{
		{Name: "Tag1"},
		{Name: "Tag2"},
		{Name: "Tag3"},
	}
	result := tagsToNotionProps(tags)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, result)
	}
}
