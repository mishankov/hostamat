export type Site = {
	id: string;
	name: string;
	subdomain?: string;
	domain: string;
	createdAt?: string;
	updatedAt: string;
};

export type DomainType = 'generated' | 'custom';
export type DomainStatus = 'pending' | 'verification_failed' | 'verified' | 'active';

export type DomainDNSRecord = {
	type: 'TXT' | 'CNAME';
	name: string;
	value: string;
};

export type DomainInstructions = {
	ownership: DomainDNSRecord;
	routing: DomainDNSRecord[];
};

export type DomainRecord = {
	id: string;
	siteId: string;
	siteName: string;
	siteDomain: string;
	hostname: string;
	type: DomainType;
	status: DomainStatus;
	verificationToken?: string;
	lastVerificationError?: string;
	lastCheckedAt?: string;
	activatedAt?: string;
	createdAt: string;
	updatedAt: string;
	instructions?: DomainInstructions;
};

export type ApiKeyScope =
	| 'sites:read'
	| 'sites:write'
	| 'sites:upload'
	| 'domains:read'
	| 'domains:write';

export type ApiKeyState = 'active' | 'expired' | 'revoked';

export type ApiKeyRecord = {
	id: string;
	name: string;
	prefix?: string;
	keyPrefix?: string;
	identifier?: string;
	scopes: ApiKeyScope[];
	createdAt: string;
	updatedAt?: string;
	expiresAt: string | null;
	lastUsedAt: string | null;
	state?: ApiKeyState;
	revokedAt?: string | null;
};
