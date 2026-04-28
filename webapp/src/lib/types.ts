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
