# 🏛️ goJudge

goJudge is a comprehensive online judge system designed to facilitate coding problem-solving and evaluation. This project encompasses various components working together to provide a seamless experience for both problem creators and solvers.

## 📚 Table of Contents

- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [🛠️ Tech Stack](#️-tech-stack)
- [🧩 Components](#-components)
  - [🖥️ API Server](#️-api-server)
  - [🤖 Worker](#-worker)
  - [🌐 Frontend](#-frontend)
- [🔄 Workflow](#-workflow)
- [🚀 Getting Started](#-getting-started)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## ✨ Features

- Dynamic problem creation with custom test cases
- Support for multiple programming languages
- Isolated code execution environment using Docker
- Real-time submission status updates
- Efficient queue management for handling multiple submissions
- Secure and scalable architecture

## 🏗️ Architecture

goJudge follows a microservices architecture, with the following main components:

1. API Server
2. Worker
3. Frontend
4. Database (PostgreSQL)
5. Queue System (Redis)
6. Storage (Amazon S3)

## 🛠️ Tech Stack

- **API Server**: Go, GORM 
- **Worker**: Go, Docker, Redis Queue
- **Database**: PostgreSQL
- **Frontend**: Next.js, TypeScript, shadcn/ui, React Query (Tanstack Query) 🔄
- **Queue**: Redis 
- **Storage**: Supabase S3 

## 🧩 Components

### 🖥️ API Server

The API server is the backbone of goJudge, built using Go and GORM. It handles:

- Problem creation and management
- User authentication and authorization
- Submission handling and result updates
- Communication with the frontend and worker

Key features:
- RESTful API design for easy integration
- Efficient database operations using GORM
- Scalable architecture to handle high loads using Redis queue

### 🤖 Worker

The worker is the powerhouse of goJudge, responsible for:

- Fetching submissions from the Redis queue
- Generating test cases based on the problem and language
- Injecting test cases into the user-submitted code
- Executing the code in an isolated Docker environment
- Sending results back to the API server via webhook and updating the submission status

Key features:
- Parallel processing of multiple submissions using Redis queue
- Secure code execution in isolated containers using Docker
- Dynamic test case generation and injection based on problem specifications

### 🌐 Frontend

The frontend is the face of goJudge, built using Next.js, TypeScript, and shadcn/ui. It provides:

- 👨‍💻 User interface for problem solving and submission
- 🔄 Real-time updates on submission status using React Query
- 🎨 Sleek and responsive design using shadcn/ui
- 👑 Problem creation interface for administrators

Key features:
- Type-safe development with TypeScript
- Efficient state management and data fetching with React Query
- Customizable UI components using shadcn/ui
- Real-time updates using Polling

## 🔄 Workflow

1. **Problem Creation** 📝:
   - Admin creates a problem, specifying:
     - Function name
     - Parameters
     - Return type
     - Test cases
   - System generates boilerplate code for supported languages
   - Boilerplate code is stored in Amazon S3 for quick access

2. **Problem Solving** 🧠:
   - User browses and selects a problem to solve
   - System fetches appropriate boilerplate code from S3
   - User writes their solution using the provided boilerplate
   - User submits their solution for evaluation

3. **Submission Processing** ⚙️:
   - Submitted code is added to the Redis queue for processing
   - GoJudge Worker picks up the submission from the queue
   - Worker generates test cases based on problem specifications
   - Test cases are injected into the submitted code
   - Code is executed in an isolated Docker environment based on the specific language for security

4. **Result Handling** 📊:
   - Worker analyzes execution results and determines correctness
   - Results are sent back to the API server via webhook
   - API server updates the submission status (Accepted/Rejected)
   - Frontend polls for results using React Query
   - UI is updated in real-time with the submission status

## 🚀 Getting Started

### Prerequisites

- Go 1.16+
- Node.js 14+
- Docker
- Redis
- Supabase account for S3 storage and PostgreSQL

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/goJudge.git
   cd goJudge
   ```

2. Set up the API server:
   ```bash
   cd go-api
   go mod tidy
   cp .env.example .env
   # Edit .env with your configuration
   go run main.go
   ```

3. Set up the worker:
   ```bash
   cd ../gojudge
   go mod tidy
   # Edit .env with your configuration
   go run main.go
   ```

4. Set up the frontend:
   ```bash
   cd ../frontend
   npm install
   cp .env.local.example .env.local
   # Edit .env.local with your configuration
   npm run dev
   ```

5. Set up the database:
   - Create a PostgreSQL database in Supbase
   

6. Configure Redis and S3:
   - Set up a Redis instance using docker compose file by running `docker-compose up -d`
   - Create an S3 bucket for storing boilerplate code
