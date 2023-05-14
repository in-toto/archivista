import * as forge from "node-forge";

export interface ChainOfTrust {
  id: string;
  root: forge.pki.Certificate;
  intermediates: forge.pki.Certificate[];
}

export function ParseFromPEM(pem: string): forge.pki.Certificate {
  /** This can throw errors... we will need to handle them **/
  return forge.pki.certificateFromPem(pem);
}

export function FriendlyString(cert: forge.pki.Certificate): string {
  let friendlyString = "";
  for (const subject of cert.subject.attributes) {
    if (subject.shortName !== undefined && typeof subject.value === "string") {
      friendlyString += `${subject.shortName}=${subject.value} `;
    }
  }

  return friendlyString;
}

export function GenerateCertId(cert: forge.pki.Certificate): string {
  return GeneratePubKeyId(cert.publicKey);
}

export function GeneratePubKeyId(pubKey: forge.pki.PublicKey): string {
  const pem = forge.pki.publicKeyToPem(pubKey).replace(/[\r]+/g, "");
  const keyDigest = forge.md.sha256.create();
  keyDigest.update(pem);
  const hexxed = keyDigest.digest().toHex();
  return forge.util.encode64(hexxed);
}

export function CertToBase64PEM(cert: forge.pki.Certificate): string {
  return forge.util.encode64(forge.pki.certificateToPem(cert));
}
