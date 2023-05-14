import * as React from "react";
import * as asn1js from "asn1js";
import * as pkijs from "pkijs";

import { Envelope, Signature } from "../../lib/Dsse";

import { AttestationCollection } from "../../generated/graphql";
import { Box } from "@mui/system";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import { DialogActions } from "@mui/material";
import Typography from "@mui/material/Typography";

interface CertModalProps {
  data?: Envelope<AttestationCollection>;
  open: boolean;
  close: () => void;
}

interface SigningTime {
  values: [{ year: number; month: number; day: number; hour: number; minute: number; second: number }];
}

interface TimestampCertificate {
  issuer: TimestampIssuer;
}

interface TimestampIssuer {
  typesAndValues?: [
    {
      type?: string;
      value: {
        getValue(): string;
      };
    }
  ];
}

const CertModal = (props: CertModalProps) => {
  if (!props.open) {
    return null;
  }

  const cert = importCert(props.data);
  if (!cert) {
    return (
      <Dialog open={props.open} onClose={props.close}>
        <Box sx={{ p: 2 }}>
          <Typography variant="h6" component="div">
            Attestation does not include a certificate
          </Typography>
          <DialogActions>
            <Button onClick={props.close}>Close</Button>
          </DialogActions>
        </Box>
      </Dialog>
    );
  }

  return (
    <Dialog open={props.open}>
      <Box
        sx={{
          padding: 2,
        }}
      >
        <Typography variant="h6" component="div">
          {cert?.subjectAlternativeNames}
        </Typography>
        <Typography variant="body1" component="div">
          Issued: {cert?.issued.toLocaleString()}
          <br />
          Expires: {cert?.expires.toLocaleString()}
          <br />
          {cert?.timestamp && (
            <>
              Timestamp: {cert?.timestamp.toLocaleString()}
              <br />
              TimeStamped By: {cert?.timestampCN}
              <br />
            </>
          )}
        </Typography>
        <Typography variant="h6" component="div">
          Subject
        </Typography>
        {/* // TODO: to can fix some html semantics here  if we use a unorder list below and style it instead of relying on <br /> tags. this will make restyling and ADA a lot better*/}
        <Typography variant="body1" component="div">
          CN: {cert?.subject.commonName}
          <br />
          O: {cert?.subject.organization}
          <br />
          OU: {cert?.subject.organizationalUnit}
          <br />
          L: {cert?.subject.city}
          <br />
          ST: {cert?.subject.state}
          <br />
          C: {cert?.subject.countryCode}
        </Typography>
        <Typography variant="h6" component="div">
          Issuer
        </Typography>
        <Typography variant="body1" component="div">
          CN: {cert?.issuer.commonName}
          <br />
          O: {cert?.issuer.organization}
          <br />
          OU: {cert?.issuer.organizationalUnit}
          <br />
          L: {cert?.issuer.city}
          <br />
          ST: {cert?.issuer.state}
          <br />
          C: {cert?.issuer.countryCode}
        </Typography>
      </Box>

      <DialogActions>
        <Button>Download</Button>
        <Button onClick={props.close}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

type Cert = {
  timestamp?: Date;
  timestampCN?: string;
  issued: Date;
  expires: Date;
  serial: string;
  keySize: number;
  subject: {
    commonName: string;
    organization: string;
    organizationalUnit: string;
    countryCode: string;
    state: string;
    city: string;
  };
  issuer: {
    commonName: string;
    organization: string;
    organizationalUnit: string;
    countryCode: string;
    state: string;
    city: string;
  };
  subjectAlternativeNames: string[];
};

const importCert = (data: Envelope<AttestationCollection> | undefined) => {
  if (!data) {
    return undefined;
  }

  if (data?.signatures?.length === 0) {
    return undefined;
  }

  const sigs = data.signatures;
  if (!sigs) {
    return undefined;
  }

  let timestampDate = undefined;
  let timestampCertCN = undefined;

  if (sigs[0]?.timestamps) {
    const timestampB64Encoded = sigs[0]?.timestamps[0]?.data;
    if (!timestampB64Encoded) return;

    const timestampBytes = Buffer.from(timestampB64Encoded, "base64");
    const timestamp = new Uint8Array(timestampBytes);

    const cmsContentSimp = pkijs.ContentInfo.fromBER(timestamp);
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const cmsSignedSimp = new pkijs.SignedData({ schema: cmsContentSimp.content });

    //get datetime of signature
    const signingTime = cmsSignedSimp?.signerInfos[0]?.signedAttrs?.attributes?.find((attr) => attr.type === "1.2.840.113549.1.9.5") as unknown as SigningTime;
    if (!signingTime) return;

    const signingTimeValue = signingTime?.values[0];
    const [year, month, day, hour, minute, second] = [
      signingTimeValue?.year,
      signingTimeValue?.month,
      signingTimeValue?.day,
      signingTimeValue?.hour,
      signingTimeValue?.minute,
      signingTimeValue?.second,
    ];

    //make a date object using GMT time
    timestampDate = new Date(Date.UTC(year, month, day, hour + 1, minute, second));

    const timestampCertData = (cmsSignedSimp?.certificates?.[0] as unknown as TimestampCertificate)?.issuer?.typesAndValues;

    timestampCertData?.forEach((item) => {
      if (item?.type == "2.5.4.3") {
        timestampCertCN = item?.value?.getValue();
      }
    });
  }

  const sig = data.signatures[0] as Signature;
  if (!sig.certificate) {
    return null;
  }
  // TODO: atob is deprecated
  const c = atob(sig.certificate);
  const b64 = c.replace(/(-----(BEGIN|END) CERTIFICATE-----|[\n\r])/g, "");
  const der = Buffer.from(b64, "base64");
  const ber = new Uint8Array(der).buffer;
  const asn1 = asn1js.fromBER(ber);
  const parsed = new pkijs.Certificate({ schema: asn1.result });

  const altNameOID = "2.5.29.17";
  const commonNameOID = "2.5.4.3";
  const countryCodeOID = "2.5.4.6";
  const localityNameOID = "2.5.4.7";
  const stateOrProvinceNameOID = "2.5.4.8";
  const organizationOID = "2.5.4.10";
  const organizationalUnitOID = "2.5.4.11";

  let altNames = [] as string[];
  parsed.extensions?.forEach(function (ext) {
    switch (ext.extnID) {
      case altNameOID:
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call
        altNames = ext.parsedValue.altNames.map((n: { value: string }) => n.value);
        break;
    }
  });

  let subjectCommonName = "";
  let subjectCountryCode = "";
  let subjectLocalityName = "";
  let subjectStateOrProvinceName = "";
  let subjectOrganization = "";
  let subjectOrganizationalUnit = "";

  parsed.subject.typesAndValues?.forEach(function (typeAndValue) {
    switch (typeAndValue.type) {
      case commonNameOID:
        subjectCommonName = typeAndValue.value.valueBlock.value;
        break;
      case countryCodeOID:
        subjectCountryCode = typeAndValue.value.valueBlock.value;
        break;
      case localityNameOID:
        subjectLocalityName = typeAndValue.value.valueBlock.value;
        break;
      case stateOrProvinceNameOID:
        subjectStateOrProvinceName = typeAndValue.value.valueBlock.value;
        break;
      case organizationOID:
        subjectOrganization = typeAndValue.value.valueBlock.value;
        break;
      case organizationalUnitOID:
        subjectOrganizationalUnit = typeAndValue.value.valueBlock.value;
        break;
    }
  });

  const issuer = {
    commonName: "",
    countryCode: "",
    organization: "",
    organizationalUnit: "",
    state: "",
    city: "",
  };

  parsed.issuer.typesAndValues?.forEach(function (typeAndValue) {
    switch (typeAndValue.type) {
      case commonNameOID:
        issuer.commonName = typeAndValue.value.valueBlock.value;
        break;
      case countryCodeOID:
        issuer.countryCode = typeAndValue.value.valueBlock.value;
        break;
      case localityNameOID:
        issuer.city = typeAndValue.value.valueBlock.value;
        break;
      case stateOrProvinceNameOID:
        issuer.state = typeAndValue.value.valueBlock.value;
        break;
      case organizationOID:
        issuer.organization = typeAndValue.value.valueBlock.value;
        break;
      case organizationalUnitOID:
        issuer.organizationalUnit = typeAndValue.value.valueBlock.value;
        break;
    }
  });

  const cert: Cert = {
    timestamp: timestampDate,
    timestampCN: timestampCertCN,
    issued: parsed.notBefore.value,
    expires: parsed.notAfter.value,
    serial: parsed.serialNumber.valueBlock.toString(),
    keySize: parsed.subjectPublicKeyInfo.subjectPublicKey.valueBlock.valueHex.byteLength * 8,
    subject: {
      commonName: subjectCommonName,
      countryCode: subjectCountryCode,
      organization: subjectOrganization,
      organizationalUnit: subjectOrganizationalUnit,
      state: subjectStateOrProvinceName,
      city: subjectLocalityName,
    },
    issuer,
    subjectAlternativeNames: altNames,
  };

  return cert;
};

export default CertModal;
