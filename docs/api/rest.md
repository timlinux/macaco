# REST API Reference

MoCaCo provides a REST API for programmatic access.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently no authentication required (local-only mode).

## Endpoints

### Health Check

```http
GET /health
```

**Response:**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime_seconds": 3600
}
```

### Sessions

#### Create Session

```http
POST /sessions
Content-Type: application/json

{
  "round_type": "beginner"
}
```

**Response (201 Created):**

```json
{
  "session_id": "uuid",
  "round_type": "beginner",
  "total_tasks": 30,
  "current_task_index": 0,
  "started_at": "2026-02-21T10:00:00Z",
  "current_task": {
    "task_id": "motion-w-001",
    "category": "motion",
    "difficulty": 1,
    "initial": "hello world from vim",
    "desired": "hello world from vim",
    "cursor_start": 0,
    "description": "Move to next word start",
    "hint": "Use 'w' to move forward one word"
  }
}
```

#### Get Session

```http
GET /sessions/:session_id
```

#### Delete Session

```http
DELETE /sessions/:session_id
```

### Keystrokes

#### Send Keystroke

```http
POST /sessions/:session_id/keystroke
Content-Type: application/json

{
  "key": "w",
  "modifiers": []
}
```

**Response:**

```json
{
  "buffer_state": "hello world from vim",
  "cursor_position": 6,
  "current_mode": "NORMAL",
  "match_status": "complete",
  "task_completed": true,
  "elapsed_time_ms": 1234
}
```

#### Send Multiple Keystrokes

```http
POST /sessions/:session_id/keystrokes
Content-Type: application/json

{
  "keys": ["d", "w"]
}
```

### Task Operations

#### Complete Task

```http
POST /sessions/:session_id/complete
```

#### Skip Task

```http
POST /sessions/:session_id/skip
```

#### Reset Task

```http
POST /sessions/:session_id/reset
```

### Statistics

#### Get Session Statistics

```http
GET /sessions/:session_id/stats
```

#### Get Lifetime Statistics

```http
GET /stats/lifetime
```

#### Export Statistics

```http
GET /stats/export?format=json
GET /stats/export?format=csv
```

### Tasks

#### Get Task by ID

```http
GET /tasks/:task_id
```

### Rounds

#### Get Round Types

```http
GET /rounds
```

**Response:**

```json
{
  "round_types": ["beginner", "intermediate", "advanced", "expert", "mixed"]
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "SESSION_NOT_FOUND",
    "message": "Session with ID 'xyz' not found"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| SESSION_NOT_FOUND | 404 | Session doesn't exist |
| TASK_NOT_FOUND | 404 | Task doesn't exist |
| INVALID_ROUND_TYPE | 400 | Unknown round type |
| INVALID_REQUEST | 400 | Malformed request |
| NO_SKIPS_REMAINING | 400 | No skips left |
| INTERNAL_ERROR | 500 | Server error |

## Running the Server

```bash
# Start server only
macaco --server --addr localhost:8080

# Or use the script
./scripts/start-backend.sh
```
