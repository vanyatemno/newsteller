package templates

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"html/template"
	"newsteller/internal/models"
)

type SingleTemplate struct {
	post *models.Post
}

func (s *SingleTemplate) GeneratePage() (string, error) {
	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		zap.L().Error("Error parsing template:", zap.Error(err))
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s.post); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

const postTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 15px;
            line-height: 1.6;
        }
        .post {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 15px;
            margin: 15px 0;
            background: #f9f9f9;
        }
        .post-header {
            border-bottom: 1px solid #eee;
            padding-bottom: 10px;
            margin-bottom: 15px;
        }
        .post-title {
            color: #333;
            margin: 0 0 10px 0;
            font-size: 1.4em;
            word-wrap: break-word;
        }
        .post-meta {
            color: #666;
            font-size: 0.85em;
            word-wrap: break-word;
        }
        .post-content {
            color: #444;
            margin: 15px 0;
            word-wrap: break-word;
            overflow-wrap: break-word;
        }
        .post-actions {
            margin-top: 20px;
            padding-top: 15px;
            border-top: 1px solid #eee;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            cursor: pointer;
            margin-right: 10px;
            font-size: 16px;
            min-height: 44px;
            touch-action: manipulation;
        }
        button:hover {
            background: #0056b3;
        }
        .btn-secondary {
            background: #6c757d;
        }
        .btn-secondary:hover {
            background: #545b62;
        }
        .htmx-indicator {
            opacity: 0;
            transition: opacity 500ms ease-in;
        }
        .htmx-request .htmx-indicator {
            opacity: 1;
        }

        /* Mobile-specific styles */
        @media (max-width: 768px) {
            body {
                padding: 10px;
            }
            .post {
                padding: 12px;
                margin: 10px 0;
                border-radius: 6px;
            }
            .post-title {
                font-size: 1.3em;
                line-height: 1.3;
            }
            .post-meta {
                font-size: 0.8em;
            }
            .post-content {
                font-size: 1em;
                line-height: 1.5;
            }
            button {
                width: 100%;
                margin-right: 0;
                margin-bottom: 10px;
                padding: 14px 20px;
                font-size: 16px;
            }
            .post-actions {
                margin-top: 15px;
                padding-top: 12px;
            }
        }

        @media (max-width: 480px) {
            body {
                padding: 8px;
            }
            .post {
                padding: 10px;
                margin: 8px 0;
            }
            .post-title {
                font-size: 1.2em;
            }
            .post-meta {
                font-size: 0.75em;
            }
            .post-content {
                font-size: 0.95em;
            }
        }
    </style>
</head>
<body>
<div id="post-container">
    <article class="post">
        <header class="post-header">
            <h1 class="post-title">{{.Title}}</h1>
            <div class="post-meta">
                <span>Created: {{.CreatedAt.Format "January 2, 2006 at 3:04 PM"}}</span>
                {{if not .UpdatedAt.IsZero}}
                <span> | Updated: {{.UpdatedAt.Format "January 2, 2006 at 3:04 PM"}}</span>
                {{end}}
            </div>
        </header>

        <div class="post-content">
            <p>{{.Content}}</p>
        </div>

        <footer class="post-actions">
            <button class="btn-secondary" onclick="window.location.href='/posts/search'">Back to posts</button>
            <span id="loading" class="htmx-indicator">Loading...</span>
        </footer>
    </article>
</div>
</body>
</html>
`

func RenderSinglePost(post *models.Post) (string, error) {
	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		zap.L().Error("Error parsing template:", zap.Error(err))
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
