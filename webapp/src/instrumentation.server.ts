import * as Sentry from '@sentry/sveltekit';

Sentry.init({
	dsn: 'https://5aa1745a65b3bea6b8ae57db7c22717a@o1060431.ingest.us.sentry.io/4510556742680576',

	tracesSampleRate: 1.0,

	// Enable logs to be sent to Sentry
	enableLogs: true

	// uncomment the line below to enable Spotlight (https://spotlightjs.com)
	// spotlight: import.meta.env.DEV,
});
