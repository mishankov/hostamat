I got this issue with this project. For every new site in create a new subdomain in caddy. When too many sites are created i got an error like this:

[lucky-jamboree-57738.hostamat.ru] Obtain: [lucky-jamboree-57738.hostamat.ru] creating new order: attempt 1: https://acme-v02.api.letsencrypt.org/acme/new-order: HTTP 429 urn:ietf:params:acme:error:rateLimited - too many certificates (50) already issued for "hostamat.ru" in the last 168h0m0s, retry after 2025-12-25 04:58:32 UTC: see https://letsencrypt.org/docs/rate-limits/#new-certificates-per-registered-domain (ca=https://acme-v02.api.letsencrypt.org/directory)

I am hitting rate limits. Solution is probably wildcard certificates, but I am not totaly sure. Help me fix this issue.
