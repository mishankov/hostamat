import { json } from '@sveltejs/kit';

export async function POST({ cookies }) {
	cookies.delete('dashboard-auth', { path: '/' });
	return json({ success: true });
}
