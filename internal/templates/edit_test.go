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

func TestEdit_GeneratePage(t *testing.T) {
	postID := primitive.NewObjectID()
	now := time.Now()
	mockPost := &models.Post{
		ID:        postID,
		Title:     "Test Post Title",
		Content:   "This is the test post content.",
		CreatedAt: now.Add(-24 * time.Hour), // Yesterday
		UpdatedAt: now,                      // Now
	}

	editPage := NewEdit(mockPost)
	html, err := editPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err, "GeneratePage should not return an error")
	assert.NotEmpty(t, html, "Generated HTML should not be empty")

	// 1. Check for main elements
	assert.Contains(t, html, "<title>Edit Post</title>", "HTML should contain the correct title")
	assert.Contains(t, html, `<h1>Edit Post</h1>`, "HTML should contain the main heading")
	assert.Contains(t, html, fmt.Sprintf(`<form id="edit-form" hx-put="/posts/%s"`, postID.Hex()), "HTML should contain the edit form with correct action URL")
	assert.Contains(t, html, `<label for="title">Title</label>`, "HTML should contain title label")
	assert.Contains(t, html, `<input type="text" id="title" name="title"`, "HTML should contain title input")
	assert.Contains(t, html, `<label for="content">Content</label>`, "HTML should contain content label")
	assert.Contains(t, html, `<textarea id="content" name="content"`, "HTML should contain content textarea")
	assert.Contains(t, html, `<button type="submit" id="save-btn" class="btn btn-primary">`, "HTML should contain save button")
	assert.Contains(t, html, `<a href="/" class="btn btn-secondary">Return to Home Page</a>`, "HTML should contain cancel button")

	// 2. Check changeable elements
	assert.Contains(t, html, fmt.Sprintf(`value="%s"`, mockPost.Title), "HTML should display the correct post title in input")
	assert.Contains(t, html, fmt.Sprintf(`>%s</textarea>`, mockPost.Content), "HTML should display the correct post content in textarea")
	assert.Contains(t, html, fmt.Sprintf(`<div class="metadata-value post-id">%s</div>`, mockPost.ID.Hex()), "HTML should display the correct Post ID")

	// Check formatted dates using the funcMap logic
	expectedCreatedAt := mockPost.CreatedAt.Format("January 2, 2006 at 3:04 PM")
	expectedUpdatedAt := mockPost.UpdatedAt.Format("January 2, 2006 at 3:04 PM")
	assert.Contains(t, html, expectedCreatedAt, "HTML should display the correct formatted CreatedAt date")
	assert.Contains(t, html, expectedUpdatedAt, "HTML should display the correct formatted UpdatedAt date")

	// 3. Check button states
	// The "Save Changes" button's disabled state is managed by JavaScript based on changes.
	// Initially, it should not be explicitly disabled in the raw HTML.
	assert.NotContains(t, html, `id="save-btn" class="btn btn-primary" disabled`, "Save button should not be disabled by default in raw HTML")
}

func TestEdit_GeneratePage_NilPost(t *testing.T) {
	// Test with a nil post, though the template might panic.
	// The `html/template` package will likely render empty strings for nil fields.
	editPage := NewEdit(nil) // Pass nil post
	html, err := editPage.GeneratePage()

	// Depending on how the template handles nil, it might error out or render empty.
	// The current template uses {{.ID.Hex}}, {{.Title}}, etc. which will cause a panic if . is nil.
	// The html/template package's Execute method recovers from such panics and returns an error.
	assert.Error(t, err, "GeneratePage should return an error when post is nil and template tries to access its fields")
	assert.True(t, strings.Contains(err.Error(), "failed to execute template"), "Error message should indicate template execution failure")
	assert.Empty(t, html, "Generated HTML should be empty on execution error")
}

func TestEdit_GeneratePage_EmptyPost(t *testing.T) {
	// Test with an empty (zero-value) post
	emptyPost := &models.Post{} // Zero values for fields
	editPage := NewEdit(emptyPost)
	html, err := editPage.GeneratePage()

	assert.NoError(t, err, "GeneratePage should not return an error with an empty post")
	assert.NotEmpty(t, html, "Generated HTML should not be empty with an empty post")

	// Check that placeholders or zero values are rendered
	assert.Contains(t, html, `value=""`, "Title input should be empty")          // Empty title
	assert.Contains(t, html, `></textarea>`, "Content textarea should be empty") // Empty content
	assert.Contains(t, html, `<div class="metadata-value post-id">000000000000000000000000</div>`, "HTML should display zero value for Post ID")

	// Check formatted zero dates
	zeroTimeFormatted := time.Time{}.Format("January 2, 2006 at 3:04 PM")
	assert.Contains(t, html, zeroTimeFormatted, "HTML should display formatted zero time for CreatedAt")
	assert.Contains(t, html, zeroTimeFormatted, "HTML should display formatted zero time for UpdatedAt")
}
