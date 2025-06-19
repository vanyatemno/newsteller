package templates

import (
	"fmt"
	"newsteller/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSingleTemplate_GeneratePage(t *testing.T) {
	postID := primitive.NewObjectID()
	now := time.Now()
	mockPost := &models.Post{
		ID:        postID,
		Title:     "Single Test Post Title",
		Content:   "This is the detailed content for the single test post. It should be displayed fully.",
		CreatedAt: now.Add(-48 * time.Hour),
		UpdatedAt: now.Add(-2 * time.Hour),
	}

	singlePage := &SingleTemplate{post: mockPost}
	html, err := singlePage.GeneratePage()

	assert.NoError(t, err, "GeneratePage should not return an error")
	assert.NotEmpty(t, html, "Generated HTML should not be empty")

	// 1. Check for main elements
	assert.Contains(t, html, fmt.Sprintf("<title>%s</title>", mockPost.Title), "HTML should contain the correct post title in the <title> tag")
	assert.Contains(t, html, `<article class="post">`, "HTML should contain the main article element")
	assert.Contains(t, html, `<header class="post-header">`, "HTML should contain post header")
	assert.Contains(t, html, `<div class="post-content">`, "HTML should contain post content div")
	assert.Contains(t, html, `<footer class="post-actions">`, "HTML should contain post actions footer")
	assert.Contains(t, html, `<button class="btn-secondary" onclick="window.location.href='/posts/search'">Back to posts</button>`, "HTML should contain back to posts button")

	// 2. Check changeable elements
	assert.Contains(t, html, fmt.Sprintf(`<h1 class="post-title">%s</h1>`, mockPost.Title), "HTML should display the correct post title in <h1>")
	assert.Contains(t, html, fmt.Sprintf(`<p>%s</p>`, mockPost.Content), "HTML should display the full post content")

	expectedCreatedAt := mockPost.CreatedAt.Format("January 2, 2006 at 3:04 PM")
	expectedUpdatedAt := mockPost.UpdatedAt.Format("January 2, 2006 at 3:04 PM")
	assert.Contains(t, html, fmt.Sprintf(`<span>Created: %s</span>`, expectedCreatedAt), "HTML should display the correct formatted CreatedAt date")
	assert.Contains(t, html, fmt.Sprintf(`<span> | Updated: %s</span>`, expectedUpdatedAt), "HTML should display the correct formatted UpdatedAt date")
}

func TestSingleTemplate_GeneratePage_NoUpdatedAt(t *testing.T) {
	postID := primitive.NewObjectID()
	now := time.Now()
	mockPost := &models.Post{
		ID:        postID,
		Title:     "Post Without Update",
		Content:   "Content here.",
		CreatedAt: now.Add(-72 * time.Hour),
	}

	singlePage := &SingleTemplate{post: mockPost}
	html, err := singlePage.GeneratePage()

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	expectedCreatedAt := mockPost.CreatedAt.Format("January 2, 2006 at 3:04 PM")
	assert.Contains(t, html, fmt.Sprintf(`<span>Created: %s</span>`, expectedCreatedAt))
	assert.NotContains(t, html, "<span> | Updated:", "HTML should not display UpdatedAt if it's zero")
}

func TestSingleTemplate_GeneratePage_NilPost(t *testing.T) {
	singlePage := &SingleTemplate{post: nil}
	html, err := singlePage.GeneratePage()

	assert.Error(t, err, "GeneratePage should return an error when post is nil")
	assert.True(t, strings.Contains(err.Error(), "failed to execute template"), "Error message should indicate template execution failure")
	assert.Empty(t, html, "Generated HTML should be empty on execution error")
}

func TestRenderSinglePost_Function(t *testing.T) {
	postID := primitive.NewObjectID()
	now := time.Now()
	mockPost := &models.Post{
		ID:        postID,
		Title:     "Render Test Title",
		Content:   "Render test content.",
		CreatedAt: now,
	}

	html, err := RenderSinglePost(mockPost)
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.Contains(t, html, mockPost.Title)
	assert.Contains(t, html, mockPost.Content)
}
