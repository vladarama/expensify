# Expensify

A modern expense tracking application built with React, Go, and PostgreSQL that helps users manage their personal finances effectively.
<img width="1412" alt="Screenshot 2024-12-02 at 15 02 41" src="https://github.com/user-attachments/assets/e756da52-e453-4545-957b-7bd64622aef1">

## Features

- **Unified Tracking**: Consolidate cash, credit, and debit expenses in one place
- **Real-time Insights**: Visualize spending patterns and budgets through interactive charts
- **Budget Management**: Set and track budgets by category
- **Income Tracking**: Monitor multiple income sources

## Tech Stack

### Frontend
- React.js
- TypeScript
- Vite
- Tailwind CSS

### Backend
- Go
- PostgreSQL

### DevOps
- Docker & Docker Compose
- Azure Container Registry (ACR)
- Azure Web Apps
- Azure DevOps Pipelines

## Data Model
<img width="562" alt="Screenshot 2024-12-02 at 15 08 18" src="https://github.com/user-attachments/assets/7fb47878-9d31-4e2f-a5d2-833f8eb4f05e">

- **Expense**: Tracks individual expenses with amount, date, and category
- **Income**: Records income sources and amounts
- **Budget**: Manages spending limits by category
- **Category**: Organizes expenses into logical groups

## Getting Started

### Prerequisites
- Node.js 16+
- Go 1.23+
- Docker and Docker Compose
- PostgreSQL

### Local Development

1. Clone the repository
2. Start the services using Docker Compose:
`docker compose up`

The application will be available at:
- Frontend: http://localhost:3000
- Backend: http://localhost:8080

### Testing
Run the Go tests:
- `cd backend`
- `go test ./internal/tests/integration`
- `go test ./internal/tests/models`

## CI/CD Pipeline
The project uses Azure DevOps Pipelines for continuous integration and deployment:

- Frontend Pipeline: Builds and deploys the React application
- Backend Pipeline: Builds, tests, and deploys the Go backend
- Automated Testing: Runs unit and integration tests
- Docker Integration: Builds and pushes container images to Azure Container Registry

## Deployment
The application is deployed using Azure Web Apps for Containers:
- Frontend: Azure Web App (React)
- Backend: Azure Web App (Go)
- Database: PostgreSQL deployed on Railway
