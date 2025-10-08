# Resumify Frontend

A modern React-based frontend application for the Resumify resume builder, built with TypeScript, Vite, and Tailwind CSS.

## Features

- **Resume Builder**: Intuitive drag-and-drop interface for creating resumes
- **Real-time Preview**: Live preview of resume changes
- **Multiple Themes**: Professional, modern, classic, and creative templates
- **Section Management**: Add, remove, and reorder resume sections
- **Rich Content Editor**: WYSIWYG editing for all resume content
- **PDF Export**: Generate and download professional PDF resumes
- **Responsive Design**: Optimized for desktop, tablet, and mobile
- **User Authentication**: Secure login and user management
- **Auto-save**: Automatic saving of resume changes

## Tech Stack

- **React 18**: Modern React with hooks and concurrent features
- **TypeScript**: Type-safe development
- **Vite**: Fast build tool and development server
- **Tailwind CSS**: Utility-first CSS framework
- **React Query**: Data fetching and caching
- **React Hook Form**: Form handling and validation
- **Zod**: Schema validation
- **React Router**: Client-side routing
- **Lucide React**: Beautiful icons

## Getting Started

### Prerequisites

- Node.js 22+
- Bun (recommended) or npm

### Installation

1. Install dependencies:

```bash
bun install
```

2. Set up environment variables:

```bash
cp .env.example .env.local
# Configure your environment variables
```

3. Start the development server:

```bash
bun dev
```

The application will be available at `http://localhost:5173`

## Environment Variables

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_CLERK_PUBLISHABLE_KEY=your-clerk-key
```

## Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable UI components
│   ├── pages/            # Page components
│   ├── hooks/            # Custom React hooks
│   ├── services/         # API service layer
│   ├── types/            # TypeScript type definitions
│   ├── utils/            # Utility functions
│   ├── styles/           # Global styles
│   └── App.tsx           # Main application component
├── public/               # Static assets
└── package.json         # Dependencies and scripts
```

## Available Scripts

```bash
bun dev          # Start development server
bun build        # Build for production
bun preview      # Preview production build
bun lint         # Run ESLint
bun type-check   # Run TypeScript compiler
```

## Components

### Core Components

- **ResumeBuilder**: Main resume editing interface
- **ResumePreview**: Live preview component
- **SectionEditor**: Individual section editing
- **ThemeSelector**: Theme selection interface
- **PDFExporter**: PDF generation and download

### Form Components

- **EducationForm**: Education entry form
- **ExperienceForm**: Work experience form
- **ProjectForm**: Project details form
- **SkillForm**: Skills and certifications form

### UI Components

- **Button**: Reusable button component
- **Input**: Form input components
- **Modal**: Modal dialog component
- **Toast**: Notification component
- **Loading**: Loading states

## API Integration

The frontend communicates with the backend API through a service layer:

```typescript
// Example API service
export const resumeService = {
  getResumes: () => api.get('/resumes'),
  createResume: (data: CreateResumeRequest) => api.post('/resumes', data),
  updateResume: (id: string, data: UpdateResumeRequest) => api.put(`/resumes/${id}`, data),
  deleteResume: (id: string) => api.delete(`/resumes/${id}`),
}
```

## State Management

The application uses React Query for server state management and React Context for client state:

- **Server State**: Cached API responses with React Query
- **Client State**: Form state, UI state, and user preferences
- **Optimistic Updates**: Immediate UI updates with rollback on error

## Styling

The application uses Tailwind CSS for styling with a custom design system:

- **Color Palette**: Professional color scheme
- **Typography**: Consistent font hierarchy
- **Spacing**: Consistent spacing scale
- **Components**: Reusable styled components
- **Responsive**: Mobile-first responsive design

## Testing

```bash
bun test              # Run unit tests
bun test:coverage     # Run tests with coverage
bun test:e2e          # Run end-to-end tests
```

## Building for Production

```bash
bun build
```

The build output will be in the `dist/` directory, ready for deployment.

## Deployment

The frontend can be deployed to various platforms:

- **Vercel**: Zero-config deployment
- **Netlify**: Static site hosting
- **AWS S3**: Static website hosting
- **Docker**: Containerized deployment

## Contributing

1. Follow React and TypeScript best practices
2. Write tests for new components
3. Use semantic commit messages
4. Ensure responsive design
5. Maintain accessibility standards

## License

See the parent project's LICENSE file.