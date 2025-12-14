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
- [ ] manifest.json to make it recognisable as a PWA


### Search & Discovery (NEXT - MVP BLOCKER)
- [ ] Full-text search in HTML content // 👤 Roberto: A search box in the NavBar that leads to a /nav/search with results and filters
  - [x] Anonymous user search (public content only)
  - [x] Logged user search (public + accessible content)
- [ ] Advanced filters
  - [ ] Date range filter
  - [ ] File type filter
  - [ ] Author filter
  - [ ] Language filter
- [ ] Search results highlighting
- [ ] Search in name, description, and HTML content // ⚠️ all objects have name and description, only page and news have html
- [ ] Pagination for search results

### Documentation
- [ ] Easy steps to get the CMS up and running on your machine
- [ ] Project description: brief description of the project, some how-to, feature list, history?
- [ ] License: gpl? lgpl? apache 2.0?

### Rich Text Editor Improvements (HIGH PRIORITY)
- [ ] Pre condition: make it a separate reusable component
- [ ] Tables support in ReactQuill
- [ ] Image resize (simply put width="...")
- [ ] Text and elements alignment
- [ ] Markdown alternative editor // ❓ where do we store the markdown, in the html field or another? if it's the same field, how do we distinguish the 2 in View and Edit?
  - [ ] Toggle between WYSIWYG and Markdown // 👤 Roberto: I love Markdown, I don't know why :)
  - [ ] Markdown preview
- [ ] Code syntax highlighting // I'd say "nice to have" (it will give a professional feeling to the end user), by now it seems redundant as we have a wysiwig editor
- [ ] Custom CSS classes selector // 👤 Roberto: YES ! Let's customize the site colors at first (should be easy with bootstrap primary etc.). I'd like to have selectable skins for the public site, but that looks too much for now?
- [ ] Emoji picker // 👤 Roberto: YES !

### OAuth Integration (HIGH PRIORITY - if not complicated)
- [ ] Google OAuth login // 👤 Roberto: YES !
- [ ] GitHub OAuth login
- [ ] Facebook OAuth login (optional)
- [ ] Link existing account with OAuth
- [ ] OAuth user creation with:  // 👤 Roberto: YES !
  - [ ] default permissions rwx------
  - [ ] it will have its private group and then will be linked to the group "Guest" that has ID "-4" (use the type UserGroup)

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

---

## 📋 TODO - Organized by Priority

### Security & Performance
- [ ] Token auto-refresh (file preview tokens expire after 15 min)
- [ ] Session storage/caching for tokens
- [ ] Rate limiting for file download
- [ ] Rate limiting for API endpoints (general)
- [ ] CSRF protection
- [ ] Content Security Policy headers
- [x] Password encryption with salt (already in plan, not implemented) // 👤 Roberto: it's done
- [ ] Rainbow table attack protection

### CMS Features
- [ ] DBLink implementation (mentioned but not implemented) // 👤 Roberto: I'm not sure if I want to port it from the old project
- [ ] DBEvent implementation (mentioned but not implemented) // 👤 Roberto: it's basically a calendar entry, with field about recurring events too
- [ ] Versioning/History for DBPage (track who modified what when)
- [ ] Draft system for content (save without publishing)
- [ ] Content scheduling (publish at specific date/time) // 👤 Roberto: simple fields with publish date start and publush date end?
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

### File Management
- [ ] File upload progress indicator
- [ ] Image resizing/thumbnails on upload (backend exists, integrate in UI) // 👤 Roberto: we have already thumbnails in the BE. Are not yet used in the FE
- [ ] File storage optimization (nested directory structure: `files/XX/YY/ZZZZ...`) // 👤 Roberto: now the structure is <father_id>/<file>
- [ ] Quota management per user/group
- [ ] File versioning
- [ ] Batch file upload (multiple files at once) // 👤 Roberto: yes
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
- [ ] Support for PostgreSQL and SQLite3

