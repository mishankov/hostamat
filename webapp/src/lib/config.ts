import { env } from '$env/dynamic/private';

export const Config = {
  files: {
    uploadUrl: `${env.API_BASE || 'http://localhost:1001'}/api/files`
  },
  sites: {
    createUrl: `${env.API_BASE || 'http://localhost:1001'}/api/sites`,
    listUrl: `${env.API_BASE || 'http://localhost:1001'}/api/sites`
  },
  uploadSessionsUrl: `${env.API_BASE || 'http://localhost:1001'}/api/uploadsessions`,
  auth: {
    loginUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/login`,
    logoutUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/logout`,
    registerUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/register`,
    changePasswordUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/change-password`,
    meUrl: `${env.API_BASE || 'http://localhost:1001'}/auth/me`
  }
} as const;
