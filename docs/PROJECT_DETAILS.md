# ρBee Project Details

## Architecture Overview

ρBee follows a modular architecture with three main components:

### Backend (Go)
- RESTful API with JWT authentication
- File upload/download with token-based access
- Database layer with MariaDB support
- Permission system (Unix-like rwx------)
- Multi-language content filtering

### Frontend (React)
- Component-based UI with Bootstrap styling
- WYSIWYG editor (ReactQuill) with custom features
- PWA capabilities
- Responsive design for mobile/desktop

### CLI Client (Go)
- Admin commands for user/file management
- API interaction for bulk operations

## Database Schema

Based on a flexible DBObject model:
- Core tables: `rprj_objects`, `rprj_users`, `rprj_groups`
- Content types: `rprj_pages`, `rprj_news`, `rprj_files`, etc.
- Audit fields: creator, last_modify, deleted tracking
- Permissions: owner/group/world with rwx flags

See `db/00_initial.sql` for full schema.

## Roadmap

### MVP Features (High Priority)
- [x] Full-text search with filters
- [ ] OAuth integration (Google, GitHub)
- [ ] Advanced editor features (tables, markdown)
- [ ] Email notifications and password reset

### Completed Features
- Multi-language support
- File embedding with JWT tokens
- Drag & drop file upload
- User/group permissions
- PWA manifest

### Future Plans
- GraphQL API
- Plugin system
- Advanced analytics
- Mobile app

See [ROADMAP.md](ROADMAP.md) for full details.

## API Documentation

### Authentication
- POST `/api/auth/login` - User login
- POST `/api/auth/register` - User registration (planned)

### Objects
- GET `/api/objects` - List objects with filters
- POST `/api/objects` - Create new object
- GET `/api/objects/{id}` - Get object details

### Files
- POST `/api/files/upload` - Upload file
- GET `/api/files/{id}/download` - Download file (with token)

Full OpenAPI/Swagger docs planned.

## Development Setup (Advanced)

### Backend Development
```bash
cd be
go mod tidy
go run main.go
```

### Frontend Development
```bash
cd fe
npm install
npm start
```

### Testing
- Backend: `go test ./...`
- Frontend: `npm test` (planned)

### Docker Development
```bash
docker-compose -f docker-compose.dev.yml up
```

## Deployment

### Production
```bash
docker-compose up -d
```

### Environment Variables
- `DB_HOST`, `DB_USER`, `DB_PASS` - Database config
- `JWT_SECRET` - Authentication secret
- `SMTP_HOST` - Email config (planned)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Submit pull request

See [ROADMAP.md](ROADMAP.md) for priority features.

## License

Apache v2.0 - See LICENSE file.