import * as AttestationModal from "../../components/attestation/Modal";
import * as React from "react";

import { AttestationCollection, Dsse, Subject } from "../../generated/graphql";

import { Attestation } from "../../lib/Witness";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CertModal from "../../components/search/CertModal";
import Chip from "@mui/material/Chip";
import { ConfigContext } from "../../App";
import { Envelope } from "../../lib/Dsse";
import Grid from "@mui/material/Grid";
import Item from "@mui/material/Grid";
import Link from "@mui/material/Link";
import { Paper } from "@mui/material";
import { Statement } from "../../lib/Attestations";
import Typography from "@mui/material/Typography";

type DsseDisplayProps = {
  commithash: linkedLabel;
  pipelineurl: linkedLabel;
  ciprovider: linkedLabel;
  repo: linkedLabel;
};

type linkedLabel = {
  link: string;
  label: string;
  hashRef: string;
};

function MakeDisplayProps(dsse: Dsse): DsseDisplayProps {
  const commitSubject = dsse.statement?.subjects?.edges?.find((subject) => {
    return subject?.node?.name.includes("https://witness.dev/attestations/git/v0.1/commithash");
  });

  const isGitLabProvider = (subject: Subject) => {
    return subject.name.includes("https://witness.dev/attestations/gitlab/v0.1/pipelineurl");
  };

  const isGitHubProvider = (subject: Subject) => {
    return subject.name.includes("https://witness.dev/attestations/github/v0.1/pipelineurl");
  };

  const pipelineSubject = dsse.statement?.subjects?.edges?.find((subject) => {
    return isGitLabProvider(subject?.node as Subject) || isGitHubProvider(subject?.node as Subject);
  });

  const isProjectSubject = (subject: Subject) => {
    return subject.name.includes("https://witness.dev/attestations/gitlab/v0.1/project");
  };

  const projectSubject = dsse.statement?.subjects?.edges?.find((subject) => {
    return isProjectSubject(subject?.node as Subject);
  });

  const projectUrl = (subject: Subject) => {
    return subject.name.replace("https://witness.dev/attestations/gitlab/v0.1/projecturl:", "");
  };

  if (commitSubject) {
    let commitLink = null;
    if (pipelineSubject) {
      if (commitSubject.node?.subjectDigests?.length == 0) {
        commitLink = "";
      } else {
        const commitHash = subjectToHash(commitSubject?.node?.name as string);
        commitLink = `${projectUrl(projectSubject?.node as Subject)}/-/commit/${commitHash}`;
      }
      return {
        commithash: {
          link: commitLink,
          label: commitSubject.node?.name.replace("https://witness.dev/attestations/git/v0.1/commithash:", ""),
          hashRef: subjectToHash(commitSubject.node?.name as string),
        },
      } as DsseDisplayProps;
    }
  }

  return {} as DsseDisplayProps;
}

