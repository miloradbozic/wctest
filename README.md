# Employee Organization Tree

A full-stack application that displays an organization's employee hierarchy in a tree structure. Built with Go (backend) and React (frontend).

## Prerequisites

- Go 1.16 or later
- Node.js 14 or later
- npm or yarn

## Setup


1. Start the backend:
   ```bash
   make backend
   ```

2. In a separate terminal, start the frontend:
   ```bash
   make frontend
   ```

The application will be available at:
- Backend API: http://localhost:8080/api/employees
- Frontend: http://localhost:3000


## API Endpoints

- `GET /api/employees` - Returns the employee organization tree