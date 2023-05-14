export type StatementType = "https://in-toto.io/Statement/v0.1";
export type SlsaProvenanceType = "https://slsa.dev/provenance/v0.1";
export type WitnessPolicyType = "https://witness.testifysec.com/policy/v0.1";
export type PredicateType = SlsaProvenanceType | WitnessPolicyType;

export type DigestHashAlgorithms = "sha224" | "sha256" | "sha384" | "sha512";
export type Digests = Record<DigestHashAlgorithms, string>;

export function AllAttestations(): Record<string, string> {
  return {
    "Gitlab Job": "https://witness.testifysec.com/attestations/gitlab/v0.1",
    "SLSA Provenance": "https://slsa.dev/provenance/v0.1",
    "Snyk SAST": "https://witness.testifysec.com/attestations/snyk/v0.1",
  };
}

/** See https://slsa.dev/provenance/v0.1 and https://github.com/in-toto/ITE/blob/master/ITE/6/README.md for the types below. **/
export interface Subject {
  name: string;
  digest: Digests;
}

export interface Statement<T> {
  type: StatementType;
  subject: Subject[];
  predicateType: PredicateType;
  predicate: T;
}

export interface SlsaProvenanceBuilder {
  id: string;
}

export interface SlsaProvenanceCompleteness {
  arguments?: boolean;
  environment?: boolean;
  materials?: boolean;
}

export interface SlsaProvenanceMetadata {
  buildInvocationId?: string;
  completeness?: SlsaProvenanceCompleteness;
  reproducible?: boolean;
  buildStartedOn?: string;
  buildFinishedOn?: string;
}

export interface SlsaProvenanceRecipe {
  type: string;
  definedInMaterial?: number;
  entryPoint?: string;
  arguments?: Record<string | number, unknown>;
  environment?: Record<string | number, unknown>;
}

export interface SlsaProvenanceMaterial {
  uri: string;
  digest: Digests;
}

export interface SlsaProvenancePredicate {
  builder: SlsaProvenanceBuilder;
  metadata?: SlsaProvenanceMetadata;
  recipe?: SlsaProvenanceRecipe;
  materials: SlsaProvenanceMaterial[];
}
