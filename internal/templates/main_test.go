package templates

import (
	"fmt"
	"newsteller/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_GeneratePage_WithPosts(t *testing.T) {
	mockPosts := createMockPosts(3) // Create 3 mock posts
	mainPage := NewMain(mockPosts)

	html, err := mainPage.GeneratePage()
	assert.NoError(t, err, "GeneratePage should not return an error")
	assert.NotEmpty(t, html, "Generated HTML should not be empty")

	// 1. Check for main elements
	assert.Contains(t, html, "<title>Blog Home</title>", "HTML should contain the correct title")
	assert.Contains(t, html, `<h1>Welcome to Your Blog</h1>`, "HTML should contain the main heading")
	assert.Contains(t, html, `<div class="actions-section">`, "HTML should contain actions section")
	assert.Contains(t, html, `<a href="/posts/create" class="action-card">`, "HTML should contain create post action card")
	assert.Contains(t, html, `<a href="/posts/edit" class="action-card">`, "HTML should contain edit posts action card")
	assert.Contains(t, html, `<a href="/posts/search" class="action-card">`, "HTML should contain all posts action card")
	assert.Contains(t, html, `<div class="recent-posts-section">`, "HTML should contain recent posts section")
	assert.Contains(t, html, `<h2 class="section-title">Recent Posts</h2>`, "HTML should contain recent posts title")

	// 2. Check changeable elements (recent posts)
	for _, post := range mockPosts {
		assert.Contains(t, html, fmt.Sprintf(`href="/posts/%s"`, post.ID.Hex()), "HTML should link to recent post")
		expectedTitle := truncateContent(post.Title, 21)
		expectedContent := truncateContent(post.Content, 32)
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-title">%s</div>`, expectedTitle), "Post title mismatch")
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-content">%s</div>`, expectedContent), "Post content mismatch")
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-date">%s</div>`, formatDate(post.CreatedAt)), "Post date mismatch")
	}
	assert.Contains(t, html, `<a href="/posts/search" class="view-all-link">View All Posts →</a>`, "HTML should contain view all posts link")
}

func TestMain_GeneratePage_NoPosts(t *testing.T) {
	mainPage := NewMain([]models.Post{})

	html, err := mainPage.GeneratePage()
	assert.NoError(t, err, "GeneratePage should not return an error")
	assert.NotEmpty(t, html, "Generated HTML should not be empty")

	assert.Contains(t, html, `<div class="no-posts">`, "HTML should contain no-posts div")
	assert.Contains(t, html, `No posts yet. <a href="/posts/create">Create your first post!</a>`, "HTML should show 'No posts yet' message with create link")
	assert.NotContains(t, html, `<div class="post-card">`, "HTML should not contain any post-card when no posts")
}
