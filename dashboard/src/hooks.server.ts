import { redirect } from '@sveltejs/kit';
import { validateAuthToken } from '$lib/auth';

const publicRoutes = ['/login'];

function isApiRoute(pathname: string): boolean {
	return pathname.startsWith('/api/');
}

export async function handle({ event, resolve }) {
	const authToken = event.cookies.get('dashboard-auth');
	const isAuthenticated = authToken && validateAuthToken(authToken);

	console.log('Auth check - path:', event.url.pathname, 'authenticated:', isAuthenticated);

	event.locals.user = isAuthenticated ? { authenticated: true } : null;

	// Don't apply authentication to API routes
	if (isApiRoute(event.url.pathname)) {
		return resolve(event);
	}

	// Redirect to login if not authenticated and trying to access protected routes
	if (!isAuthenticated && !publicRoutes.includes(event.url.pathname)) {
		console.log('Redirecting to login - not authenticated');
		throw redirect(302, '/login');
	}

	// Redirect to dashboard if authenticated and trying to access login
	if (isAuthenticated && event.url.pathname === '/login') {
		console.log('Redirecting to dashboard - already authenticated');
		throw redirect(302, '/');
	}

	return resolve(event);
}
