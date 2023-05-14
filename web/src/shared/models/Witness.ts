export interface AttestationCollection {
  predicateType: 'https://witness.testifysec.com/attestation-collection/v0.1';
  name: string;
  attestations: Attestation[];
}

export interface Attestation {
  type: string;
  attestation: unknown;
}

export interface GitlabAttestor {
  jwt: JwtAttestor;
  ciconfigpath: string;
  jobid: string;
  jobimage: string;
  jobname: string;
  jobstage: string;
  joburl: string;
  pipelineid: string;
  pipelineurl: string;
  projectid: string;
  projecturl: string;
  runnerid: string;
  cihost: string;
}

export interface JwtAttestor {
  claims: { [key: string]: unknown };
  verifiedBy: VerificationInfo;
  jwksurl: string;
  token: string;
}

export interface VerificationInfo {
  jwksUrl: string;
  jwk: { [key: string]: string };
}

export interface GcpIitAttestor {
  jwt: JwtAttestor;
  project_id: string;
  project_number: string;
  zone: string;
  instace_id: string;
  instance_hostname: string;
  instance_creation_timestamp: string;
  instance_confidentiality: string;
  licence_id: string[];
  cluster_name: string;
  cluster_uid: string;
  cluster_location: string;
}

export interface GitStatus {
  staging: string;
  worktree: string;
}

export interface GitAttestor {
  commithash: string;
  status?: { [key: string]: GitStatus };
}

export type HashType = 'sha256' | 'sha1' | 'md5';

export type DigestSet = Record<HashType, string>;

export type MaterialAttestor = Record<string, DigestSet>;

export interface CommandRunAttestor {
  cmd: string[];
  exitcode: number;
  stdout?: string;
  stderr?: string;
  processes?: ProcessInfo[];
}

export interface ProcessInfo {
  program?: string;
  processid: number;
  parentpid: number;
  programdigest?: DigestSet;
  comm?: string;
  cmdline?: string;
  exedigest?: DigestSet;
  openedfiles?: Record<string, DigestSet>;
  environ?: string;
  specbypassisvuln?: boolean;
}

export type ProductAttestor = Record<string, Product>;

export interface Product {
  mime_type: string;
  digest: DigestSet;
}

export interface SarifAttestor {
  report: {
    version: string;
    $schema: string;
    runs: string;
  };
  reportFileName: string;
  reportDigestSet: DigestSet;
}

export interface AwsAttestor {
  devpayProductCodes: string;
  marketplaceProductCodes: string;
  availabilityZone: string;
  privateIp: string;
  version: string;
  region: string;
  instanceId: string;
  billingProducts: string;
  instanceType: string;
  accountId: string;
  pendingTime: string;
  imageId: string;
  kernalId: string;
  ramdiskId: string;
  architecture: string;
  rawiid: string;
  rawsig: string;
  publickey: string;
}

export interface EnvironmentAttestor {
  os: string;
  hostname: string;
  username: string;
  variables: Record<string, string>;
}

export interface OciAttestor {
  tardigest: DigestSet;
  manifest: Manifest[];
  imagetags: string[];
  diffids: DigestSet[];
  imageid: DigestSet;
  manifestraw: string;
}

export interface Manifest {
  Config: string;
  RepoTags: string[];
  Layers: string[];
}

export interface MavenAttestor {
  groupid: string;
  artifactid: string;
  version: string;
  projectname: string;
  dependencies: MavenDependency[];
}

export interface MavenDependency {
  groupid: string;
  artifactid: string;
  version: string;
  scope: string;
}

export interface GithubAttestor {
  claims: {
    actor: string;
    actor_id: string;
    aud: string;
    base_ref: string;
    event_name: string;
    exp: number;
    head_ref: string;
    iat: number;
    iss: string;
    job_workflow_ref: string;
    jti: string;
    nbf: number;
    ref: string;
    ref_type: string;
    repository: string;
    repository_id: string;
    repository_owner: string;
    repository_owner_id: string;
    repository_visibility: string;
    run_attempt: string;
    run_id: string;
    run_number: string;
    sha: string;
    sub: string;
    workflow: string;
  };
  verifiedBy?: {
    jwksUrl: string;
    jwk: {
      use: string;
      kty: string;
      alg: string;
      n: string;
      e: string;
      x5c: string[];
      x5t: string;
    };
  };

  ciconfigpath: string;
  pipelineid: string;
  pipelinename: string;
  pipelineurl: string;
  projecturl: string;
  runnerid: string;
  cihost: string;
  ciserverurl: string;
  runnerarch: string;
  runneros: string;
}

export class WitnessApiService {
  readonly #baseUrl: string;
  #token: string;

  constructor(baseUrl: string, token: string) {
    this.#baseUrl = baseUrl;
    this.#token = token;
  }

  async fetchTestEndpoint(): Promise<unknown> {
    const url = new URL('/token', this.#baseUrl);

    const response = await fetch(url.toString(), {
      headers: {
        Authorization: `Bearer ${this.#token}`,
      },
    });

    return response.json();
  }
}
