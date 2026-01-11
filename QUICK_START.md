# IronTrack - Quick Start Guide

## Running the Full Stack

### Terminal 1: Start Backend
```bash
cd /Users/hankmendix/side-projects/IronTrack-AI-Backend
# Make sure environment variables are set
export DATABASE_URL="your-database-url"
export GEMINI_API_KEY="your-api-key"
export JWT_SECRET="your-secret"
export ALLOWED_ORIGINS="http://localhost:5173"

# Run backend
go run cmd/server/main.go
# Or if using a main.go at root
go run main.go
```

Backend will be available at: `http://localhost:8080`

Health check: `curl http://localhost:8080/health`

### Terminal 2: Start Admin Dashboard
```bash
cd /Users/hankmendix/side-projects/IronTrack-AI-Backend/admin-app
npm install  # First time only
npm run dev
```

Admin dashboard will be available at: `http://localhost:5173`

## Initial Setup

### 1. Create Admin User (Choose One)

**Option A: Direct Database Query**
```sql
-- After creating a regular user, update them to admin:
UPDATE users SET is_admin = true WHERE id = '<user-id>';
```

**Option B: Via Backend API**
First, create a regular user:
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@example.com",
    "password": "password123"
  }'
```

Then manually update database:
```sql
UPDATE users SET is_admin = true WHERE email = 'admin@example.com';
```

### 2. Login to Dashboard
1. Open `http://localhost:5173/login`
2. Enter admin credentials
3. You'll be redirected to dashboard

### 3. Start Using Admin Panel
- **Dashboard**: View summary metrics
- **Users**: Create/edit/delete users
- **Plans**: Create and manage workout plans
- **Exercises**: Define global and user-specific exercises
- **AI Requests**: Monitor AI usage

## Available Admin Routes

### Summary
```
GET /api/admin/summary
Returns: { users, plans, exercises, aiRequests }
```

### Users
```
GET /api/admin/users                    # List all users
POST /api/admin/users                   # Create user
PUT /api/admin/users/:id                # Update user
DELETE /api/admin/users/:id             # Delete user
```

### Plans
```
GET /api/admin/plans                    # List all plans
POST /api/admin/plans                   # Create plan
DELETE /api/admin/plans/:id             # Delete plan
```

### Exercises
```
GET /api/admin/exercises                # List all exercises
POST /api/admin/exercises               # Create exercise
DELETE /api/admin/exercises/:id         # Delete exercise
```

### AI Requests
```
GET /api/admin/ai-requests              # View AI usage log
```

## Environment Variables

### Backend (.env or system)
```
DATABASE_URL=postgresql://user:pass@localhost/irontrack
DB_HOST=localhost
DB_PORT=5432
DB_USER=hankmendix
DB_PASSWORD=your_password
DB_NAME=irontrack
DB_SSLMODE=disable

GEMINI_API_KEY=your_gemini_key
JWT_SECRET=your_secret_key_change_in_production
ALLOWED_ORIGINS=http://localhost:5173,https://yourdomain.com
TRUSTED_PROXIES=

# Optional for Render.com or cloud deployment
# TRUSTED_PROXIES=203.0.113.0/24,198.51.100.0/24
```

### Frontend (admin-app/.env.local)
```
VITE_API_URL=http://localhost:8080/api
```

For production:
```
VITE_API_URL=https://api.irontrack.com/api
```

## Troubleshooting

### Admin Dashboard shows "Unauthorized"
- Make sure user has `isAdmin: true` in database
- Check browser console for JWT errors
- Verify backend is running and accessible

### Backend 401 errors
- Check JWT_SECRET is same in backend
- Verify token hasn't expired (30 days)
- Clear localStorage in browser and re-login

### CORS errors
- Add frontend URL to ALLOWED_ORIGINS in backend
- For local dev: `ALLOWED_ORIGINS=http://localhost:5173`

### AI Features not working
- Verify GEMINI_API_KEY is set
- Check Gemini API quota
- View backend logs for errors

### Database connection fails
- Verify PostgreSQL is running
- Check DATABASE_URL format
- Ensure database exists and user has permissions

## Project Structure Overview

```
IronTrack-AI-Backend/
├── admin-app/               # React admin dashboard
│   ├── src/
│   ├── package.json
│   ├── tailwind.config.js
│   └── SETUP.md
├── internal/
│   ├── auth/               # JWT authentication
│   ├── database/           # Database connection
│   ├── handlers/           # API handlers (including admin)
│   ├── models/             # Data models
│   └── router/             # Route definitions
├── cmd/
│   ├── server/main.go      # Server entrypoint
│   └── ...
├── Dockerfile              # Container configuration
├── go.mod                  # Go dependencies
├── render.yaml             # Render deployment config
└── ADMIN_DASHBOARD_IMPLEMENTATION.md
```

## Development Workflow

### Making changes to backend
1. Modify handlers or routes
2. Backend auto-reloads (if using hot reload)
3. Or restart: `Ctrl+C` then `go run cmd/server/main.go`

### Making changes to frontend
1. Edit React components
2. Browser auto-refreshes (Vite HMR)
3. No rebuild needed during development

### Testing
Backend: `curl` or Postman for API endpoints
Frontend: Browser DevTools for component inspection

## Deployment

### Backend
Use render.yaml configuration:
```bash
# Deploy to Render
git push origin main
```

Or use Docker:
```bash
docker build -t irontrack-backend .
docker run -p 8080:8080 irontrack-backend
```

### Frontend (Admin Dashboard)
```bash
# Build for production
cd admin-app
npm run build
```

Deploy `dist/` folder to:
- Vercel: `vercel deploy --prod`
- Netlify: `netlify deploy --prod`
- Or any static hosting service

## Key Features Implemented

### Backend
✅ User management (CRUD)
✅ Workout plans (CRUD)
✅ Exercises (global + user-specific)
✅ Workout logs
✅ AI plan generation
✅ AI progress reports
✅ Admin dashboard APIs
✅ AI request tracking

### Frontend
✅ Authentication (JWT)
✅ Role-based access (admin only)
✅ Responsive design
✅ User management UI
✅ Plan management UI
✅ Exercise management UI
✅ AI monitoring
✅ Real-time dashboard

## Support & Resources

- Backend Router: [internal/router/router.go](internal/router/router.go)
- Admin Handlers: [internal/handlers/admin_handler.go](internal/handlers/admin_handler.go)
- Admin App Setup: [admin-app/SETUP.md](admin-app/SETUP.md)
- API Client: [admin-app/src/lib/api.ts](admin-app/src/lib/api.ts)

## Next Steps

1. ✅ Set up database and backend
2. ✅ Create admin user
3. ✅ Start backend server
4. ✅ Start admin dashboard
5. ✅ Login and start managing!
6. Consider: Adding more features, deploying to production, monitoring usage

---

**Happy building! 🚀**
