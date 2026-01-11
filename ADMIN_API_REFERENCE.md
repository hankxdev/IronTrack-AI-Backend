# IronTrack Admin API - Complete Reference

## Authentication

### Login
```
POST /api/login
Content-Type: application/json

Request:
{
  "email": "admin@example.com",
  "password": "password123"
}

Response 200:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-string",
    "email": "admin@example.com",
    "name": "Admin User",
    "isAdmin": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### Get Current User
```
GET /api/me
Authorization: Bearer <token>

Response 200:
{
  "id": "uuid-string",
  "email": "admin@example.com",
  "name": "Admin User",
  "isAdmin": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

## Admin Summary

### Get Dashboard Summary
```
GET /api/admin/summary
Authorization: Bearer <token>
```

Response 200:
```json
{
  "users": 45,
  "plans": 120,
  "exercises": 256,
  "aiRequests": 890
}
```

---

## User Management

### List All Users
```
GET /api/admin/users
Authorization: Bearer <token>

Response 200:
[
  {
    "id": "uuid-1",
    "name": "John Doe",
    "email": "john@example.com",
    "isAdmin": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  ...
]
```

### Create User
```
POST /api/admin/users
Authorization: Bearer <token>
Content-Type: application/json

Request:
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123",
  "isAdmin": false
}

Response 201:
{
  "id": "uuid-new",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "isAdmin": false,
  "createdAt": "2024-01-11T12:00:00Z",
  "updatedAt": "2024-01-11T12:00:00Z"
}
```

### Update User
```
PUT /api/admin/users/:id
Authorization: Bearer <token>
Content-Type: application/json

Request (all fields optional):
{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "password": "newpassword123",
  "isAdmin": true
}

Response 200:
{
  "id": "uuid-1",
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "isAdmin": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-11T12:00:00Z"
}
```

### Delete User
```
DELETE /api/admin/users/:id
Authorization: Bearer <token>

Response 200:
{
  "message": "User deleted"
}
```

---

## Plan Management

### List All Plans
```
GET /api/admin/plans
Authorization: Bearer <token>

Response 200:
[
  {
    "id": "plan-uuid-1",
    "userId": "user-uuid-1",
    "name": "Push Day",
    "description": "Chest, shoulders, and triceps workout",
    "targetGoal": "Strength",
    "isAiGenerated": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "exercises": [
      {
        "name": "Bench Press",
        "defaultSets": 4,
        "defaultReps": 6,
        "muscleGroup": "Chest",
        "instructions": "Barbell bench press"
      },
      ...
    ]
  },
  ...
]
```

### Create Plan
```
POST /api/admin/plans
Authorization: Bearer <token>
Content-Type: application/json

Request:
{
  "userId": "user-uuid-1",
  "name": "Full Body",
  "description": "Complete full body workout",
  "targetGoal": "Muscle Gain",
  "exercises": [
    {
      "name": "Squats",
      "defaultSets": 4,
      "defaultReps": 8,
      "muscleGroup": "Legs",
      "instructions": "Barbell back squats"
    },
    {
      "name": "Deadlift",
      "defaultSets": 3,
      "defaultReps": 5,
      "muscleGroup": "Back",
      "instructions": "Conventional deadlift"
    }
  ]
}

Response 201:
{
  "id": "plan-uuid-new",
  "userId": "user-uuid-1",
  "name": "Full Body",
  "description": "Complete full body workout",
  "targetGoal": "Muscle Gain",
  "isAiGenerated": false,
  "createdAt": "2024-01-11T12:00:00Z",
  "exercises": [
    {
      "name": "Squats",
      "defaultSets": 4,
      "defaultReps": 8,
      "muscleGroup": "Legs",
      "instructions": "Barbell back squats"
    },
    ...
  ]
}
```

### Delete Plan
```
DELETE /api/admin/plans/:id
Authorization: Bearer <token>

Response 200:
{
  "message": "Plan deleted"
}
```

---

## Exercise Management

### List All Exercises
```
GET /api/admin/exercises
Authorization: Bearer <token>

Response 200:
[
  {
    "id": "ex-uuid-1",
    "userId": null,
    "name": "Bench Press",
    "muscleGroup": "Chest",
    "instructions": "Barbell bench press at 45 degrees",
    "isGlobal": true
  },
  {
    "id": "ex-uuid-2",
    "userId": "user-uuid-5",
    "name": "Custom Exercise",
    "muscleGroup": "Arms",
    "instructions": "User-defined exercise",
    "isGlobal": false
  },
  ...
]
```

### Create Exercise
```
POST /api/admin/exercises
Authorization: Bearer <token>
Content-Type: application/json

Request (for global):
{
  "name": "Pull-ups",
  "muscleGroup": "Back",
  "instructions": "Bodyweight pull-ups",
  "isGlobal": true
}

Request (for user-specific):
{
  "name": "Dumbbell Curls",
  "muscleGroup": "Arms",
  "instructions": "12lb dumbbells",
  "isGlobal": false,
  "userId": "user-uuid-1"
}

Response 201:
{
  "id": "ex-uuid-new",
  "userId": null,
  "name": "Pull-ups",
  "muscleGroup": "Back",
  "instructions": "Bodyweight pull-ups",
  "isGlobal": true
}
```

### Delete Exercise
```
DELETE /api/admin/exercises/:id
Authorization: Bearer <token>

Response 200:
{
  "message": "Exercise deleted"
}

Note: Admin users can delete any exercise (global or user-specific)
Non-admin users can only delete their own exercises
```

---

## AI Requests Monitoring

### List AI Requests
```
GET /api/admin/ai-requests
Authorization: Bearer <token>

Response 200:
[
  {
    "id": "req-uuid-1",
    "userId": "user-uuid-5",
    "type": "generate_plan",
    "createdAt": "2024-01-11T10:30:00Z"
  },
  {
    "id": "req-uuid-2",
    "userId": "user-uuid-3",
    "type": "generate_report",
    "createdAt": "2024-01-11T10:25:00Z"
  },
  {
    "id": "req-uuid-3",
    "userId": "user-uuid-5",
    "type": "generate_plan",
    "createdAt": "2024-01-11T10:20:00Z"
  },
  ...
]

Notes:
- Returns up to 200 most recent requests
- Sorted by createdAt (newest first)
- Types: "generate_plan", "generate_report"
```

---

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Invalid token"
}
```

### 403 Forbidden
```json
{
  "error": "Admin access required"
}
```

### 400 Bad Request
```json
{
  "error": "Invalid email or password"
}
```

### 404 Not Found
```json
{
  "error": "User not found"
}
```

### 409 Conflict
```json
{
  "error": "User already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create user"
}
```

---

## Data Models

### User
```typescript
{
  id: string;              // UUID
  email: string;           // Unique
  password: string;        // Hashed (bcrypt)
  name: string;
  isAdmin: boolean;        // Default: false
  createdAt: string;       // ISO 8601
  updatedAt: string;       // ISO 8601
}
```

### WorkoutPlan
```typescript
{
  id: string;              // UUID
  userId: string;          // UUID of plan owner
  name: string;
  description: string;
  targetGoal: string;
  isAiGenerated: boolean;
  createdAt: string;       // ISO 8601
  exercises: PlanExercise[];
}
```

### PlanExercise
```typescript
{
  name: string;
  defaultSets: number;     // Integer > 0
  defaultReps: number;     // Integer > 0
  muscleGroup?: string;
  instructions?: string;
}
```

### ExerciseDefinition
```typescript
{
  id: string;              // UUID
  userId?: string;         // Null if global
  name: string;
  muscleGroup: string;
  instructions?: string;
  isGlobal: boolean;
}
```

### AIRequestLog
```typescript
{
  id: string;              // UUID
  userId: string;          // Who used the AI
  type: string;            // "generate_plan" | "generate_report"
  createdAt: string;       // ISO 8601
}
```

### AdminSummary
```typescript
{
  users: number;           // Total user count
  plans: number;           // Total plan count
  exercises: number;       // Total exercise count
  aiRequests: number;      // Total AI usage count
}
```

---

## Headers

### Required Headers
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

### Response Headers
```
Content-Type: application/json
Access-Control-Allow-Origin: <ALLOWED_ORIGINS>
Access-Control-Allow-Credentials: true
```

---

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid/missing token |
| 403 | Forbidden - Admin access required |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server issue |

---

## Rate Limiting

No official rate limiting implemented. Consider adding for production:

```go
// Example rate limiting config
const MaxRequestsPerMinute = 60
const TokenRefreshWindow = 24 * time.Hour
```

---

## Security Notes

1. **JWT Tokens**: Valid for 30 days
2. **Passwords**: Hashed with bcrypt
3. **CORS**: Configured via `ALLOWED_ORIGINS` env var
4. **Admin Check**: Verified on every admin request
5. **SQL Injection**: Protected via GORM ORM
6. **HTTPS**: Recommended for production

---

## Example Requests

### Create Admin User
```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Admin",
    "email": "newadmin@example.com",
    "password": "secure_password",
    "isAdmin": true
  }'
```

### Get Dashboard Stats
```bash
curl -X GET http://localhost:8080/api/admin/summary \
  -H "Authorization: Bearer <token>"
```

### Create Workout Plan
```bash
curl -X POST http://localhost:8080/api/admin/plans \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-uuid-1",
    "name": "Beginner Full Body",
    "description": "3-day full body routine",
    "targetGoal": "Strength & Conditioning",
    "exercises": [
      {
        "name": "Squats",
        "defaultSets": 3,
        "defaultReps": 8,
        "muscleGroup": "Legs"
      }
    ]
  }'
```

---

## Testing with Postman

1. Import endpoints into Postman
2. Set `{{baseUrl}}` variable to `http://localhost:8080/api`
3. Set `{{token}}` variable after login
4. Add token to Authorization header: `Bearer {{token}}`
5. All admin requests will use admin collection with auto-token injection

---

## Changelog

### Version 1.0 (Current)
- ✅ Complete admin CRUD for users, plans, exercises
- ✅ AI request logging and monitoring
- ✅ Dashboard summary statistics
- ✅ Admin-only route protection
- ✅ JWT authentication
