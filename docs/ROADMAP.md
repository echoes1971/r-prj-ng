# R-Project NextGen: Roadmap.


References:
- Link to the old project: https://github.com/echoes1971/r-prj
  See the php subfolder.
- db/00_initial.sql

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
  - all documents that are public access have to be shown to non logged users
  - users can navigate the site plus the content they have access to (personal and of the assigned groups)
  - if the user is logged in and has rights to modify the content is visiting, show a button to bring him to the edit form


Frontend
- handle error messages
  - translation in the 4 languages
- unit tests: React Testing Library

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
