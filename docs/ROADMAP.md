# R-Project NextGen: Roadmap

## Project Name: ρBee (rhobee)
A lightweight CMS framework based on a flexible DBObject architecture.

References:
- Link to the old project: https://github.com/echoes1971/r-prj
  See the php subfolder.
- db/00_initial.sql

---

## 🎯 MVP - Priority Features

- [ ] ObjectList: as Admin, I have a filter to see all deleted objects and massively select the delete
- [x] manifest.json to make it recognisable as a PWA


### Search & Discovery (NEXT - MVP BLOCKER)
- [ ] Advanced filters
  - [ ] Date range filter // Roberto: a generic range can be implemented, passing [_from_<name attribute>, _to_<name attribute>] in the metadata. These will be handled by SearchObjectsHandler that passes them to DBRepository.Search in the metadata of the search object
  - [ ] File type filter
  - [x] Author filter
  - [x] Language filter
- [ ] Search results highlighting
- [ ] Search in name, description, and HTML content // ⚠️ all objects have name and description, only page and news have html
- [ ] Pagination for search results

### Documentation
- [ ] Easy steps to get the CMS up and running on your machine
- [ ] Project description: brief description of the project, some how-to, feature list, history?
- [x] License: gpl? lgpl? apache 2.0?


### OAuth Integration (HIGH PRIORITY - if not complicated)
- [ ] Google OAuth login // 👤 Roberto: YES !
- [ ] OAuth user creation with:  // 👤 Roberto: YES !
  - [ ] default permissions rwx------
  - [ ] it will have its private group and then will be linked to the group "Guest" that has ID "-4" (use the type UserGroup)
- [ ] Link existing account with OAuth
- [ ] GitHub OAuth login
- [ ] Facebook OAuth login (optional)

---

## ✅ Completed Features

### CMS Core
- [x] DBFolder with index pages (multi-language)
- [x] DBFolder children sort order (drag & drop)
- [x] DBPage with WYSIWYG editor
- [x] DBNote
- [x] DBFile upload/download
- [x] DBPerson
- [x] DBCompany
- [x] File embedding system with JWT tokens
- [x] FileSelector with write permission filtering
- [x] Multi-language support (EN, IT, DE, FR)

### Search & Discovery
- [x] Full-text search in HTML content // 👤 Roberto: A search box in the NavBar that leads to a /nav/search with results and filters
  - [x] Anonymous user search (public content only)
  - [x] Logged user search (public + accessible content)

### Rich Text Editor Improvements
- [x] Pre condition: make it a separate reusable component
- [x] Text and elements alignment
- [x] Image resize
- [x] Emoji picker

### Site Navigation
- [x] Tree view sidebar with expandable nodes
- [x] Breadcrumb navigation
- [x] Language-based content filtering
- [x] Public/private content access control
- [x] URL-based object access (/<object-id>)
- [x] Edit button for authorized users

### User Management
- [x] Login/logout
- [x] JWT authentication
- [x] User profile page
- [x] Group management
- [x] Permissions system (rwx-------)
- [x] Language preference per user

### File System
- [x] File upload with drag & drop
- [x] Image preview
- [x] File download with token-based auth
- [x] Write permission check for embedding
- [x] MIME type filtering
- [x] Alternative link support

### Security & Performance
- [x] Password encryption with salt (already in plan, not implemented) // 👤 Roberto: it's done


---

## 📋 TODO - Organized by Priority

### Security & Performance
- [ ] Token auto-refresh (file preview tokens expire after 15 min)
- [ ] Session storage/caching for tokens
- [ ] Rate limiting for file download
- [ ] Rate limiting for API endpoints (general)
- [ ] CSRF protection
- [ ] Content Security Policy headers
- [ ] Rainbow table attack protection

### CMS Features
- [x] DBLink implementation
- [ ] DBEvent implementation
  - [x] basic implementation done
  - [ ] advanced: calendar view, compute recurring events
- [ ] Versioning/History for DBPage (track who modified what when) // 👤 Roberto: how?
- [ ] Draft system for content (save without publishing)
- [ ] Content scheduling (publish at specific date/time) // 👤 Roberto: simple fields with publish date start and publish date end?
- [ ] Bulk operations // 👤 Roberto: yes
  - [ ] Delete multiple objects
  - [ ] Move multiple objects
  - [ ] Change permissions for multiple
- [ ] Content duplication/cloning
- [ ] Recently viewed/edited list
- [ ] Favorites/bookmarks system // 👤 Roberto: nice to have, but requires db modifications
- [ ] Tags system for better categorization
- [ ] Content templates
- [ ] Ollama integration? For assisted document redacting or automatic translation? llama3.2 seems light and efficient enough. Open to suggestions

### OAuth Integration (HIGH PRIORITY - if not complicated)


### Rich Text Editor Improvements
- [ ] Custom CSS classes selector // 👤 Roberto: YES ! Let's customize the site colors at first (should be easy with bootstrap primary etc.). I'd like to have selectable skins for the public site, but that looks too much for now?
- [ ] Markdown alternative editor // ❓ where do we store the markdown, in the html field or another? if it's the same field, how do we distinguish the 2 in View and Edit?
  - [ ] Toggle between WYSIWYG and Markdown // 👤 Roberto: I love Markdown, I don't know why :)
  - [ ] Markdown preview
- [ ] Code syntax highlighting // I'd say "nice to have" (it will give a professional feeling to the end user), by now it seems redundant as we have a wysiwig editor
- [ ] ~~Tables support in ReactQuill~~ quill-table not compatible with the current quill version