export function DsseEnvelope(props: { result: Dsse }) {
  const [certModalOpen, setCertModalOpen] = React.useState(Boolean);
  const [attestationModalOpen, setAttestationModalOpen] = React.useState(Boolean);
  const [attestationData, setAttestationData] = React.useState<Envelope<AttestationCollection> | null>(null);
  const [modalData, setModalData] = React.useState<Attestation | null>(null);
  const config = React.useContext(ConfigContext);

  const handleCertModalOpen = () => {
    if (!attestationData) {
      const r = getAttestationData(props.result.gitoidSha256, config.archivist_svc);
      void r.then((data) => {
        setAttestationData(data);
        setCertModalOpen(true);
      });
    } else {
      setCertModalOpen(true);
    }
  };

  const handleAttestationModalOpen = () => {
    setAttestationModalOpen(true);
  };

  const handleCertModalClose = () => setCertModalOpen(false);
  const handleAttestationModalClose = () => setAttestationModalOpen(false);
  const getAttestationChip = (attestation: Attestation | null) => {
    if (!attestation) {
      return null;
    }

    const color = "default";
    const attestationType = attestation.type;

    const clickAction = () => {
      if (!attestationData?.payload) {
        const r = getAttestationData(props.result.gitoidSha256, config.archivist_svc);
        void r.then((data) => {
          setAttestationData(data);
          const raw = data.payload as unknown as string;
          const payload = Buffer.from(raw, "base64").toString("utf-8");
          const attestation = JSON.parse(payload) as AttestationCollection;
          console.log(attestation);
          const stmt = JSON.parse(payload) as Statement<AttestationCollection>;
          const match = stmt.predicate.attestations?.find((att) => {
            return att?.type === attestationType;
          });
          if (match) {
            setModalData(match as unknown as Attestation);
            handleAttestationModalOpen();
          }
        });
      } else {
        const p = attestationData.payload as unknown as string;
        const payload = Buffer.from(p, "base64").toString("utf-8");
        const stmt = JSON.parse(payload) as Statement<AttestationCollection>;
        const match = stmt.predicate.attestations?.find((att) => {
          return att?.type === attestationType;
        });
        if (match) {
          setModalData(match as unknown as Attestation);
          handleAttestationModalOpen();
          console.log(match);
        }
      }
    };

    if (attestationType.includes("gitlab")) {
      return <Chip avatar={<Avatar>G</Avatar>} label="GitLab" color={color} onClick={clickAction} />;
    }
    if (attestationType.includes("gcp-iit")) {
      return <Chip avatar={<Avatar>G</Avatar>} label="Google Cloud" color={color} onClick={clickAction} />;
    }
    if (attestationType === "https://witness.dev/attestations/git/v0.1") {
      return <Chip avatar={<Avatar>G</Avatar>} label="Git" color={color} onClick={clickAction} />;
    }
    if (attestationType.includes("material")) {
      return <Chip avatar={<Avatar>M</Avatar>} label="Materials" color={color} onClick={clickAction} />;
    }
    if (attestationType.includes("command-run")) {
      return <Chip avatar={<Avatar>C</Avatar>} label="Execution Trace" color={color} onClick={clickAction} />;
    }
    if (attestationType.includes("product")) {
      return <Chip avatar={<Avatar>P</Avatar>} label="Products" color={color} onClick={clickAction} />;
    }

    if (attestationType === "https://witness.dev/attestations/github/v0.1") {
      return <Chip avatar={<Avatar>G</Avatar>} label="GitHub" color={color} onClick={clickAction} />;
    }

    return <Chip avatar={<Avatar>?</Avatar>} label="Unknown" color={color} onClick={clickAction} />;
  };

  const dp = MakeDisplayProps(props.result);
  const { result } = props;
  return (
    <div>
      <Card sx={{ minWidth: 600 }}>
        <div key={result?.gitoidSha256}>
          <div>
            <CardContent>
              <Grid container spacing={0}>
                <Grid item xs={4}>
                  <Item>
                    <Typography sx={{ fontSize: 14 }} color="text.secondary">
                      <div>Step Name: {result?.statement?.attestationCollections?.name}</div>
                      <div>
                        {dp.commithash && (
                          <Link href={dp.commithash?.link} target="_blank">
                            Commit: {dp.commithash?.hashRef}
                          </Link>
                        )}
                      </div>
                    </Typography>
                  </Item>
                </Grid>
                <Grid item xs={8}>
                  <Item>
                    <Typography sx={{ fontSize: 14 }} color="text.secondary" align="right">
                      <div>Envelope: {result?.gitoidSha256?.slice(-10)}</div>
                    </Typography>
                  </Item>
                </Grid>
              </Grid>
            </CardContent>
          </div>
          <div>
            <CardContent>
              <Paper
                elevation={1}
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  flexWrap: "wrap",
                  p: 0.5,
                  m: 0.5,
                }}
              >
                <Box sx={{ display: "flex", alignItems: "center" }}>
                  <Typography variant="h6" component="div" justifyContent="center">
                    Attestations:
                  </Typography>
                </Box>

                {result?.statement?.attestationCollections?.attestations?.map((attestation) => (
                  <div key={attestation?.type} style={{ margin: "0.5em" }}>
                    {getAttestationChip(attestation as unknown as Attestation)}
                  </div>
                ))}
              </Paper>
            </CardContent>
          </div>
        </div>
        <CardActions>
          <Button href={`https://archivist.testifysec.localhost/download/${result?.gitoidSha256}`} size="small">
            Download
          </Button>
          <Button size="small" onClick={handleCertModalOpen}>
            View Certificate
          </Button>
          <Button size="small">Share</Button>
        </CardActions>
        <CertModal data={attestationData} open={certModalOpen} close={handleCertModalClose} />
        <AttestationModal.Modal data={modalData} open={attestationModalOpen} close={handleAttestationModalClose} />
      </Card>
      <br />
    </div>
  );
}

function getAttestationData(oid: string, archivistURL: string) {
  const url = `${archivistURL}/download/${oid}`;
  const response = fetch(url).then((response) => response.json());
  return response as Promise<Envelope<AttestationCollection>>;
}

function subjectToHash(s: string) {
  const parts = s.split(":");
  const subject = parts[parts.length - 1] as string;

  //check if subject is commit hash
  if (subject?.length === 40) {
    return subject;
  }

  //else hash it and return it
  void sha256(subject).then((hash) => {
    return hash;
  });

  return "";
}

async function sha256(s: string) {
  const utf8 = new TextEncoder().encode(s);
  const hashBuffer = await crypto.subtle.digest("SHA-256", utf8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map((bytes) => bytes.toString(16).padStart(2, "0")).join("");
  return hashHex;
}