### Developer Experience
- [ ] API documentation improvements
- [ ] GraphQL endpoint (alternative to REST)? // 👤 Roberto: interesting, I need to learn about this new (for me) tool
- [ ] Webhook system for events (onCreate, onUpdate, onDelete)
- [ ] Plugin/extension system // 👤 Roberto: "nice to have" how can we make the project extendable, both in BE and in FE?
- [ ] CLI tools for admin tasks
- [ ] Docker compose for development // 👤 Roberto: ongoing?
- [ ] Hot reload for backend (air or similar) // 👤 Roberto: is active, check the .dev compose file

### Nice to Have
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
- [ ] RSS/Atom feeds for content // 👤 Roberto: YES! it should be easy to implement
- [ ] Sitemap generation (XML for SEO)
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

---

## OLD NOTES (To be organized/removed)

## TODO

FE
- ~~edit index page of a folder~~ ✅ DONE




## Features

Features
- Framework
  - [DBObject](#the-dbobject)
    - is a base abstract class for almost all other etities in the project
    - how do you handle subclasses in GO that share or override methods from the superclass ?
    - frontend
      - a React reusable component can be a good idea in the frontend
    - backend
      - reusable or copy paste handlers and db logic
- Contacts
  - Country
  - Company
  - Person
- CMS
  - DBFolder
  - DBNote
  - DBPage
  - DBFile
    - Upload and download
    - binary or image: content type
  - DBNews (=to page)
  - DBLink
  - DBEvent
- User
  - each time a user select a language
    - IF the user is logged in, it is saved in the users table
    - IF the user is not logged in, as soon as he logs in, it is saved in DB
    - a field must be added to store the language (file db/02_users_language.sql? is it automatically executed on an existing DB?)
  - the user has a profile section in the dropdown under his name
  - will use both users and dbperson. See the old project
  - a user photo ?
  - all subclasses of dbobject accessible from the profile page?
- Site navigation
  - the site is organized in folders and content pages, like a tree
  - a read-only navigation has to be implemented
  - all documents are returned to the frontend with the classname in the metadata, so the UI can select the appropriate view
  - documents (dbobjects) have to be accessible via url: http://sitename.com/<object id in the format of xxxx-xxxxxxxx-xxxx>
  - all documents that are public access have to be shown to non logged users
  - the site root is a folder specified in an environment parameter passed to the react application or, given the app name, the react app can search the backend for a folder with name <app name>
  - a folder has a default main page (a child) with simply the name "index"
  - a tree view is diplayed on the side with expandable nodes, to access the site structure and sub-nodes
    - the tree view will disappear on mobile?
  - a breadcrumb path is displayed on the top, use the father_id to display up to the root folder
  - a root folder is a folder without a father
  - a search field must be displayed on the top (at the left of the login button?): it allows to research objects by name, description, content (html or other)
  - in the same folder, a page can have the same name of another only if the specified language is different.
  - users can navigate the site plus the content they have access to (personal and of the assigned groups)
  - if the user is logged in and has rights to modify the content is visiting, show a button to bring him to the edit form


Frontend
- refine the UI for all data types (objects)
- handle error messages
  - translation in the 4 languages
- unit tests: React Testing Library
- registration process
  - a non logged user can register a new account
  - a new user has only his own private group and can add objects only visible to himself rwx------
  - send confirmation email to activate the user: how?
  - add a field user_enabled to the table?

Backend
- add swagger
- db transactionality when writing
- error handling
  - log
  - pass error to the UI
- NO, NOT NOW: ridondare il controllo che solo gli admins possono fare CRUD su utenti e gruppi, ad eccezione che un utente puo modificare il suo utente ed il suo gruppo primario solamente.
- unit tests
  - db
  - handlers
- logging
  - any suggestion on the strategy: besides using log.print (the how) what info should be logged (the what, more important)?
- rate limiting: suggestions?
- password encryption:
  - strategy: IF salt is empty => password not encrypted ELSE password is encrypted, this ensure compatibility until the adoption of the new model
  - any user created or modified with this tool MUST have the password encrypted
  - beware of rainbow tables attacks
- pagination
- support for postgresql and sqlite3


## The framework: a detailed analisys of the backend

The framework is based on the main abstract class DBObject, and its specializations (ie. DBPerson, DBNote, etc.)

The purpose is to create a light ORM.

Each operation, use the user_id and the group_id list extracted from the JWT (inserted there by the login, we trust the JWT to not be manipulated, right?)

There are 3 main parts in the framework, in the backend:
1. the DBEntity
2. the basic CRUD (not so basic as you will see later)
3. the search engine
4. the DBObject itself

### 1. The DBEntity

Abstract class: it represents a generic table.

Has methods (to be overridden) that can perform actions on the DB:
- before_insert
- after_insert
- before_update
- after_update
- before_delete
- after_delete
and other minor but useful methods like:
- readFKFrom which populates the foreign keys of the DBEntity with value from an instance of another entity
- writeFKTo: viceversa

Example of readFKFrom:
```php
dbuser = User(ID:xxx, login:"myname")
dbgroup = Group(ID:YYY, Name: "My Group")
dbuser.readFKFrom(dbgroup)

Now dbuser = User(ID:xxx, login:"myname", GroupID: "YYY")
```

### 2. The basic CRUD

It's implemented as follows.

insert(dbe):
1. execute dbe.before_insert on the db
2. execute the INSERT INTO the db
3. execute dbe.after_insert

It is transactional. Example:
```GO
tx.Begin()
dbe.before_insert(tx)
tx.Exec("INSERT INTO ...")
dbe.after_insert(tx)
tx.Commit()
```


The creation of the 'INSERT INTO' string is dynamic: it takes the populated values in the dbe and build the SQL only for those.

Example:

person := DBPerson{FirstName: "Mario"}

becomes: INSERT INTO RPRJ_PERSON (firstname) VALUES ("Mario")

Same for update and delete.

### 3. The search engine

How is performed a search? Simple (quite so):

1. Create an instance of the DBE you want to search, i.e. DBEUser
2. populate the fields you want to search
3. build the query based on the populated fields:
   1. if they are basic type like number or date, use `<colonna> = <value>`
   2. if it is a string, use `<colonna> like '%<value>%`, unless the user request an exact match

Permissions do not really exist at this level. They'll be handled with DBObject classes.

### 4. The DBObject

Is subclass of DBEntity.


DBEntity and DBObject are abstract:
- dbentity has subclasses like dbuser with different attributes each time
- dbobject has subclasses with those fields in common, plus other fields dependent on the subclass, such as first and last names for a person.
  - in the db, the common fields are replicated for each table (we are not in postgresql :( )

The DBObject has fields:
- id: 'uuid','not null'
- owner: 'uuid','not null'
- group_id: 'uuid','not null'
- permissions: 'char(9)',"not null default 'rwx------'" // user,group,all
- creator: 'uuid','not null'
- creation_date: 'datetime','default null'
- last_modify: 'uuid','not null'
- last_modify_date: 'datetime','default null'
- deleted_by: 'uuid','default null'
- deleted_date: 'datetime',"not null default '0000-00-00 00:00:00'"
- father_id: 'uuid','default null'
- name: 'varchar(255)','not null'
- description: 'text','default null'

Redefines the methods:
- before_insert
- before_update
- before_delete
here it checks permissions before executing the insert, update and delete as above (super.insert, etc.)

Has methods to be overridden:
- canRead(U/G/A), canWrite, canRead: they check the requested permission for User, Group, All; returns true if undefined.
- setDefaultValues

Permissions.

It's a char(9) field i.e. rwxrw-r-- . you read it exactly as a unix permission: the user can do all (rwx), the group can read and write (rw-), and everybody can read it (r--).


The CRUD and the search are handled as above for the DBEntity. The only exception here is that there will be permissions checks.

The search:
- IF the search object is a subclass of DBObject -> search as above, only in the selected table
- IF the search is of class DBObject itself:
  - select the DBObject fields only on each table of the subclasses.
  - returns the union of the results
- return only the results the user has right to read: as a user OR has the group OR the object is public

Clarification for the search.

In the old framework, I also included the type name in the union:
- I could either settle for the common data (name, description, etc.) and display it on screen.
- Or, for each object, I required creating an instance of the specific type (e.g., DBPerson) with ALL the fields populated (alas, triggering a query for each result).
Why? In most cases, I only needed the common data to display, and when the user selected one of these objects, I could drill down and retrieve all the data.
