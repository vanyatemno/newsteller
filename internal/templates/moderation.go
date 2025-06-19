package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"newsteller/internal/models"
)

const moderationHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Edit Posts</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1000px;
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
            font-size: 2.2em;
            font-weight: 600;
            margin: 0 0 10px 0;
        }

        .header p {
            color: #666;
            font-size: 1em;
            margin: 0;
        }

        .back-link {
            display: inline-block;
            margin-bottom: 20px;
            color: #007bff;
            text-decoration: none;
            font-weight: 500;
            transition: color 0.2s ease;
        }

        .back-link:hover {
            color: #0056b3;
            text-decoration: underline;
        }

        .posts-container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .posts-header {
            background: #f8f9fa;
            padding: 20px;
            border-bottom: 1px solid #e9ecef;
        }

        .posts-header h2 {
            margin: 0;
            color: #333;
            font-size: 1.3em;
            font-weight: 600;
        }

        .posts-list {
            list-style: none;
            margin: 0;
            padding: 0;
        }

        .post-item {
            display: flex;
            align-items: center;
            padding: 20px;
            border-bottom: 1px solid #e9ecef;
            transition: background-color 0.2s ease;
        }

        .post-item:hover {
            background-color: #f8f9fa;
        }

        .post-item:last-child {
            border-bottom: none;
        }

        .post-info {
            flex: 1;
            min-width: 0;
        }

        .post-title {
            font-size: 1.1em;
            font-weight: 600;
            color: #333;
            margin: 0 0 5px 0;
            line-height: 1.3;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .post-date {
            color: #666;
            font-size: 0.9em;
            margin: 0;
        }

        .post-actions {
            display: flex;
            gap: 10px;
            margin-left: 20px;
        }

        .btn {
            padding: 8px 16px;
            border: none;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s ease;
            text-decoration: none;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            min-width: 70px;
        }

        .btn-edit {
            background: #007bff;
            color: white;
        }

        .btn-edit:hover {
            background: #0056b3;
            transform: translateY(-1px);
        }

        .btn-delete {
            background: #dc3545;
            color: white;
        }

        .btn-delete:hover {
            background: #c82333;
            transform: translateY(-1px);
        }

        .btn:disabled {
            background: #6c757d;
            cursor: not-allowed;
            transform: none;
        }

        .btn.loading {
            position: relative;
            color: transparent;
        }

        .btn.loading::after {
            content: '';
            position: absolute;
            top: 50%;
            left: 50%;
            width: 14px;
            height: 14px;
            margin: -7px 0 0 -7px;
            border: 2px solid transparent;
            border-top: 2px solid white;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .pagination {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 15px;
            margin-top: 30px;
            padding: 20px;
            background: white;
            border-radius: 12px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .pagination button {
            padding: 10px 20px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: background-color 0.2s ease;
        }

        .pagination button:hover:not(:disabled) {
            background: #0056b3;
        }

        .pagination button:disabled {
            background: #ddd;
            cursor: not-allowed;
            color: #999;
        }

        .page-info {
            font-weight: 500;
            color: #333;
        }

        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #666;
        }

        .empty-state h3 {
            margin: 0 0 10px 0;
            color: #333;
        }

        .empty-state p {
            margin: 0 0 20px 0;
        }

        .empty-state a {
            color: #007bff;
            text-decoration: none;
            font-weight: 500;
        }

        .empty-state a:hover {
            text-decoration: underline;
        }

        .loading-indicator {
            text-align: center;
            padding: 20px;
            color: #666;
            display: none;
        }

        .message {
            padding: 12px 16px;
            border-radius: 6px;
            margin-bottom: 20px;
            font-size: 14px;
            font-weight: 500;
        }

        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        /* Confirmation Modal */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.5);
            animation: fadeIn 0.3s ease;
        }

        .modal.show {
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .modal-content {
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            max-width: 400px;
            width: 90%;
            text-align: center;
            animation: slideIn 0.3s ease;
        }

        .modal h3 {
            margin: 0 0 15px 0;
            color: #333;
        }

        .modal p {
            margin: 0 0 25px 0;
            color: #666;
        }

        .modal-actions {
            display: flex;
            gap: 10px;
            justify-content: center;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        @keyframes slideIn {
            from { transform: translateY(-20px); opacity: 0; }
            to { transform: translateY(0); opacity: 1; }
        }

        @media (max-width: 768px) {
            body {
                padding: 15px;
            }

            .post-item {
                flex-direction: column;
                align-items: flex-start;
                gap: 15px;
            }

            .post-actions {
                margin-left: 0;
                width: 100%;
            }

            .btn {
                flex: 1;
            }

            .pagination {
                flex-wrap: wrap;
                gap: 10px;
            }
        }
    </style>
</head>
<body>
<a href="/" class="back-link">← Back to Home</a>

<div class="header">
    <h1>Edit Posts</h1>
    <p>Manage and edit your blog posts</p>
</div>

<div id="messages"></div>

<div id="posts-content">
    <div class="posts-container">
        <div class="posts-header">
            <h2>All Posts ({{.TotalPosts}} total)</h2>
        </div>

        {{if .Posts}}
        <ul class="posts-list">
            {{range .Posts}}
            <li class="post-item">
                <div class="post-info">
                    <h3 class="post-title">{{truncateContent .Title 45}}</h3>
                    <p class="post-date">Created: {{formatDate .CreatedAt}}</p>
                </div>
                <div class="post-actions">
                    <a href="/posts/{{.ID.Hex}}/edit" class="btn btn-edit">Edit</a>
                    <button class="btn btn-delete"
                            onclick="confirmDelete('{{.ID.Hex}}', '{{.Title}}')">
                        Delete
                    </button>
                </div>
            </li>
            {{end}}
        </ul>
        {{else}}
        <div class="empty-state">
            <h3>No posts found</h3>
            <p>You haven't created any posts yet.</p>
            <a href="/posts/create">Create your first post</a>
        </div>
        {{end}}
    </div>

    {{if gt .TotalPages 1}}
    <div class="pagination">
        <button
                hx-get="/posts/edit?page={{.PrevPage}}"
                hx-target="body"
                hx-indicator=".loading-indicator"
                {{if not .HasPrev}}disabled{{end}}>
            Previous
        </button>

        <span class="page-info">Page {{.CurrentPage}} of {{.TotalPages}}</span>

        <button
                hx-get="/posts/edit?page={{.NextPage}}"
                hx-target="body"
                hx-indicator=".loading-indicator"
                {{if not .HasNext}}disabled{{end}}>
            Next
        </button>
    </div>
    {{end}}
</div>

<div class="loading-indicator">
    Loading...
</div>

<!-- Delete Confirmation Modal -->
<div id="deleteModal" class="modal">
    <div class="modal-content">
        <h3>Confirm Delete</h3>
        <p>Are you sure you want to delete "<span id="deletePostTitle"></span>"?</p>
        <p style="color: #dc3545; font-size: 0.9em;">This action cannot be undone.</p>
        <div class="modal-actions">
            <button class="btn btn-delete" id="confirmDeleteBtn">Delete</button>
            <button class="btn" style="background: #6c757d; color: white;" onclick="closeDeleteModal()">Cancel</button>
        </div>
    </div>
</div>

<script>
    let deletePostId = null;

    // Handle HTMX loading states
    document.body.addEventListener('htmx:beforeRequest', function() {
        document.querySelector('.loading-indicator').style.display = 'block';
    });

    document.body.addEventListener('htmx:afterRequest', function() {
        document.querySelector('.loading-indicator').style.display = 'none';
    });

    // Delete confirmation functions
    function confirmDelete(postId, postTitle) {
        deletePostId = postId;
        document.getElementById('deletePostTitle').textContent = postTitle;
        document.getElementById('deleteModal').classList.add('show');
    }

    function closeDeleteModal() {
        document.getElementById('deleteModal').classList.remove('show');
        deletePostId = null;
    }

    // Handle delete confirmation
    document.getElementById('confirmDeleteBtn').addEventListener('click', function() {
        if (!deletePostId) return;

        const btn = this;
        btn.classList.add('loading');
        btn.disabled = true;

        // Send delete request
        fetch(` +
	"`/posts/${deletePostId}`" + `, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            }
        })
            .then(response => {
                if (response.ok) {
                    showMessage('Post deleted successfully', 'success');
                    // Reload the current page
                    window.location.reload();
                } else {
                    throw new Error('Delete failed');
                }
            })
            .catch(error => {
                showMessage('Error deleting post. Please try again.', 'error');
                console.error('Delete error:', error);
            })
            .finally(() => {
                btn.classList.remove('loading');
                btn.disabled = false;
                closeDeleteModal();
            });
    });

    // Close modal when clicking outside
    document.getElementById('deleteModal').addEventListener('click', function(e) {
        if (e.target === this) {
            closeDeleteModal();
        }
    });

    // Close modal with Escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape' && deletePostId) {
            closeDeleteModal();
        }
    });

    function showMessage(text, type) {
        const messagesContainer = document.getElementById('messages');
        messagesContainer.innerHTML = ` + "`" + `<div class="message ${type}">${text}</div>` + "`" + `;

        // Auto-hide success messages
        if (type === 'success') {
            setTimeout(() => {
                const message = messagesContainer.querySelector('.message.success');
                if (message) {
                    message.remove();
                }
            }, 4000);
        }
    }

    // Enhanced button interactions
    document.querySelectorAll('.btn').forEach(btn => {
        btn.addEventListener('mousedown', function() {
            if (!this.disabled) {
                this.style.transform = 'translateY(0) scale(0.98)';
            }
        });

        btn.addEventListener('mouseup', function() {
            if (!this.disabled) {
                this.style.transform = 'translateY(-1px) scale(1)';
            }
        });

        btn.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0) scale(1)';
        });
    });
