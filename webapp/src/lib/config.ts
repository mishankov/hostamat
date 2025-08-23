import { env } from '$env/dynamic/private';

export const Config = {
  files: {
    uploadUrl: `${env.API_BASE || 'http://localhost:1001'}/files/api`
  },
  sites: {
    createUrl: `${env.API_BASE || 'http://localhost:1001'}/sites/api`,
    listUrl: `${env.API_BASE || 'http://localhost:1001'}/sites/api`
  }
} as const;
