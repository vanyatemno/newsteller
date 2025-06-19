package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"newsteller/internal/models"
)

type homePageData struct {
	RecentPosts []models.Post `json:"recent_posts"`
}

type Main struct {
	posts []models.Post
}

func (m *Main) GeneratePage() (string, error) {
	tmpl, err := template.New("main").Funcs(funcMap).Parse(mainPage)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(
		&buf,
		homePageData{
			RecentPosts: m.posts,
		},
	); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func NewMain(posts []models.Post) *Main {
	return &Main{posts}
}

const mainPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog Home</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        
        .header {
            text-align: center;
            margin-bottom: 40px;
        }
        
        .header h1 {
            color: #333;
            font-size: 2.5em;
            font-weight: 600;
            margin: 0 0 10px 0;
        }
        
        .header p {
            color: #666;
            font-size: 1.1em;
            margin: 0;
        }
        
        .actions-section {
            background: white;
            border-radius: 12px;
            padding: 25px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 40px;
        }
        
        .actions-title {
            font-size: 1.3em;
            font-weight: 600;
            color: #333;
            margin-bottom: 20px;
        }
        
        .actions-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }
        
        .action-card {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
            text-decoration: none;
            color: inherit;
            transition: all 0.2s ease;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 10px;
        }
        
        .action-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
            text-decoration: none;
            color: inherit;
            border-color: #007bff;
        }
        
        .action-icon {
            width: 40px;
            height: 40px;
            background: #007bff;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 18px;
            font-weight: bold;
        }
        
        .action-title {
            font-weight: 600;
            color: #333;
            margin: 0;
        }
        
        .action-description {
            font-size: 0.9em;
            color: #666;
            margin: 0;
        }
        
        .recent-posts-section {
            background: white;
            border-radius: 12px;
            padding: 25px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        
        .section-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 25px;
        }
        
        .section-title {
            font-size: 1.3em;
            font-weight: 600;
            color: #333;
            margin: 0;
        }
        
        .view-all-link {
            color: #007bff;
            text-decoration: none;
            font-weight: 500;
            font-size: 0.9em;
            transition: color 0.2s ease;
        }
        
        .view-all-link:hover {
            color: #0056b3;
            text-decoration: underline;
        }
        
        .posts-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
            gap: 20px;
        }
        
        .post-card {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            transition: all 0.2s ease;
            cursor: pointer;
            text-decoration: none;
            color: inherit;
            display: block;
        }
        
        .post-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
            text-decoration: none;
            color: inherit;
            border-color: #007bff;
        }
        
        .post-title {
            font-size: 1.2em;
            font-weight: 600;
            margin-bottom: 10px;
            color: #333;
            line-height: 1.3;
        }
        
        .post-content {
            color: #666;
            line-height: 1.5;
            margin-bottom: 12px;
            font-size: 0.9em;
        }
        
        .post-date {
            color: #999;
            font-size: 0.8em;
            font-weight: 500;
        }
        
        .no-posts {
            text-align: center;
            color: #666;
            font-style: italic;
            padding: 40px 20px;
        }
        
        .loading {
            text-align: center;
            padding: 20px;
            color: #666;
            display: none;
        }
        
        @media (max-width: 768px) {
            body {
                padding: 15px;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .actions-grid {
                grid-template-columns: 1fr;
            }
            
            .posts-grid {
                grid-template-columns: 1fr;
            }
            
            .section-header {
                flex-direction: column;
                gap: 10px;
                align-items: flex-start;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Welcome to Your Blog</h1>
        <p>Manage your posts and share your thoughts</p>
    </div>

    <div class="actions-section">
        <h2 class="actions-title">Quick Actions</h2>
        <div class="actions-grid">
            <a href="/posts/create" class="action-card">
                <div class="action-icon">+</div>
                <h3 class="action-title">Create Post</h3>
                <p class="action-description">Write a new blog post</p>
            </a>
            
            <a href="/posts/edit" class="action-card">
                <div class="action-icon">✎</div>
                <h3 class="action-title">Edit Posts</h3>
                <p class="action-description">Manage existing posts</p>
            </a>
            
            <a href="/posts/search" class="action-card">
                <div class="action-icon">📄</div>
                <h3 class="action-title">All Posts</h3>
                <p class="action-description">Browse all your posts</p>
            </a>
        </div>
    </div>

    <div class="recent-posts-section">
        <div class="section-header">
            <h2 class="section-title">Recent Posts</h2>
            <a href="/posts/search" class="view-all-link">View All Posts →</a>
        </div>
        
        <div id="recent-posts-container">
            <div class="posts-grid">
                {{range .RecentPosts}}
                <a href="/posts/{{.ID.Hex}}" class="post-card">
                    <div class="post-title">{{truncateContent .Title 21}}</div>
                    <div class="post-content">{{truncateContent .Content 32}}</div>
                    <div class="post-date">{{formatDate .CreatedAt}}</div>
                </a>
                {{else}}
                <div class="no-posts">
                    No posts yet. <a href="/posts/create">Create your first post!</a>
                </div>
                {{end}}
            </div>
        </div>
    </div>

    <div class="loading">
        Loading...
    </div>

    <script>
        // Handle HTMX loading states
        document.body.addEventListener('htmx:beforeRequest', function() {
            document.querySelector('.loading').style.display = 'block';
        });
        
        document.body.addEventListener('htmx:afterRequest', function() {
            document.querySelector('.loading').style.display = 'none';
        });

        // Add some interactive feedback
        document.querySelectorAll('.action-card, .post-card').forEach(card => {
            card.addEventListener('mousedown', function() {
                this.style.transform = 'translateY(-1px) scale(0.98)';
            });
            
            card.addEventListener('mouseup', function() {
                this.style.transform = 'translateY(-2px) scale(1)';
            });
            
            card.addEventListener('mouseleave', function() {
                this.style.transform = 'translateY(0) scale(1)';
            });
        });
    </script>
</body>
</html>`
