# Hostamat Dashboard

A secure dashboard for monitoring Hostamat system metrics.

## Features

- **Authentication**: Simple token-based authentication
- **System Metrics**: Real-time statistics including:
  - Registered users count
  - Sites created count
  - Upload sessions count
- **Responsive Design**: Built with DaisyUI and Tailwind CSS
- **Secure**: Protected routes with cookie-based authentication

## Setup

1. Copy `.env` file and configure:

   ```bash
   cp .env .env.local
   ```

2. Update environment variables in `.env.local`:

   ```env
   # Database connection
   DATABASE_URL=postgresql://username:password@localhost:5432/hostamat

   # Dashboard authentication token
   DASHBOARD_AUTH_TOKEN=your-secret-token-here
   ```

3. Install dependencies:

   ```bash
   bun install
   ```

4. Start development server:
   ```bash
   bun run dev
   ```

## Usage

1. Navigate to `/login` and enter your authentication token
2. Access the dashboard at `/` to view system metrics
3. Use the logout button to end your session

## Database Schema

The dashboard queries the following tables:

- `sites` - Count of all created sites
- `uploadsessions` - Count of all upload sessions
- Distinct users from both tables for user count

## Development

- Built with SvelteKit
- Uses DaisyUI for components
- Tailwind CSS for styling
- PostgreSQL for data storage
