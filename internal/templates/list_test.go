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

func createMockPosts(count int) []models.Post {
	posts := make([]models.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = models.Post{
			ID:        primitive.NewObjectID(),
			Title:     fmt.Sprintf("Post Title %d", i+1),
			Content:   fmt.Sprintf("This is the content for post %d. It might be a bit longer to test truncation.", i+1),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
	}
	return posts
}

func TestList_GeneratePage_WithPostsAndPagination(t *testing.T) {
	mockPosts := createMockPosts(5)          // 5 posts
	listPage := NewList(mockPosts, 1, 10, 3) // Current page 1, 10 total posts, 3 per page

	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	// 1. Check for main elements
	assert.Contains(t, html, `<div id="posts-content">`, "HTML should contain posts-content div")
	assert.Contains(t, html, `<div class="posts-container">`, "HTML should contain posts-container div")
	assert.Contains(t, html, `<div class="pagination">`, "HTML should contain pagination div")

	// 2. Check changeable elements (posts)
	for _, post := range mockPosts {
		assert.Contains(t, html, fmt.Sprintf(`href="/posts/%s"`, post.ID.Hex()), "HTML should link to post")
		// Check truncated title and content (using the logic from funcMap)
		expectedTitle := truncateContent(post.Title, 25)
		expectedContent := truncateContent(post.Content, 32)
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-title">%s</div>`, expectedTitle))
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-content">%s</div>`, expectedContent))
		assert.Contains(t, html, fmt.Sprintf(`<div class="post-date">%s</div>`, formatDate(post.CreatedAt)))
	}

	// Check pagination info
	assert.Contains(t, html, `Page 1 of 4`, "HTML should show correct page info (10 posts, 3 per page -> 4 pages)") // 10/3 = 3.33 -> 4 pages

	// 3. Check button states
	// Previous button should be disabled on page 1
	assert.Contains(t, html, `hx-get="/posts?page=0"`, "Previous button link incorrect") // PrevPage is CurrentPage - 1
	assert.Contains(t, html, `disabled`, "Previous button should be disabled on page 1")
	// Next button should be enabled
	assert.Contains(t, html, `hx-get="/posts?page=2"`, "Next button link incorrect")
	assert.NotContains(t, html, `hx-get="/posts?page=2" [...] disabled`, "Next button should be enabled")
}

func TestList_GeneratePage_NoPosts(t *testing.T) {
	listPage := NewList([]models.Post{}, 1, 0, 5) // No posts

	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, "No posts found.", "HTML should show 'No posts found' message")
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present if no posts or only one page")
}

func TestList_GeneratePage_SinglePageOfPosts(t *testing.T) {
	mockPosts := createMockPosts(2)
	listPage := NewList(mockPosts, 1, 2, 5) // 2 posts, 5 per page -> 1 page total

	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	for _, post := range mockPosts {
		assert.Contains(t, html, post.ID.Hex())
	}
	assert.NotContains(t, html, `<div class="pagination">`, "Pagination should not be present for a single page of results")
}

func TestList_GeneratePage_LastPage(t *testing.T) {
	mockPosts := createMockPosts(1) // last post on the page
	listPage := NewList(mockPosts, 4, 10, 3)

	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, `Page 4 of 4`)
	assert.Contains(t, html, `hx-get="/posts?page=3"`)
	assert.NotContains(t, html, `hx-get="/posts?page=3" [...] disabled`)
	assert.Contains(t, html, `hx-get="/posts?page=5"`)
	assert.Contains(t, html, `disabled`, "Next button should be disabled on the last page")
}

func TestList_GeneratePage_MiddlePage(t *testing.T) {
	mockPosts := createMockPosts(3)
	listPage := NewList(mockPosts, 2, 10, 3) // Page 2 of 4

	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotEmpty(t, html)

	assert.Contains(t, html, `Page 2 of 4`)
	assert.Contains(t, html, `hx-get="/posts?page=1"`)
	assert.NotContains(t, html, `hx-get="/posts?page=1" [...] disabled`)
	assert.Contains(t, html, `hx-get="/posts?page=3"`)
	assert.NotContains(t, html, `hx-get="/posts?page=3" [...] disabled`)
}

func TestList_GeneratePage_TotalPagesCalculation(t *testing.T) {
	listPage := NewList(createMockPosts(6), 1, 6, 3) // 6 posts, 3 per page -> 2 pages
	html, err := listPage.GeneratePage()
	assert.NoError(t, err)
	assert.Contains(t, html, `Page 1 of 2`)

	listPage = NewList([]models.Post{}, 1, 0, 3)
	html, err = listPage.GeneratePage()
	assert.NoError(t, err)
	assert.NotContains(t, html, `<div class="pagination">`)
	assert.Contains(t, html, "No posts found.")
}

func truncateContent(content string, maxCharacters int) string {
	if len(content) < maxCharacters {
		return content
	}
	words := strings.Fields(content)
	if len(words) > 6 && len(strings.Join(words[:6], " ")) < maxCharacters {
		return strings.Join(words[:6], " ") + "..."
	}
	return content[:maxCharacters] + "..."
}

func formatDate(t time.Time) string {
	return t.Format("Jan 02, 2006")
}
