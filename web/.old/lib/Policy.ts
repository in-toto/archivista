import { formatISO } from "date-fns";

import { CertConstraint, Key } from "./InToto";
import { CertToBase64PEM, ChainOfTrust, GenerateCertId } from "./ChainOfTrust";
/**
 * This is based upon the in-toto layout spec defined at https://github.com/in-toto/docs/blob/master/in-toto-spec.md#43-file-formats-layout
 * With ITE-5, ITE-6, and ITE-7 (https://github.com/in-toto/ITE/pull/21).
 * Let's not be afraid to make changes where they make sense for our uses.  ITE-6 currently breaks the idea of in-toto layouts as defined in the official spec
 * currently, so don't try too hard to maintain backward compatibility where it doesn't make sense.
 * See the Witness Architecture document on HackMD for the plan for this.
 **/
export interface Policy {
  expires: string;
  /**
   * I'm re-using the key type from in-toto here.. but I think we could remove the `type` and `scheme` fields on the type and just detect key types and make
   * sane choices for the scheme.  We could make the fields optional and use them if provided and fall back on defaults later... punting.
   **/
  roots?: Record<string, Root>;
  keys?: Record<string, Key>;
  steps: Step[];
}

export interface Root {
  certificate: string;
  intermediates?: string[];
}

/**
 * Each step of a CI/CD process may have multiple attestations we want to instrument.  For instance a base Gitlab CI job attestation that records basic Gitlab
 * job information like image used, job name, the stage the job belongs to, environment vars, each command run for the job, etc.
 * If the TPM is expected to be mounted in the job runner can add a tpm attestation as well. If a specific compiler runs and we can instrument it that may be
 * another attestation for the step.
 * Right now I'm putting the allowed functionary as part of the step, so a functionary would be allowed to run each of the attestations for the step.
 **/
export interface Step {
  name: string;
  functionaries: Functionary[];
  attestations: Attestation[];
}

export type PublicKeyFunctionaryType = "https://witness.testifysec.com/functionaries/publickey/v0.1";
export type CertConstraintFunctionaryType = "https://witness.testifysec.com/functionaries/certconstraint/v0.1";
export type FunctionaryType = PublicKeyFunctionaryType | CertConstraintFunctionaryType;
export type Functionary = PublicKeyFunctionary | CertConstraintFunctionary;

export interface PublicKeyFunctionary {
  type: PublicKeyFunctionary;
  keyId: string;
}

export interface CertConstraintFunctionary {
  type: CertConstraintFunctionaryType;
  certConstraint: CertConstraint;
}

export interface Attestation {
  /** The predicateType we expect this attestation to have during instrumentation */
  predicate: string;
  /**
   * Base64 encoded rego policies that will be run against the predicate during verification of the policy.
   * Each policy must pass for the attestation to be considered verified.
   **/
  policies: string[];
}

export function NewPolicy(expires: Date, chainsOfTrust: ChainOfTrust[], steps: Step[]): Policy {
  const roots: Record<string, Root> = {};
  for (const chain of chainsOfTrust) {
    const keyId = GenerateCertId(chain.root);
    const encodedRoot = CertToBase64PEM(chain.root);
    const encodedIntermediates = chain.intermediates.map((intermediate) => CertToBase64PEM(intermediate));
    roots[keyId] = { certificate: encodedRoot, intermediates: encodedIntermediates };
  }

  // temporarily just set the step's functionaries to allow certs from any of the roots in the layout
  for (const step of steps) {
    step.functionaries = Object.keys(roots).map((rootId) => ({
      type: "https://witness.testifysec.com/functionaries/certconstraint/v0.1",
      certConstraint: {
        common_name: "*",
        dns_names: ["*"],
        emails: ["*"],
        organizations: ["*"],
        uris: ["*"],
        roots: [rootId],
      },
    }));
  }

  // for now we're going to just put the

  return {
    expires: formatISO(expires),
    steps: steps,
    roots: roots,
  };
}
