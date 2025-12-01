# File Embedding in DBPages

## Overview

This system allows embedding DBFile objects (images, documents) inside DBPage HTML content with automatic permission handling and temporary JWT tokens.

## How It Works

### For Public Files
Files with public read permissions (e.g., `rwxr-xr--` where "others" can read) work directly without tokens.

### For Private Files
Files restricted to specific users/groups require temporary JWT tokens (valid 15 minutes) that are:
- Generated automatically when viewing/editing pages
- Never stored in the database
- User-specific (each user gets their own tokens based on their permissions)

## Usage

### In WYSIWYG Editor
Simply insert images or links as usual. The editor will display them correctly with automatic token injection.

### In HTML Source
Use the `data-dbfile-id` attribute to mark files:

```html
<!-- Image -->
<img src="/api/files/abc123/download" 
     data-dbfile-id="abc123" 
     alt="My Image" />

<!-- Download Link -->
<a href="/api/files/def456/download" 
   data-dbfile-id="def456">
   Download Document
</a>
```

### Important Notes
- **Always use `/api/files/{ID}/download`** as the base URL
- **Always include `data-dbfile-id="{ID}"`** attribute
- Do NOT add `?token=...` manually (it's handled automatically)

## Workflow

### When Editing a Page:
1. System extracts all `data-dbfile-id` from HTML
2. Requests tokens for all files (only those you can access)
3. Injects tokens into URLs for WYSIWYG preview
4. Before saving, removes all tokens (keeps only base URL + data-dbfile-id)

### When Viewing a Page:
1. System extracts all `data-dbfile-id` from HTML
2. Requests tokens for all files (based on current user's permissions)
3. Injects tokens into URLs for display
4. Files refresh automatically with new tokens (15 min expiry)

## Security Features

- ✅ Tokens are user-specific and include user_id in JWT payload
- ✅ Backend validates permissions before generating tokens
- ✅ Tokens expire after 15 minutes
- ✅ No tokens stored in database (always generated on-demand)
- ✅ Users without permission get no token (file won't load)
- ✅ Public files work without tokens (fallback to permission check)

## API Endpoints

### Generate Tokens (Protected)
```
POST /api/files/preview-tokens
Authorization: Bearer {JWT}
Content-Type: application/json

{
  "file_ids": ["abc123", "def456", ...]
}

Response:
{
  "success": true,
  "tokens": {
    "abc123": "eyJhbGc...",
    "def456": "eyJhbGc..."
  }
}
```

### Download File
```
GET /api/files/{id}/download
GET /api/files/{id}/download?token={JWT}

- If token present: validates JWT and serves file
- If no token: checks user permissions and serves if authorized
```

## Example HTML

```html
<h2>My Page with Embedded Files</h2>

<p>Here's a public image:</p>
<img src="/api/files/public123/download" 
     data-dbfile-id="public123" 
     alt="Public Image" />

<p>And here's a private document:</p>
<p>
  <a href="/api/files/private456/download" 
     data-dbfile-id="private456">
    Download Confidential Report (PDF)
  </a>
</p>

<p>Multiple images:</p>
<div class="gallery">
  <img src="/api/files/img1/download" data-dbfile-id="img1" />
  <img src="/api/files/img2/download" data-dbfile-id="img2" />
  <img src="/api/files/img3/download" data-dbfile-id="img3" />
</div>
```

## Implementation Details

### Frontend (React)
- `extractFileIDs(html)`: Parses HTML and extracts all data-dbfile-id values
- `requestFileTokens(fileIDs)`: Calls API to get tokens for multiple files
- `injectTokensForEditing(html, tokens)`: Adds tokens to URLs for WYSIWYG
- `injectTokensForViewing(html, tokens)`: Adds tokens to URLs for display
- `cleanTokensBeforeSave(html)`: Removes tokens before saving to DB

### Backend (Go)
- `GenerateFileTokensHandler`: Validates permissions and generates JWT tokens
- `DownloadFileHandler`: Serves files with token validation or permission check
- JWT payload: `{id: fileID, user_id: userID, exp: timestamp}`

## Token Refresh

Tokens expire after 15 minutes. To refresh:
- **In edit mode**: Close and reopen the page
- **In view mode**: Page reload will request new tokens automatically
- Future enhancement: Auto-refresh tokens before expiry

## Troubleshooting

### Image Not Loading
1. Check if file exists and user has read permission
2. Open browser console to see token request errors
3. Verify `data-dbfile-id` matches actual file ID
4. Check if token expired (15 min limit)

### Token Errors
- "Invalid or expired token": Token expired, reload page
- "Token does not match file ID": data-dbfile-id mismatch
- "Unauthorized": User has no read permission on file

## Future Enhancements

- [ ] Auto-refresh tokens before expiry (using timer in React)
- [ ] Batch token requests for large pages
- [ ] Cache tokens in sessionStorage
- [ ] Visual indicator for private vs public files in editor
- [ ] File picker UI integration with data-dbfile-id insertion
