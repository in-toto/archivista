/**
 * This is the current InToto layout spec with ITE-7.
 */
export interface InTotoLayout {
  expires: string;
  roots?: Record<string, Key>;
  intermediates: Record<string, Key>;
  keys?: Record<string, Key>;
  steps: Step[];
  inspect: Inspection[];
}

export type RsaKeyType = 'rsa';
export type Ed25519KeyType = 'ed25519';
export type EcdsaKeyType = 'ecdsa';
export type KeyType = RsaKeyType | Ed25519KeyType | EcdsaKeyType;

export type RsaSsaPssSha256Scheme = 'rsassa-pss-sha256';
export type EcdsaSha2Nistp224Scheme = 'ecdsa-sha2-nistp224';
export type EcdsaSha2Nistp256Scheme = 'ecdsa-sha2-nistp256';
export type EcdsaSha2Nistp384Scheme = 'ecdsa-sha2-nistp384';
export type EcdsaSha2Nistp521Scheme = 'ecdsa-sha2-nistp521';
export type EcdsaScheme = EcdsaSha2Nistp521Scheme | EcdsaSha2Nistp384Scheme | EcdsaSha2Nistp256Scheme | EcdsaSha2Nistp224Scheme;
export type Ed25519Scheme = 'ed25519';
export type KeyScheme = RsaSsaPssSha256Scheme | EcdsaScheme | Ed25519Scheme;

export interface Key {
  keytype: KeyType;
  scheme: KeyScheme;
  keyval: KeyVal;
}

export interface KeyVal {
  public: string;
  certificate?: string;
}

/**
 * Do we even want to use this in witness policies?
 * https://github.com/in-toto/docs/blob/master/in-toto-spec.md#433-artifact-rules
 **/
export type MatchArtifactRule = ['MATCH', string, 'WITH', 'MATERIALS' | 'PRODUCTS', 'FROM', string];
export type MatchInArtifactRule = ['MATCH', string, 'IN', string, 'WITH', 'MATERIALS' | 'PRODUCTS', 'IN', string, 'FROM', string];
export type CreateArtifactRule = ['CREATE', string];
export type DeleteArtifactRule = ['DELETE', string];
export type ModifyArtifactRule = ['MODIFY', string];
export type AllowArtifactRule = ['ALLOW', string];
export type RequireArtifactRule = ['REQUIRE', string];
export type DisallowArtifactRule = ['DISALLOW', string];
export type ArtifactRule =
  | MatchArtifactRule
  | MatchInArtifactRule
  | CreateArtifactRule
  | DeleteArtifactRule
  | ModifyArtifactRule
  | AllowArtifactRule
  | RequireArtifactRule
  | DisallowArtifactRule;

export interface Step {
  name: string;
  threshold?: number;
  expected_materials: ArtifactRule[];
  expected_products: ArtifactRule[];
  expected_command: string[];
  pubkeys?: string[];
  cert_constraints?: CertConstraint[];
}

export interface CertConstraint {
  common_name: string;
  dns_names: string[];
  emails: string[];
  organizations: string[];
  roots: string[];
  uris: string[];
}

export interface Inspection {
  name: string;
  expected_materials: string[];
  expected_products: string[];
  run: string[];
}

export interface Statement<T> {
  _type: 'https://in-toto.io/Statement/v0.1';
  predicate: T;
  predicateType: string;
  subject: Subject[];
}

export interface Subject {
  name: string;
  digest: Record<string, string>;
}
