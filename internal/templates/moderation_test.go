package templates

import (
	"fmt"
	"math"
	"newsteller/internal/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModeration_GeneratePage_WithPostsAndPagination(t *testing.T) {
	mockPosts := createMockPosts(5)
	moderationPage := NewModeration(mockPosts, 1, 3, 10)
	html, err := moderationPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	// 1. Check for main elements
	assert.Contains(t, html, "<title>Edit Posts</title>", "HTML should contain the correct title") // Title is "Edit Posts"
	assert.Contains(t, html, `<h1>Edit Posts</h1>`, "HTML should contain the main heading")
	assert.Contains(t, html, `href="/" class="back-link"`, "HTML should contain back to home link")
	assert.Contains(t, html, `<div id="posts-content">`, "HTML should contain posts-content div")
	assert.Contains(t, html, `<ul class="posts-list">`, "HTML should contain posts-list ul")
	assert.Contains(t, html, `<div class="pagination">`, "HTML should contain pagination div")
	assert.Contains(t, html, `<h2>All Posts (10 total)</h2>`, "HTML should show total posts count")

	// 2. Check changeable elements (posts)
	for _, post := range mockPosts {
		assert.Contains(t, html, fmt.Sprintf(`<h3 class="post-title">%s</h3>`, truncateContent(post.Title, 45)), "Post title mismatch or not found")
		assert.Contains(t, html, fmt.Sprintf(`<p class="post-date">Created: %s</p>`, formatDate(post.CreatedAt)), "Post date mismatch or not found")
		assert.Contains(t, html, fmt.Sprintf(`href="/posts/%s/edit" class="btn btn-edit">Edit</a>`, post.ID.Hex()), "Edit button link incorrect")
		assert.Contains(t, html, fmt.Sprintf(`onclick="confirmDelete('%s', '%s')`, post.ID.Hex(), post.Title), "Delete button onclick incorrect")
	}

	// Check pagination info
	totalPages := int(math.Ceil(float64(10) / float64(3))) // 4
	assert.Contains(t, html, fmt.Sprintf(`Page 1 of %d`, totalPages), "HTML should show correct page info")

	// 3. Check button states
	// Previous button should be disabled on page 1
	assert.Contains(t, html, `hx-get="/posts/edit?page=0"`, "Previous button link incorrect") // PrevPage is CurrentPage - 1
	assert.Contains(t, html, `disabled`, "Previous button should be disabled on page 1")
	// Next button should be enabled
	assert.Contains(t, html, `hx-get="/posts/edit?page=2"`, "Next button link incorrect")
	assert.NotContains(t, html, `hx-get="/posts/edit?page=2" [...] disabled`, "Next button should be enabled") // This check is tricky with string contains.
	// A more robust check would parse the HTML or use regex. For now, check it's not disabled.
	nextButtonHtml := fmt.Sprintf(`<button hx-get="/posts/edit?page=2" hx-target="body" hx-indicator=".loading-indicator" > Next </button>`)
	assert.True(t, strings.Contains(html, nextButtonHtml) || strings.Contains(html, strings.ReplaceAll(nextButtonHtml, "\"", "'")), "Next button HTML structure is not as expected or is disabled")

}

func TestModeration_GeneratePage_NoPosts(t *testing.T) {
	moderationPage := NewModeration([]models.Post{}, 1, 5, 0) // No posts

	html, err := moderationPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, `<div class="empty-state">`, "HTML should show empty state")
	assert.Contains(t, html, `<h3>No posts found</h3>`, "HTML should show 'No posts found' message")
	assert.Contains(t, html, `<a href="/posts/create">Create your first post</a>`, "HTML should contain create post link in empty state")
	assert.NotContains(t, html, `<ul class="posts-list">`, "Posts list should not be present")
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present if no posts or only one page")
	assert.Contains(t, html, `<h2>All Posts (0 total)</h2>`, "HTML should show 0 total posts count")
}

func TestModeration_GeneratePage_SinglePageOfPosts(t *testing.T) {
	mockPosts := createMockPosts(2)
	moderationPage := NewModeration(mockPosts, 1, 5, 2)

	html, err := moderationPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	for _, post := range mockPosts {
		assert.Contains(t, html, post.ID.Hex())
	}
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present for a single page of results")
	assert.Contains(t, html, `<h2>All Posts (2 total)</h2>`, "HTML should show 2 total posts count")
}

func TestModeration_GeneratePage_LastPage(t *testing.T) {
	mockPosts := createMockPosts(1)
	totalPages := int(math.Ceil(float64(10) / float64(3))) // 4
	moderationPage := NewModeration(mockPosts, totalPages, 3, 10)
	html, err := moderationPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, fmt.Sprintf(`Page %d of %d`, totalPages, totalPages))
	assert.Contains(t, html, fmt.Sprintf(`hx-get="/posts/edit?page=%d"`, totalPages-1))
	prevButtonHtml := fmt.Sprintf(`<button hx-get="/posts/edit?page=%d" hx-target="body" hx-indicator=".loading-indicator" > Previous </button>`, totalPages-1)
	assert.True(t, strings.Contains(html, prevButtonHtml) || strings.Contains(html, strings.ReplaceAll(prevButtonHtml, "\"", "'")), "Previous button HTML structure is not as expected or is disabled")

	assert.Contains(t, html, fmt.Sprintf(`hx-get="/posts/edit?page=%d"`, totalPages+1))
	assert.Contains(t, html, `disabled`, "Next button should be disabled on the last page")
}

func TestModeration_GeneratePage_MiddlePage(t *testing.T) {
	mockPosts := createMockPosts(3)
	currentPage := 2
	totalPages := int(math.Ceil(float64(10) / float64(3))) // 4
	moderationPage := NewModeration(mockPosts, currentPage, 3, 10)
	html, err := moderationPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, fmt.Sprintf(`Page %d of %d`, currentPage, totalPages))
	assert.Contains(t, html, fmt.Sprintf(`hx-get="/posts/edit?page=%d"`, currentPage-1))
	prevButtonHtml := fmt.Sprintf(`<button hx-get="/posts/edit?page=%d" hx-target="body" hx-indicator=".loading-indicator" > Previous </button>`, currentPage-1)
	assert.True(t, strings.Contains(html, prevButtonHtml) || strings.Contains(html, strings.ReplaceAll(prevButtonHtml, "\"", "'")), "Previous button HTML structure is not as expected or is disabled")

	assert.Contains(t, html, fmt.Sprintf(`hx-get="/posts/edit?page=%d"`, currentPage+1))
	nextButtonHtml := fmt.Sprintf(`<button hx-get="/posts/edit?page=%d" hx-target="body" hx-indicator=".loading-indicator" > Next </button>`, currentPage+1)
	assert.True(t, strings.Contains(html, nextButtonHtml) || strings.Contains(html, strings.ReplaceAll(nextButtonHtml, "\"", "'")), "Next button HTML structure is not as expected or is disabled")
}
