package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

type CreatePage struct{}

func NewCreatePage() *CreatePage {
	return &CreatePage{}
}

//var createPageTemplate = func() string {
//	bts, err := os.ReadFile("./internal/templates/html/create.html")
//	if err != nil {
//		zap.L().Fatal("error while reading html template", zap.Error(err))
//		return ""
//	}
//	return string(bts)
//}()

const createPageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <meta name="viewport" content="width=device-width, initial-scale=1.0">
   <title>Create New Post</title>
   <script src="https://unpkg.com/htmx.org@1.9.10"></script>
   <style>
       body {
           font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
           max-width: 800px;
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
           margin: 0;
       }

       .form-container {
           background: white;
           border-radius: 12px;
           padding: 30px;
           box-shadow: 0 2px 10px rgba(0,0,0,0.1);
           margin-bottom: 20px;
       }

       .form-group {
           margin-bottom: 25px;
       }

       .form-group label {
           display: block;
           margin-bottom: 8px;
           color: #333;
           font-weight: 500;
           font-size: 14px;
       }

       .form-group input[type="text"],
       .form-group textarea {
           width: -webkit-fill-available;
           padding: 12px 16px;
           border: 1px solid #ddd;
           border-radius: 6px;
           font-size: 16px;
           font-family: inherit;
           transition: border-color 0.2s ease, box-shadow 0.2s ease;
           background: white;
       }

       .form-group input[type="text"]:focus,
       .form-group textarea:focus {
           outline: none;
           border-color: #007bff;
           box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
       }

       .form-group textarea {
           resize: vertical;
           min-height: 120px;
           line-height: 1.5;
       }

       .button-group {
           display: flex;
           gap: 12px;
           justify-content: flex-end;
           margin-top: 30px;
       }

       .btn {
           padding: 12px 24px;
           border: none;
           border-radius: 6px;
           font-size: 14px;
           font-weight: 500;
           cursor: pointer;
           transition: background-color 0.2s ease, transform 0.1s ease;
           text-decoration: none;
           display: inline-flex;
           align-items: center;
           justify-content: center;
       }

       .btn-primary {
           background: #007bff;
           color: white;
       }

       .btn-primary:hover:not(:disabled) {
           background: #0056b3;
       }

       .btn-secondary {
           background: #6c757d;
           color: white;
       }

       .btn-secondary:hover {
           background: #545b62;
       }

       .btn:disabled {
           background: #ddd;
           cursor: not-allowed;
           color: #999;
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
           width: 16px;
           height: 16px;
           margin: -8px 0 0 -8px;
           border: 2px solid transparent;
           border-top: 2px solid white;
           border-radius: 50%;
           animation: spin 1s linear infinite;
       }

       @keyframes spin {
           0% { transform: rotate(0deg); }
           100% { transform: rotate(360deg); }
       }

       .message {
           padding: 12px 16px;
           border-radius: 6px;
           margin-bottom: 20px;
           font-size: 14px;
           font-weight: 500;
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

       .field-error {
           color: #dc3545;
           font-size: 13px;
           margin-top: 5px;
           display: none;
       }

       .form-group.has-error input,
       .form-group.has-error textarea {
           border-color: #dc3545;
       }

       .form-group.has-error .field-error {
           display: block;
       }

       .loading-indicator {
           text-align: center;
           padding: 20px;
           color: #666;
           display: none;
       }

       @media (max-width: 768px) {
           body {
               padding: 15px;
           }

           .form-container {
               padding: 20px;
           }

           .button-group {
               flex-direction: column-reverse;
           }

           .btn {
               width: 100%;
               justify-content: center;
           }
       }
   </style>
</head>
<body>
<a href="/" class="back-link">← Back to Home</a>
<div class="header">
   <h1>Create New Post</h1>
</div>

<div id="messages"></div>

<div class="form-container">
   <form id="post-form"
         hx-post="/posts"
         hx-include="#title, #content"
   >
   <!--          hx-trigger="submit"-->
   <!--          hx-target="#form-response"-->
   <!--    hx-indicator=".loading-indicator"-->

       <div class="form-group" id="title-group">
           <label for="title">Title</label>
           <input type="text"
                  id="title"
                  name="title"
                  placeholder="Enter post title..."
                  required>
           <div class="field-error" id="title-error"></div>
       </div>

       <div class="form-group" id="content-group">
           <label for="content">Content</label>
           <textarea id="content"
                     name="content"
                     placeholder="Write your post content here..."
                     required></textarea>
           <div class="field-error" id="content-error"></div>
       </div>

       <div class="button-group">
           <a href="/home" class="btn btn-secondary">Main Menu</a>
           <button type="submit" id="submit-btn" class="btn btn-primary">
               Create Post
           </button>
       </div>
   </form>

   <div id="form-response"></div>
</div>
</body>
</html>
`

func (c *CreatePage) GeneratePage() (string, error) {
	tmpl, err := template.New("posts").Funcs(funcMap).Parse(createPageTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