### File Management
- [x] File upload progress indicator
- [x] Batch file upload (multiple files at once) // 👤 Roberto: yes
- [ ] Image resizing/thumbnails on upload (backend exists, integrate in UI) // 👤 Roberto: we have already thumbnails
- [ ] File storage optimization (nested directory structure: `files/XX/YY/ZZZZ...`) // 👤 Roberto: now the structure is <father_id>/<file>
- [ ] Quota management per user/group
- [ ] File versioning // 👤 Roberto: how?
- [ ] Preview for more file types (PDF viewer, video player) // 👤 Roberto: yes! how?
- [ ] Image editing tools (crop, rotate, filters) // 👤 Roberto: if easy. "nice to have"

### User Experience
- [ ] Mobile responsive improvements // 👤 Roberto: it doesn't look so bad now in mobile, does it?
- [ ] Dark mode polish
- [ ] Accessibility improvements // 👤 Roberto: "nice to have"
  - [ ] ARIA labels
  - [ ] Keyboard navigation
  - [ ] Screen reader support
- [ ] Undo/Redo system for editors
- [ ] Auto-save drafts (local storage) // 👤 Roberto: why not :)
- [ ] Copy/paste improvements in editor
- [ ] Drag & drop file insertion in editor // 👤 Roberto: yes

### Administration
- [ ] Admin dashboard with statistics
  - [ ] User activity // 👤 Roberto: how?
  - [ ] Content statistics
  - [ ] Storage usage // 👤 Roberto: should be easy
  - [ ] Popular pages // 👤 Roberto: needs db support
- [ ] Audit log (comprehensive who/what/when tracking) // 👤 Roberto: not easy
- [ ] User activity monitoring // 👤 Roberto: not easy / how?
- [ ] Backup/restore functionality // 👤 Roberto: mariadb dump/restore or something smarter?
- [ ] Database migrations management
- [ ] System health check endpoint
- [ ] Email configuration for notifications // 👤 Roberto: hah! I think a local smtp server will not work, right?

### Frontend
- [ ] Handle error messages refinement
- [ ] Error translation in 4 languages (partial)
- [ ] Unit tests with React Testing Library // 👤 Roberto: I need at least one as example, so I can work on it in my spare time
- [ ] Registration process // 👤 Roberto: YES! See also OAuth
  - [ ] Non-logged user can register
  - [ ] Email confirmation to activate account
  - [ ] Add user_enabled field to table
  - [ ] New users start with private group only (rwx------)

### Backend
- [ ] Add Swagger/OpenAPI documentation // 👤 Roberto: if easy, I'd say to put it in place ASAP
- [ ] Database transactionality for writes // 👤 Roberto: we have it, haven't we?
- [ ] Transaction isolation level configuration // 👤 Roberto: we have it, haven't we?
- [ ] Error handling improvements
  - [ ] Structured logging
  - [ ] Error messages to UI
- [ ] Logging strategy
  - [ ] What to log (access, errors, changes, etc.)
  - [ ] Log rotation
  - [ ] Log levels (debug, info, warn, error)
- [ ] Unit tests
  - [ ] Database layer tests
  - [ ] Handler tests
  - [ ] Permission tests
- [ ] Pagination for large result sets
- [ ] DB: Add indexes for name, description and html content to support text search
- [ ] Support for PostgreSQL and SQLite3

### Developer Experience
- [ ] API documentation improvements
- [ ] GraphQL endpoint (alternative to REST)? // 👤 Roberto: interesting, I need to learn about this new (for me) tool
- [ ] Webhook system for events (onCreate, onUpdate, onDelete)
- [ ] Plugin/extension system // 👤 Roberto: "nice to have" how can we make the project extendable, both in BE and in FE?
- [x] CLI tools for admin tasks
- [x] Docker compose for development // 👤 Roberto: ongoing?
- [x] Hot reload for backend (air or similar) // 👤 Roberto: is active, check the .dev compose file

### Nice to Have
- [ ] RSS/Atom feeds for content // 👤 Roberto: YES! it should be easy to implement
- [ ] Sitemap generation (XML for SEO)
- [ ] Comments system for pages
- [ ] Sharing links with expiry date // 👤 Roberto: nice
- [ ] Email notifications // 👤 Roberto: I fear we need a provider
  - [ ] Content published
  - [ ] User mentioned
  - [ ] Permission granted
- [ ] Two-factor authentication (2FA)
- [ ] ~~OAuth providers (Google, GitHub, etc.)~~ (moved to MVP priorities)
- [ ] Export content (PDF, ZIP, JSON) // 👤 Roberto: YES!
  - [ ] Export single pages, multiple pages in a ZIP file
  - [ ] Export a selected folder with all its subelements in a zip file that the user can download
- [ ] Import content (from WordPress, other CMS)
- [ ] OpenGraph meta tags for social sharing
- [ ] Print-friendly page styles // 👤 Roberto: YES !

---

## 🔧 Technical Debt & Known Issues

### File Storage
- Current: Flat directory structure per object ID
- Issue: Scalability concerns with many files
- Solution: Nested structure (files/XX/YY/ZZZZ...)

### Database
- Transaction isolation levels need definition
- Add indexes for performance (especially search)
- Consider migration strategy for future schema changes

### Error Handling
- Standardize error responses across API
- Better error messages for users
- Error tracking/monitoring system

### Code Organization
- Extract token management into custom React hook
- Separate API client from components
- Backend: Consider hexagonal architecture
