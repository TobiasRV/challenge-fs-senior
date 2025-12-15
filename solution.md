# Task Managaer - Solution Documentation

## Quick Start

Get the entire application running with a single command:

```bash
git clone https://github.com/TobiasRV/challenge-fs-senior.git
cd challenge-fs-senior
cp .env.example .env
docker-compose up
```

The application will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: PostgreSQL on port 5432

## Technology Stack

### Backend
- **Language**: Golang
- **Framework**: Fiber
- **Database**: PostgreSQL with sqlc and squirrel
- **Authentication**: JWT with refresh tokens

### Frontend
- **Language**: TypeScript
- **Framework**: React
- **State Management**: Zustand
- **Styling**: Tailwind CSS
- **Forms**: React Hook Form
- **HTTP Client**: Axios with interceptors
- **UI Components**: Shadcn, Radix-ui

### DevOps & Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for local development
- **Database Migrations**: goose migrations and seeds
- **Environment Management**: Environment variables

## Architecture Overview

### System Design

This task management platform consists of a simple but functional web page capable of function in screens of any size that connect to a RESTful api with jwt for authentication and a postgreSQL database for reliability and performance.

### Database Schema Design
The database follows a normalized structure with proper foreign key relationships:

<img width="1166" height="768" alt="database schema" src="https://github.com/user-attachments/assets/2a27b4f3-765a-43d1-8a7d-3ef4aaaadefb" />

## Setup Instructions

### Prerequisites
- Docker and Docker Compose
- Go (for local development without Docker)
- Node (for local development without Docker)
- PostgreSQL(for local development without Docker)

### Docker Setup (Recommended)

1. **Clone and Setup**:
   ```bash
   git clone <your-repository-url>
   cd challenge-fs-senior
   cp .env.example .env
   ```

2. **Start Services**:
   ```bash
   docker-compose up -d
   ```

3. **Verify Services**:
   ```bash
   docker-compose ps
   ```

4. **View Logs**:
   ```bash
   docker-compose logs -f
   ```

## API Documentation

The project has a swagger file for documentation it can be access in the /swagger route

### Authentication Flow
1. **Registration/Login** → Returns JWT + Refresh Token
2. **API Requests** → Include `Authorization: Bearer <token>`
3. **Token Refresh** → Use refresh token when JWT expires
4. **Logout** → Invalidate refresh token

## Database Schema

### Core Tables

**users**
- Primary user accounts with role-based access
- Has a relation with the team in case of non admin users in the field team_id

**teams**
- Organization units for grouping projects
- Has a relation with the user table in the owner field

**projects**
- Work containers within teams
- Status tracking
- Has a relation with the user table in the manager field
- Has relation with the team in the team_id field

**tasks**
- Core work items
- Has status tracking
- Has relation with the assigned user in the user_id field
- Has relation with the project it belong with the project_id field

## Demo Access

### Live Demo
- **URL**: [[Deployed Application URL](https://challenge-fs-senior-91jox8nf9-tobiasrvs-projects.vercel.app/)]

### Test Accounts
```
Admin Account:
  Email: admin@admin.com
  Password: 1234
  
Manager Account:
  Email: manager@manager.com  
  Password: 1234

Manager Account:
    Email: manager2@manager.com  
    Password: 1234
  
Member Account:
  Email: member@demo.com
  Password: 1234

Member Account:
  Email: member2@demo.com
  Password: 1234
```
## Known Issues & Limitations

## Folder structure changes
- The database migrations/seeds are inside the backend/internal/sql folder instead of in a separate folder like the documentation because is used by the backend for the queries using sqlc.

## Software limitations

- When you register you can only create an admin user, users of other roles can only be created by the admin
- The admin user can't create project or tasks, it can only see them
- The members can't edit the 

### Technical Debt
1. **Test Coverage**: Frontend testing and E2E testing
2. **Logging and Monitoring**: Adding a logger to have better visibility of erros in the backend
3. **Backup Strategy**: Automated database backups and disaster recovery procedures

## Future Improvements (features not listed on the requirements)

1. Allowing user to modify the status of its assigned tasks
2. Allow admin to have access to all CRUDs
3. Password creation email when admin create new user
5. Allow update of project manager

## Time Investment Breakdown

- **High level functionality design**: 3 hours
- **Database modeling**: 2 hours
- **Backend development**: 12 hours
- **Backend testing**: 5 hours
- **Frontend development**: 8 hours
- **Documentation**: 3 hours
- **Docker Configuration**: 2 hours
- **Deployment**: 1 hour


**Total Investment**: ~36 hours over 4 days
