# Contributing to ρBee

Thank you for your interest in contributing to ρBee! This document outlines how to contribute to the project.

## Code of Conduct
- Be respectful and inclusive
- Focus on constructive feedback
- Follow the roadmap priorities

## How to Contribute

### Reporting Issues
- Use GitHub Issues for bugs and feature requests
- Include steps to reproduce, expected vs actual behavior
- Tag with appropriate labels (bug, enhancement, etc.)

### Development Setup
See [README.md](../README.md) for basic setup, and [PROJECT_DETAILS.md](PROJECT_DETAILS.md) for advanced configuration.

### Pull Requests
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Add tests if applicable
5. Ensure code builds: `npm run build` (frontend), `go build` (backend)
6. Submit PR with clear description

### Code Style
- **Go**: Follow standard Go formatting (`gofmt`)
- **React**: Use ESLint rules, functional components preferred
- **Commits**: Use conventional commits (feat:, fix:, docs:, etc.)

### Testing
- Backend: Add unit tests for new functions
- Frontend: Add React Testing Library tests for components
- Run tests before submitting PR

## Areas for Contribution
See [ROADMAP.md](ROADMAP.md) for current priorities:
- OAuth integration
- Email system
- Advanced search filters
- Mobile responsiveness
- Unit test coverage

## Questions?
Open a GitHub Discussion or contact the maintainers.