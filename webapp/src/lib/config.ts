import { env } from '$env/dynamic/private';

export const Config = {
  files: {
    uploadUrl: `${env.API_BASE}/files/api/upload`
  }
};