import { env } from '$env/dynamic/private';

export const Config = {
  files: {
    uploadUrl: `${env.API_BASE || 'http://localhost:1001'}/api/files`
  },
  sites: {
    createUrl: `${env.API_BASE || 'http://localhost:1001'}/api/sites`,
    listUrl: `${env.API_BASE || 'http://localhost:1001'}/api/sites`
  }
} as const;
