
# Project Overview

This project includes a **backend** built in Go and a **frontend** using React with TypeScript, styled with Tailwind CSS. It leverages Docker for containerization and uses `docker-compose` for managing services.

## Folder Structure
- **backend**: Go-based server code with internal modules.
- **frontend**: React app with Tailwind CSS and Vite for development.
- **azure-pipelines.yml**: Azure DevOps CI/CD pipeline configuration.
- **docker-compose.yml**: Multi-service setup for backend and frontend.

## Quick Start
1. Build and run services:
   ```bash
   docker-compose up --build
   ```
2. Access the frontend at `http://localhost:3000` and backend at `http://localhost:8080`.
