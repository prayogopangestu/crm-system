# Railway Deployment Guide for Monolithic CRM System

## Overview
This project uses a monolithic architecture with:
- **Frontend**: Next.js 16 (React 19)
- **Backend**: Go API with Chi router
- **Databases**: PostgreSQL + Redis

Railway supports both frontend and backend deployment in a single project.

## Prerequisites
- Railway account (sign up at railway.app)
- GitHub repository with this project
- Railway CLI (optional but recommended)

## Deployment Steps

### Option 1: Railway Dashboard

1. **Create New Project**
   - Go to [railway.app/new](https://railway.app/new)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository

2. **Configure Services**
   Railway will automatically detect and create services:
   - Frontend (Next.js)
   - Backend (Go)
   - PostgreSQL Database
   - Redis Cache

3. **Environment Variables**
   Set these in each service:

   **Frontend Service:**
   ```
   NODE_ENV=production
   NEXT_PUBLIC_API_URL=${{ backend.RAILWAY_PRIVATE_DOMAIN }}/api
   ```

   **Backend Service:**
   ```
   GO_ENV=production
   DB_HOST=${{ postgres.RAILWAY_PRIVATE_DOMAIN }}
   DB_PORT=5432
   DB_USER=${{ postgres.POSTGRES_USER }}
   DB_PASSWORD=${{ postgres.POSTGRES_PASSWORD }}
   DB_NAME=${{ postgres.POSTGRES_DB }}
   REDIS_HOST=${{ redis.RAILWAY_PRIVATE_DOMAIN }}
   REDIS_PORT=6379
   REDIS_PASSWORD=${{ redis.REDIS_PASSWORD }}
   JWT_SECRET=${{ JWT_SECRET }}
   ```

4. **Build & Deploy**
   - Click "Deploy" button
   - Railway will build both services automatically
   - Wait for deployment to complete

### Option 2: Railway CLI

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login to Railway
railway login

# Initialize project
railway init
railway add

# Create services
railway add frontend
railway add backend
railway add postgres
railway add redis

# Configure environment variables
railway variables set NODE_ENV=production
railway variables set NEXT_PUBLIC_API_URL=${{ backend.RAILWAY_PRIVATE_DOMAIN }}/api

# Deploy
railway up
```

## Service Configuration

### Frontend Service
- **Type**: Web Service
- **Build**: Next.js production build
- **Port**: 3000
- **Health Check**: `/`
- **Environment**: Node.js 20+

### Backend Service
- **Type**: Web Service  
- **Build**: Go binary
- **Port**: 8080
- **Health Check**: `/health`
- **Environment**: Go 1.25+

### Database Services
- **PostgreSQL**: Main application database
- **Redis**: Caching and session storage

## Environment Variables

### Required Variables

Create these in Railway Dashboard (add as reference variables):

```
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### Service References

Railway automatically creates service references:
- `${{ postgres.RAILWAY_PRIVATE_DOMAIN }}` - PostgreSQL host
- `${{ redis.RAILWAY_PRIVATE_DOMAIN }}` - Redis host  
- `${{ backend.RAILWAY_PRIVATE_DOMAIN }}` - Backend API URL

## Domain Configuration

1. **Frontend Domain**
   - Go to frontend service settings
   - Click "Generate Domain" or add custom domain
   - Update DNS records if using custom domain

2. **Backend Domain**
   - Backend automatically gets private domain
   - Accessible only via frontend service (for security)
   - Public domain available if needed for API access

## Railway Configuration File

`railway.toml` configuration is included in the project root for automatic service detection.

## Database Migration

Run database migrations after first deployment:

```bash
# Access backend service console
railway shell backend

# Run migrations
go run cmd/migrate/main.go

# Or seed data
go run cmd/seed/main.go
```

## Scaling

### Frontend Scaling
- Minimum: 1 instance
- Maximum: Auto-scales based on traffic
- Memory: 512MB - 2GB

### Backend Scaling  
- Minimum: 1 instance
- Maximum: Auto-scales based on API load
- Memory: 512MB - 2GB

### Database Scaling
- PostgreSQL: Automatic backups, auto-scaling storage
- Redis: Basic tier (can upgrade to Pro)

## Monitoring & Logs

### View Logs
```bash
# All services
railway logs

# Specific service
railway logs frontend
railway logs backend

# Real-time logs
railway logs -f
```

### Metrics
- Go to Railway Dashboard
- View metrics tab for each service
- Monitor CPU, memory, and request metrics

## Troubleshooting

### Build Errors
1. Check build logs in Railway Dashboard
2. Ensure Node.js 18+ and Go 1.25+ versions
3. Verify all dependencies in package.json and go.mod

### Database Connection Issues
1. Verify database environment variables
2. Check database service is running
3. Test connection from backend console

### Frontend-Backend Communication
1. Ensure `NEXT_PUBLIC_API_URL` points to correct backend domain
2. Check backend CORS configuration
3. Verify both services are deployed successfully

### Deployment Fails
1. Clear cache: Railway Dashboard > Settings > Clear Build Cache
2. Check for syntax errors in configuration files
3. Review build logs for specific errors

## Backup & Recovery

### Automatic Backups
- PostgreSQL: Daily automatic backups
- Redis: No automatic backups (data is cache)

### Manual Backup
```bash
# Export PostgreSQL
railway exec postgres pg_dump > backup.sql

# Restore PostgreSQL  
cat backup.sql | railway exec postgres psql
```

## Cost Management

### Free Tier Limits
- **Frontend**: 500 hours/month free
- **Backend**: 500 hours/month free  
- **PostgreSQL**: 5GB storage free
- **Redis**: Basic tier free

### Upgrade Options
- Upgrade to Pro for:
  - More resources
  - Better performance
  - Priority support
  - Custom domains

## Best Practices

1. **Environment Variables**: Never commit secrets to git
2. **Database Backups**: Regular backups of PostgreSQL
3. **Monitoring**: Set up alerts for critical errors
4. **Scaling**: Start with minimum, scale as needed
5. **Security**: Keep dependencies updated

## Continuous Deployment

Railway automatically deploys on push to main branch:
1. Push changes to GitHub
2. Railway triggers new build
3. Services update automatically
4. Rollback available if needed

## Support

- Railway Documentation: [docs.railway.app](https://docs.railway.app)
- Support: [railway.app/contact](https://railway.app/contact)
- GitHub Issues: [github.com/railwayapp](https://github.com/railwayapp)

## Next Steps

1. Deploy to Railway using this guide
2. Configure custom domains
3. Set up monitoring and alerts
4. Test all functionality end-to-end
5. Set up CI/CD pipeline