</script>
</body>
</html>`

type Moderation struct {
	posts      []models.Post
	page       int
	limit      int
	totalPosts int
}

type editPageData struct {
	Posts       []models.Post `json:"posts"`
	CurrentPage int           `json:"current_page"`
	TotalPages  int           `json:"total_pages"`
	TotalPosts  int           `json:"total_posts"`
	HasPrev     bool          `json:"has_prev"`
	HasNext     bool          `json:"has_next"`
	PrevPage    int           `json:"prev_page"`
	NextPage    int           `json:"next_page"`
}

func NewModeration(
	posts []models.Post,
	page, limit int,
	totalPosts int,
) *Moderation {
	return &Moderation{
		posts:      posts,
		page:       page,
		limit:      limit,
		totalPosts: totalPosts,
	}
}

func (m *Moderation) GeneratePage() (string, error) {
	totalPages := int(math.Ceil(float64(m.totalPosts) / float64(m.limit)))

	data := &editPageData{
		Posts:       m.posts,
		CurrentPage: m.page,
		TotalPages:  totalPages,
		TotalPosts:  m.totalPosts,
		HasPrev:     totalPages > 1,
		HasNext:     m.page < totalPages,
		PrevPage:    m.page - 1,
		NextPage:    m.page + 1,
	}

	tmpl, err := template.New("post-moderation").Funcs(funcMap).Parse(moderationHTML)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
