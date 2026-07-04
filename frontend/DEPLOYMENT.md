# Vercel Deployment Guide

## Prerequisites
- Vercel account (sign up at vercel.com)
- GitHub repository with this project
- Backend service URL

## Environment Variables

Set these in Vercel Dashboard > Project Settings > Environment Variables:

### Required
```
NEXT_PUBLIC_API_URL=https://your-backend-url.com
```

### Optional (for development)
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Deployment Steps

### Option 1: Deploy from Vercel Dashboard
1. Go to [vercel.com/new](https://vercel.com/new)
2. Import your GitHub repository
3. Set **Root Directory** to `frontend`
4. Set **Framework Preset** to `Next.js`
5. Configure environment variables
6. Click **Deploy**

### Option 2: Deploy with Vercel CLI

```bash
# Install Vercel CLI
npm install -g vercel

# Login to Vercel
vercel login

# Deploy to preview
cd frontend
vercel

# Deploy to production
vercel --prod
```

### Option 3: GitHub Integration

1. Connect GitHub repository to Vercel
2. Configure project settings:
   - Root Directory: `frontend`
   - Build Command: `npm run build`
   - Output Directory: `.next`
3. Add environment variables
4. Every push to main branch will auto-deploy

## Configuration Files

- `vercel.json` - Vercel-specific configuration
- `.env.example` - Environment variable template
- `next.config.ts` - Next.js configuration

## Production Build

Locally test production build before deploying:

```bash
cd frontend
npm install
npm run build
npm run start
```

## Deployment Checklist

- [ ] Backend service is deployed and accessible
- [ ] Environment variables configured in Vercel
- [ ] `NEXT_PUBLIC_API_URL` points to production backend
- [ ] Test deployment in preview environment
- [ ] Verify all API calls work correctly
- [ ] Deploy to production

## Common Issues

### 1. API Connection Failed
- Verify `NEXT_PUBLIC_API_URL` is correct
- Check backend CORS settings
- Ensure backend is running and accessible

### 2. Build Errors
- Clear Vercel cache: Dashboard > Settings > Functions > Clear Cache
- Check Node.js version compatibility (needs 18.x+)
- Review build logs in Vercel Dashboard

### 3. Runtime Errors
- Check environment variables are set correctly
- Verify all required dependencies are installed
- Review browser console for specific errors

## Monitoring

View deployment logs and analytics in Vercel Dashboard:
- **Deployments**: Deployment history and logs
- **Functions**: Function execution logs
- **Analytics**: Performance metrics
- **Settings**: Domain, environment variables, and configuration

## Domain Configuration

1. Go to Project Settings > Domains
2. Add your custom domain
3. Update DNS records as instructed
4. SSL certificate is automatically provided

## Rollback

If you need to rollback to a previous deployment:
1. Go to Deployments tab
2. Find the deployment to rollback to
3. Click the three-dot menu
4. Select "Promote to Production"