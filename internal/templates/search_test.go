package templates

import (
	"fmt"
	"newsteller/internal/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchHome_GeneratePage_WithPostsAndPagination(t *testing.T) {
	mockPosts := createMockPosts(5)
	searchPage := NewHome(mockPosts, 1, 10, 3)
	html, err := searchPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	// 1. Check for main elements
	assert.Contains(t, html, "<title>Posts</title>", "HTML should contain the correct title")
	assert.Contains(t, html, `<h1>Posts</h1>`, "HTML should contain the main heading")
	assert.Contains(t, html, `href="/" class="back-link"`, "HTML should contain back to home link")
	assert.Contains(t, html, `<input type="text" class="search-field"`, "HTML should contain search field")
	assert.Contains(t, html, `<div id="posts-content">`, "HTML should contain posts-content div")
	assert.Contains(t, html, `<div class="posts-container">`, "HTML should contain posts-container div")
	assert.Contains(t, html, `<div class="pagination">`, "HTML should contain pagination div")

	// 2. Check changeable elements (posts)
	for _, post := range mockPosts {
		assert.Contains(t, html, fmt.Sprintf(`href="/posts/%s"`, post.ID.Hex()), "HTML should link to post")
		expectedTitle := truncateContent(post.Title, 25)
		expectedContent := truncateContent(post.Content, 32)
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-title">%s</div>`, expectedTitle))
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-content">%s</div>`, expectedContent))
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-date">%s</div>`, formatDate(post.CreatedAt)))
	}

	assert.Contains(t, html, `Page 1 of 4`, "HTML should show correct page info")

	// 3. Check button states
	assert.Contains(t, html, `hx-get="/posts?page=0"`, "Previous button link incorrect")
	assert.Contains(t, html, `disabled`, "Previous button should be disabled on page 1")
	assert.Contains(t, html, `hx-get="/posts?page=2"`, "Next button link incorrect")
	assert.NotContains(t, html, `disabled> Next </button>`, "Next button should be enabled")
}

func TestSearchHome_GeneratePage_NoPosts(t *testing.T) {
	searchPage := NewHome([]models.Post{}, 1, 0, 5) // No posts

	html, err := searchPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, "No posts found.", "HTML should show 'No posts found' message")
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present if no posts or only one page")
}

func TestSearchHome_GeneratePage_SinglePageOfPosts(t *testing.T) {
	mockPosts := createMockPosts(2)
	searchPage := NewHome(mockPosts, 1, 2, 5)

	html, err := searchPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	for _, post := range mockPosts {
		assert.Contains(t, html, post.ID.Hex())
	}
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present for a single page of results")
}

func TestSearchHome_GeneratePage_LastPage(t *testing.T) {
	mockPosts := createMockPosts(1)
	searchPage := NewHome(mockPosts, 4, 10, 3)
	html, err := searchPage.GeneratePage()
	html = strings.Join(strings.Fields(html), " ")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, `Page 4 of 4`)
	assert.Contains(t, html, `hx-get="/posts?page=3"`)
	prevButtonHtml := `<button hx-get="/posts?page=3" hx-target="#posts-content" hx-indicator=".loading" hx-include=".search-field" > Previous </button>`
	assert.Contains(t, html, prevButtonHtml, "Previous button HTML structure is not as expected or is disabled")

	assert.Contains(t, html, `hx-get="/posts?page=5"`)
	assert.Contains(t, html, `disabled`, "Next button should be disabled on the last page")
}

func TestSearchHome_GeneratePage_CurrentPageNormalization(t *testing.T) {
	searchPage := NewHome([]models.Post{}, 0, 10, 3)
	html, err := searchPage.GeneratePage()
	assert.NoError(t, err)
	assert.Contains(t, html, `Page 1 of 4`, "Current page should be normalized to 1 if less than 1")

	searchPage = NewHome([]models.Post{}, 5, 10, 3) // Total pages is 4
	html, err = searchPage.GeneratePage()
	assert.NoError(t, err)
	assert.Contains(t, html, `Page 4 of 4`, "Current page should be normalized to totalPages if greater")
}
