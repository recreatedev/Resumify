# Resumify

A modern resume builder application that allows users to create, customize, and manage professional resumes with ease. Built with a Go backend and TypeScript frontend, featuring real-time editing, multiple themes, and PDF export capabilities.

## Features

- **Resume Management**: Create, edit, and organize multiple resumes
- **Real-time Editing**: Live preview with instant updates
- **Multiple Themes**: Professional, modern, classic, and creative templates
- **Section Control**: Customize which sections to include and their order
- **Rich Content**: Support for education, experience, projects, skills, and certifications
- **PDF Export**: Generate professional PDF resumes
- **User Authentication**: Secure user management with Clerk
- **Responsive Design**: Works seamlessly on desktop and mobile
- **Monorepo Structure**: Organized with Turborepo for efficient development

## Project Structure

```
resumify/
├── apps/
│   ├── backend/          # Go backend API
│   └── frontend/         # React frontend application
├── packages/             # Shared packages
│   ├── emails/          # Email templates
│   ├── openapi/         # API contracts
│   └── zod/             # Validation schemas
├── package.json         # Monorepo configuration
├── turbo.json           # Turborepo configuration
└── README.md            # This file
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Node.js 22+ and Bun
- PostgreSQL 16+
- Redis 8+

### Installation

1. Clone the repository:

```bash
git clone https://github.com/recreatedev/Resumify.git
cd Resumify
```

2. Install dependencies:

```bash
# Install frontend dependencies
bun install

# Install backend dependencies
cd apps/backend
go mod download
```

3. Set up environment variables:

```bash
cp apps/backend/.env.example apps/backend/.env
# Edit apps/backend/.env with your configuration
```

4. Start the database and Redis:

```bash
cd apps/backend
docker compose up -d
```

5. Run database migrations:

```bash
cd apps/backend
task migrations:up
```

6. Start the development server:

```bash
# From root directory
bun dev

# Or just the backend
cd apps/backend
task run
```

The API will be available at `http://localhost:8080` and the frontend at `http://localhost:5173`

## Development

### Available Commands

```bash
# Backend commands (from backend/ directory)
task help              # Show all available tasks
task run               # Run the application
task migrations:new    # Create a new migration
task migrations:up     # Apply migrations
task test              # Run tests
task tidy              # Format code and manage dependencies

# Frontend commands (from root directory)
bun dev                # Start development servers
bun build              # Build all packages
bun lint               # Lint all packages
```

### Environment Variables

The backend uses environment variables prefixed with `RESUMIFY_`. Key variables include:

- `RESUMIFY_DATABASE_*` - PostgreSQL connection settings
- `RESUMIFY_SERVER_*` - Server configuration
- `RESUMIFY_AUTH_*` - Authentication settings
- `RESUMIFY_REDIS_*` - Redis connection
- `RESUMIFY_INTEGRATION_*` - Third-party service configuration
- `RESUMIFY_OBSERVABILITY_*` - Monitoring settings

See `apps/backend/.env.example` for a complete list.

## Architecture

This application follows clean architecture principles:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Models**: Domain entities (Resume, Education, Experience, etc.)
- **Infrastructure**: External services (database, cache, email)

## Resume Features

### Resume Sections

- **Personal Information**: Name, contact details, summary
- **Education**: Institution, degree, dates, achievements
- **Experience**: Company, position, responsibilities, dates
- **Projects**: Project details, technologies, links
- **Skills**: Technical and soft skills with proficiency levels
- **Certifications**: Professional certifications and credentials

### Customization Options

- **Themes**: Multiple professional themes
- **Section Ordering**: Drag-and-drop section reordering
- **Visibility Control**: Show/hide sections
- **Content Formatting**: Rich text editing capabilities

## Testing

```bash
# Run backend tests
cd apps/backend
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires Docker)
go test -tags=integration ./...
```

## Production Considerations

1. Use environment-specific configuration
2. Enable production logging levels
3. Configure proper database connection pooling
4. Set up monitoring and alerting
5. Use a reverse proxy (nginx, Caddy)
6. Enable rate limiting and security headers
7. Configure CORS for your domains
8. Set up automated backups
9. Implement CDN for static assets
10. Configure SSL/TLS certificates

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
