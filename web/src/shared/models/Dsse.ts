/**
 * https://github.com/secure-systems-lab/dsse/blob/master/envelope.md
 * The Dead Simple Signature Envelope, plus some additions for certificates (and timestamps eventually?)
 **/
export interface Envelope<T> {
  payload: T;
  payloadType: string;
  signatures: Signature[];
}

export interface Signature {
  keyid: string;
  sig: string;
  /** certificate and intermediates are not part of the DSSE spec, though there are discussions about it:
   * https://github.com/secure-systems-lab/dsse/issues/42
   *
   * and timestamps:
   * https://github.com/secure-systems-lab/dsse/issues/33
   *
   * We should put more thought into how we represent the leaf and certificate chains.  my thought is for certificate is to be the leaf certificate
   * base64 encode, and intermediates being the same for any necessary intermediates to build trust back to a root.  i'm not currently including the
   * root here since we have it in our policy, but the dsse team may have opinions on that.
   **/
  certificate?: string;
  intermediates?: string[];
  timestamps?: [
    {
      data?: string;
    }
  ];
}
