package templates

import (
	"bytes"
	"html/template"
	"newsteller/internal/models"
)

const listHTML = `<div id="posts-content">
        <div class="posts-container">
            {{range .Posts}}
            <a href="/posts/{{.ID.Hex}}" class="post-card">
                <div class="post-title">{{truncateContent .Title 25}}</div>
                <div class="post-content">{{truncateContent .Content 32}}</div>
                <div class="post-date">{{formatDate .CreatedAt}}</div>
            </a>
            {{else}}
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: #666;">
                No posts found.
            </div>
            {{end}}
        </div>

        {{if gt .TotalPages 1}}
        <div class="pagination">
            <button 
                hx-get="/posts?page={{.PrevPage}}" 
                hx-target="#posts-content"
                hx-indicator=".loading"
								hx-include=".search-field"
                {{if not .HasPrev}}disabled{{end}}>
                Previous
            </button>
            
            <span class="page-info">Page {{.CurrentPage}} of {{.TotalPages}}</span>
            
            <button 
                hx-get="/posts?page={{.NextPage}}" 
                hx-target="#posts-content"
                hx-indicator=".loading"
								hx-include=".search-field"
                {{if not .HasNext}}disabled{{end}}>
                Next
            </button>
        </div>
        {{end}}
    </div>`

type List struct {
	posts        []models.Post
	currentPage  int
	totalPosts   int
	postsPerPage int
}

func NewList(
	posts []models.Post,
	currentPage int,
	totalPosts int,
	postsPerPage int,
) *List {
	return &List{
		posts:        posts,
		currentPage:  currentPage,
		totalPosts:   totalPosts,
		postsPerPage: postsPerPage,
	}
}

func (l *List) GeneratePage() (string, error) {
	totalPages := (l.totalPosts + l.postsPerPage - 1) / l.postsPerPage

	if totalPages == 0 {
		totalPages = 1
	}

	pageData := pageData{
		Posts:       l.posts,
		CurrentPage: l.currentPage,
		TotalPages:  totalPages,
		HasPrev:     l.currentPage > 1,
		HasNext:     l.currentPage < totalPages,
		PrevPage:    l.currentPage - 1,
		NextPage:    l.currentPage + 1,
	}

	tmpl, err := template.New("list").Funcs(funcMap).Parse(listHTML)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pageData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
