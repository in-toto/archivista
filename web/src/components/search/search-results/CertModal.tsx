/* eslint-disable @typescript-eslint/restrict-template-expressions */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import * as React from 'react';
import * as pkijs from 'pkijs';

import { Box } from '@mui/system';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import { Envelope } from '../../../shared/models/Dsse';
import Typography from '@mui/material/Typography';
import { PeculiarCertificateViewer } from '@peculiar/certificates-viewer-react';
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';
import Collapse from '@mui/material/Collapse';
import Tooltip from '@mui/material/Tooltip';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import DoneIcon from '@mui/icons-material/Done';
import { green } from '@mui/material/colors';

import * as asn1js from 'asn1js';

interface CertModalProps {
  data?: CertificateInfo[];
  open: boolean;
  onClose: () => void; // Change this from "close" to "onClose"
}

function CertModal({ data, open, onClose }: CertModalProps) {
  const [expanded, setExpanded] = React.useState<boolean[]>([]);
  const [certs, setCerts] = React.useState<CertificateInfo[]>([]);

  React.useEffect(() => {
    if (open && data) {
      const cert = importCert(data as unknown as Envelope<any>);
      if (cert) {
        setCerts(cert);
        setExpanded(Array(cert.length).fill(false));
      }
    }
  }, [data, open]);

  const handleClose = () => {
    setCerts([]);
    onClose();
  };

  const handleCopy = (cert: string) => {
    void navigator.clipboard.writeText(cert);
  };

  const toggleCollapse = (index: number) => {
    setExpanded((prevExpanded) => prevExpanded.map((state, idx) => (idx === index ? !state : state)));
  };

  const renderCert = (cert: CertificateInfo, index: number) => {
    const leafCert = cert.leafCert;
    const leafCertDerBase64 = convertPemToDerBase64(leafCert);
    const parsedCert = new pkijs.Certificate({ schema: asn1js.fromBER(fromBase64(leafCertDerBase64)).result });

    const issuerOrg = parsedCert.issuer.typesAndValues.find((typeAndValue) => typeAndValue.type === '2.5.4.10')?.value.valueBlock.value || 'Unknown';
    const issuerCN = parsedCert.issuer.typesAndValues.find((typeAndValue) => typeAndValue.type === '2.5.4.3')?.value.valueBlock.value || 'Unknown';

    const signature = cert.signature;
    const truncatedSignature = signature.substring(0, 15) + '...';

    return (
      <React.Fragment key={index}>
        <Box display="flex" flexDirection="column" justifyContent="space-between">
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Box display="flex" alignItems="center">
              <Typography variant="body1" style={{ marginLeft: 8, fontWeight: 'bold' }}>
                Signature {index + 1}:
              </Typography>
              <Typography variant="body2" color="textSecondary" style={{ marginLeft: 4 }}>
                {truncatedSignature}
              </Typography>
            </Box>
            <Box display="flex" alignItems="center">
              {cert.timestamps.map((timestampInfo, timestampIndex) => (
                <Box display="flex" alignItems="center" style={{ marginLeft: 8 }} key={`timestamp-${timestampIndex}`}>
                  <DoneIcon fontSize="small" style={{ color: green[500] }} />
                  <Typography variant="body2" style={{ marginLeft: 4, color: green[500] }}>
                    {timestampInfo.date.toLocaleString()} ({timestampInfo.issuer})
                  </Typography>
                </Box>
              ))}
              <Tooltip title="Copy signature to clipboard">
                <IconButton onClick={() => handleCopy(signature)} style={{ marginLeft: 8 }}>
                  <FileCopyIcon />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
          <Box display="flex" alignItems="center" justifyContent="space-between" mt={1}>
            <Box display="flex" alignItems="center">
              <Button onClick={() => toggleCollapse(index)}>
                <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                  Issuer:
                </Typography>
              </Button>
              <Box display="flex" alignItems="center" marginLeft={4}>
                <Typography variant="body2">
                  {issuerOrg} - {issuerCN}
                </Typography>
              </Box>
            </Box>
            <Box display="flex" alignItems="center">
              <Tooltip title="Copy leaf certificate to clipboard">
                <IconButton onClick={() => handleCopy(leafCert)}>
                  <FileCopyIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title={expanded[index] ? 'Collapse' : 'Expand'}>
                <IconButton onClick={() => toggleCollapse(index)} style={{ marginLeft: 8 }}>
                  {expanded[index] ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
          <Collapse in={expanded[index]}>
            <PeculiarCertificateViewer certificate={leafCertDerBase64} />
          </Collapse>
        </Box>
      </React.Fragment>
    );
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" px={2} py={1}>
        <Typography variant="h6">Signatures</Typography>
        <IconButton onClick={handleClose}>
          <CloseIcon />
        </IconButton>
      </Box>
      {certs.map((cert, index) => renderCert(cert, index))}
    </Dialog>
  );
}

interface TimestampInfo {
  date: Date;
  issuer: string;
}

interface CertificateInfo {
  signature: string;
  intermediates: string[];
  leafCert: string;
  timestamps: TimestampInfo[];
}

const importCert = (data?: Envelope<any>): undefined | CertificateInfo[] => {
  if (!data || !data.signatures || data.signatures.length === 0) {
    return undefined;
  }

  const certs: CertificateInfo[] = [];

  data.signatures.forEach((sig, index) => {
    const intermediates: string[] = [];
    let leafCert = '';
    const timestamps: TimestampInfo[] = [];

    // Process signature
    const signature = sig.sig;

    // Process timestamps
    if (sig.timestamps) {
      sig.timestamps.forEach((timestampInfo, timestampIndex) => {
        const timestampData = timestampInfo.data;

        if (timestampData) {
          const timestamp = Uint8Array.from(atob(timestampData), (c) => c.charCodeAt(0));

          const cmsContentSimp = pkijs.ContentInfo.fromBER(timestamp.buffer);
          const cmsSignedSimp = new pkijs.SignedData({ schema: cmsContentSimp.content });

          const signingTime = cmsSignedSimp?.signerInfos?.[0]?.signedAttrs?.attributes.find((attr) => attr.type === '1.2.840.113549.1.9.5');

          if (signingTime) {
            const [year, month, day, hour, minute, second] = [
              signingTime?.values[0]?.year,
              signingTime?.values[0]?.month,
              signingTime?.values[0]?.day,
              signingTime?.values[0]?.hour,
              signingTime?.values[0]?.minute,
              signingTime?.values[0]?.second,
            ];

            const timestampDate = new Date(Date.UTC(year, month, day, hour, minute, second));
            const timestampIssuer = (cmsSignedSimp?.certificates?.[0] as any)?.issuer?.typesAndValues?.[0]?.value?.valueBlock?.value;

            timestamps.push({ date: timestampDate, issuer: timestampIssuer });
            console.log(`Timestamp ${timestampIndex + 1} for signature ${index + 1}: ${timestampDate}, Issuer: ${timestampIssuer}`);
          }

          const timestampCert = new pkijs.Certificate({ schema: cmsSignedSimp?.certificates?.[0].toSchema() });
          intermediates.push(timestampCert.toSchema().toBER(false).toString());
        }
      });
    }

    if (!sig.certificate) {
      console.log(`No leaf certificate found for signature ${index + 1}`);
      return null;
    }

    leafCert = sig.certificate;

    if (sig.intermediates) {
      sig.intermediates.forEach((intermediate) => {
        intermediates.push(intermediate);
      });
    }

    const cert = {
      signature: signature,
      intermediates: intermediates,
      leafCert: leafCert,
      timestamps: timestamps,
    };

    certs.push(cert);
  });

  return certs;
};

function convertPemToDerBase64(base64Pem: string): string {
  const pem = atob(base64Pem);
  const pemHeaderFooterRegex = /-----BEGIN CERTIFICATE-----|-----END CERTIFICATE-----/g;
  const base64Der = pem.replace(pemHeaderFooterRegex, '').replace(/\n/g, '');
  return base64Der;
}

function fromBase64(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
  const raw = atob(base64);
  const uint8Array = new Uint8Array(raw.length);
  for (let i = 0; i < raw.length; i++) {
    uint8Array[i] = raw.charCodeAt(i);
  }
  return uint8Array;
}

export default CertModal;
