export interface LogEntry {
  attestation: {
    data: string;
  };
  body: string;
  integratedTime: number;
  logID: string;
  logIndex: number;
  verification: LogEntryVerification;
}

export interface LogEntryVerification {
  signedEntryTimestamp: string;
  inclusionProof: LogEntryInclusionProof;
}

export interface LogEntryInclusionProof {
  hashes: Record<number, string>;
  logIndex: number;
  rootHash: string;
  treeSize: number;
}

export type RekorKinds = Dssev001;

export interface BaseRekorKind {
  apiVersion: string;
  kind: string;
  spec: RekorKinds;
}

export interface DigestSet {
  algorithm: string;
  value: string;
}

export interface Dssev001Signature {
  intermediates?: string[];
  keyid: string;
  publicKey: string;
}

export interface Dssev001 {
  payloadHash: DigestSet;
  payloadType: string;
  signatures: Dssev001Signature[];
}

export class ApiService {
  #baseUrl: string;

  constructor(baseUrl: string) {
    this.#baseUrl = baseUrl;
  }

  async searchIndexByHash(hash: string): Promise<string[]> {
    const url = new URL('/api/v1/index/retrieve', this.#baseUrl);
    const response = await fetch(url.toString(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ hash: `sha256:${hash}` }),
    });

    return (await response.json()) as string[];
  }

  async searchIndexByString(search: string): Promise<string[]> {
    const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(search));
    const hexString = Array.from(new Uint8Array(digest))
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('');
    return this.searchIndexByHash(hexString);
  }

  async getLogEntryByUUID(uuid: string): Promise<LogEntry> {
    const url = new URL(`/api/v1/log/entries/${uuid}`, this.#baseUrl);
    const response = await fetch(url.toString(), {
      headers: {
        Accept: 'application/json',
      },
    });

    const record = (await response.json()) as Record<string, LogEntry>;
    return record[uuid];
  }

  getLogEntryUrl(uuid: string): string {
    return new URL(`/api/v1/log/entries/${uuid}`, this.#baseUrl).toString();
  }
}
