
# ρBee (rhobee) - Lightweight Headless CMS Framework

[![Headless CMS](https://img.shields.io/badge/Type-Headless%20CMS-blue)](https://github.com/your-repo/r-prj-ng)
[![Go](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org/)
[![React](https://img.shields.io/badge/React-19+-blue)](https://reactjs.org/)
[![License](https://img.shields.io/github/license/echoes1971/rhobee)](LICENSE)

**ρBee** is a lightweight, flexible Headless Content Management System (CMS) built with a modular DBObject architecture. It supports multi-language content, user permissions, file management, and more – ideal for small to medium websites or as a foundation for custom web applications.

## Features
- **Multi-Language Support**: EN, IT, DE, FR with language-based content filtering.
- **User Management**: JWT authentication, groups, Unix-like permissions (rwx------).
- **Rich Text Editor**: WYSIWYG with image resizing, emoji picker, and file embedding.
- **File Management**: Upload/download with drag & drop, token-based access, MIME filtering.
- **Search & Discovery**: Full-text search with filters (author, date, language).
- **PWA Ready**: Installable as a Progressive Web App.
- **API-Driven**: RESTful backend in Go, React frontend, CLI client in Go.

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/r-prj-ng.git
   cd r-prj-ng
   ```

2. Start the services:
   ```bash
   docker-compose up -d
   ```
   This launches the backend (Go), frontend (React), database (MariaDB), and proxy (Nginx).

3. Access the application:
   - Frontend: http://localhost:8080
   - Backend API: http://localhost:8080/api/
   - Admin CLI: Use `./client` commands (build from `client/` if needed)
   - Authenticate with the user 'adm' and password 'mysecretpass': change it at the first login.

<!-- 4. Create an admin user:
   - Via CLI: `./client login` (or build and run from `client/` directory)
   - Or use API endpoints for registration. -->

### Development Mode
For hot-reload during development:
- Backend: Use `docker-compose.dev.yml` with Air for Go hot-reload.
- Frontend: Use `docker-compose.dev.yml` for React hot-reload.

### Troubleshooting
- **Ports occupied?** Edit `.env` to change ports.
- **DB issues?** Check logs: `docker logs r-prj-ng_db`
- **Build errors?** Ensure Docker is running and ports are free.

## Project Structure
- `be/`: Go backend (REST API, file handling)
- `fe/`: React frontend (UI components)
- `client/`: Go CLI for admin tasks
- `db/`: SQL schemas and migrations
- `docs/`: Detailed documentation
- `proxy/`: Nginx configuration

## History
Evolved from the original PHP-based r-prj project (see [old repo](https://github.com/echoes1971/r-prj)). As early as 2005, the project explored API-driven decoupling with XML-RPC and JSON/REST to separate content engine from presentation, including multi-platform clients in C++/Qt and Python/Qt. Attempts with Bootstrap and AJAX for modern frontends led to the NextGen rewrite in Go and React for enhanced performance, scalability, and maintainability.

## Contributing
See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

## License
Apache v2.0 (see LICENSE file).

---

For detailed documentation, roadmap, and advanced setup, see [docs/PROJECT_DETAILS.md](docs/PROJECT_DETAILS.md).




