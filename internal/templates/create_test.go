package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// todo
func TestCreatePage_GeneratePage(t *testing.T) {
	page := NewCreatePage()
	html, err := page.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err, "GeneratePage should not return an error")
	assert.NotEmpty(t, html, "Generated HTML should not be empty")

	assert.Contains(t, html, "<title>Create New Post</title>", "HTML should contain the correct title")
	assert.Contains(t, html, `<h1>Create New Post</h1>`, "HTML should contain the main heading")
	assert.Contains(t, html, `<form id="post-form"`, "HTML should contain the post form")
	assert.Contains(t, html, `<label for="title">Title</label>`, "HTML should contain title label")
	assert.Contains(t, html, `<input type="text" id="title" name="title"`, "HTML should contain title input")
	assert.Contains(t, html, `<label for="content">Content</label>`, "HTML should contain content label")
	assert.Contains(t, html, `<textarea id="content" name="content"`, "HTML should contain content textarea")
	assert.Contains(t, html, `<button type="submit" id="submit-btn" class="btn btn-primary">`, "HTML should contain submit button")
	assert.Contains(t, html, `<a href="/home" class="btn btn-secondary">Main Menu</a>`, "HTML should contain main menu button")

	assert.NotContains(t, html, `<button type="submit" id="submit-btn" class="btn btn-primary" disabled>`, "Submit button should not be disabled by default")
}
