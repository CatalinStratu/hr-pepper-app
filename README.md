# Pepper HR Interview Application

A comprehensive HR interview system powered by the SoftBank Pepper robot, featuring:
- **Pepper Android App**: Interactive interview interface on the robot
- **Golang Backend Server**: REST API for data management and response analysis
- **React Dashboard**: Web interface for HR staff to manage interviews

## Architecture

```
┌─────────────────┐     HTTP/REST      ┌─────────────────┐
│  Pepper Robot   │ ◄──────────────────► │  Golang Server  │
│  (Android App)  │                      │  (REST API)     │
└─────────────────┘                      └────────┬────────┘
                                                  │
                                                  │
                                         ┌────────▼────────┐
                                         │   MySQL DB      │
                                         └────────┬────────┘
                                                  │
┌─────────────────┐     HTTP/REST      ┌─────────▼────────┐
│ React Dashboard │ ◄──────────────────► │  Golang Server  │
│   (Web UI)      │                      │                 │
└─────────────────┘                      └─────────────────┘
```

## Components

### 1. Backend Server (`/server`)
- Go 1.21+
- REST API for interviews, questions, and responses
- MySQL database with connection pooling
- Response analysis with sentiment scoring
- CORS enabled for cross-origin requests

### 2. Pepper Android App (`/pepper-app`)
- Android SDK with QiSDK for Pepper robot
- Displays questions on tablet display
- Robot speaks questions using Text-to-Speech
- Speech recognition for voice responses
- Real-time response submission to server
- Timer for each question

### 3. Dashboard (`/dashboard`)
- React 18 with TypeScript
- Vite for fast development
- Interview management and analytics
- Question CRUD operations
- Score visualization with charts
- Candidate response review

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# Dashboard: http://localhost:3000
# API Server: http://localhost:8080
# MySQL: localhost:3306
```

### Manual Setup

#### Prerequisites
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Android Studio (for Pepper app)

#### MySQL Database

```sql
CREATE DATABASE hr_interview;
CREATE USER 'hruser'@'localhost' IDENTIFIED BY 'hrpassword';
GRANT ALL PRIVILEGES ON hr_interview.* TO 'hruser'@'localhost';
FLUSH PRIVILEGES;
```

#### Backend Server

```bash
cd server

# Set environment variables
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=hruser
export DB_PASSWORD=hrpassword
export DB_NAME=hr_interview
export PORT=8080

# Run the server
go mod download
go run main.go
```

#### Dashboard

```bash
cd dashboard
npm install
npm run dev
# Dashboard runs on http://localhost:3000
```

#### Pepper App

1. Open `/pepper-app` in Android Studio
2. Update `SERVER_URL` in `app/build.gradle` to your server IP
3. Connect to Pepper robot (or emulator)
4. Build and deploy

## API Endpoints

### Questions
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/questions | Get all active questions |
| GET | /api/questions/:id | Get question by ID |
| POST | /api/questions | Create new question |
| PUT | /api/questions/:id | Update question |
| DELETE | /api/questions/:id | Delete question (soft delete) |

### Interviews
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/interviews | Get all interviews (filter by status) |
| GET | /api/interviews/:id | Get interview with responses |
| POST | /api/interviews | Start new interview |
| PUT | /api/interviews/:id | Update interview status/notes |
| DELETE | /api/interviews/:id | Delete interview |
| GET | /api/interviews/:id/report | Get detailed interview report |

### Responses
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/responses | Submit interview response |

### Analytics
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/analytics | Get interview statistics |

## Environment Variables

### Server
| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | MySQL host |
| DB_PORT | 3306 | MySQL port |
| DB_USER | root | MySQL username |
| DB_PASSWORD | (empty) | MySQL password |
| DB_NAME | hr_interview | Database name |
| PORT | 8080 | Server port |

### Dashboard
| Variable | Default | Description |
|----------|---------|-------------|
| VITE_API_URL | /api | Backend API URL |

## Features

### Pepper Robot App
- Welcome greeting for candidates
- Displays and speaks interview questions
- Timer countdown for each question
- Voice recording with speech-to-text
- Manual text entry option
- Real-time feedback after each response
- Completion summary with score

### Dashboard
- **Analytics Dashboard**: Overview of interview statistics, score distribution, question performance
- **Interview Management**: View, filter, and manage all interviews
- **Interview Detail**: Review candidate responses with analysis
- **Question Management**: Add, edit, delete interview questions

### Response Analysis
The server performs basic text analysis on responses:
- Sentiment scoring based on positive/negative keywords
- Confidence scoring based on response structure
- Keyword detection for professional terms
- Feedback generation for each response

## Project Structure

```
hr-pepper-app/
├── server/                 # Go backend
│   ├── handlers/          # API handlers
│   ├── models/            # Data models
│   ├── database/          # MySQL connection
│   ├── middleware/        # HTTP middleware
│   └── main.go            # Entry point
├── dashboard/             # React frontend
│   ├── src/
│   │   ├── pages/        # Page components
│   │   ├── components/   # Reusable components
│   │   ├── services/     # API services
│   │   └── types/        # TypeScript types
│   └── package.json
├── pepper-app/            # Android app
│   └── app/src/main/
│       ├── java/         # Java source code
│       └── res/          # Android resources
├── docker-compose.yml
└── README.md
```

## License

MIT License
