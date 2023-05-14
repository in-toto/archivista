import * as React from "react";

import {
  Attestation,
  AwsAttestor,
  CommandRunAttestor,
  EnvironmentAttestor,
  GcpIitAttestor,
  GitAttestor,
  GitlabAttestor,
  JwtAttestor,
  MaterialAttestor,
  MavenAttestor,
  OciAttestor,
  ProductAttestor,
} from "../../lib/Witness";
import { Aws, CommandRun, Environment, GcpIit, Git, Gitlab, Jwt, Material, Maven, Oci, Product } from "./";

import Box from "@mui/material/Box";
import { Button } from "@mui/material";
import { DialogActions } from "@mui/material";
import { default as MuiModal } from "@mui/material/Modal";
import Typography from "@mui/material/Typography";

const style = {
  position: "absolute" as const,
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: "60%",
  height: "60%",
  display: "block",
  overflow: "scroll",
  bgcolor: "background.paper",
  boxShadow: 10,
  p: 4,
  outline: 0,
  fontSize: "0.875rem",
  ul: {
    margin: "0",
  },
};

interface AttestationModalProps {
  data: Attestation | null;
  open: boolean;
  close: () => void;
}

export function Modal(props: AttestationModalProps) {
  if (!props.open) {
    return null;
  }

  if (!props.data) {
    return (
      <MuiModal open={props.open} onClose={props.close}>
        <Box sx={style}>
          <Typography variant="h6" component="div">
            No attestation found
          </Typography>
          <DialogActions>
            <Button onClick={props.close}>Close</Button>
          </DialogActions>
        </Box>
      </MuiModal>
    );
  }

  // Switch for different attestation types will render different component
  const getAttestation = () => {
    if (!props.data) {
      return null;
    }

    console.log(props.data.type);

    switch (props.data.type) {
      case "https://witness.dev/attestations/gitlab/v0.1":
        return <Gitlab data={props.data.attestation as GitlabAttestor} />;
      // case "https://witness.dev/attestations/github/v0.1":
      //   return <Github data={props.data.attestation as GithubAttestor} />;
      case "https://witness.dev/attestations/gcp-iit/v0.1":
        return <GcpIit data={props.data.attestation as GcpIitAttestor} />;
      case "https://witness.dev/attestations/git/v0.1":
        return <Git data={props.data.attestation as GitAttestor} />;
      case "https://witness.dev/attestations/material/v0.1":
        return <Material data={props.data.attestation as MaterialAttestor} />;
      case "https://witness.dev/attestations/command-run/v0.1":
        return <CommandRun attestation={props.data.attestation as CommandRunAttestor} />;
      case "https://witness.dev/attestations/product/v0.1":
        return <Product data={props.data.attestation as ProductAttestor} />;
      // case "https://witness.dev/attestations/sarif/v0.1":
      //   return <Sarif data={props.data.attestation as SarifAttestor} />;
      case "https://witness.dev/attestations/jwt/v0.1":
        return <Jwt data={props.data.attestation as JwtAttestor} />;
      case "https://witness.dev/attestations/aws/v0.1":
        return <Aws data={props.data.attestation as AwsAttestor} />;
      case "https://witness.dev/attestations/environment/v0.1":
        return <Environment data={props.data.attestation as EnvironmentAttestor} />;
      case "https://witness.dev/attestations/oci/v0.1":
        return <Oci data={props.data.attestation as OciAttestor} />;
      case "https://witness.dev/attestations/maven/v0.1":
        return <Maven data={props.data.attestation as MavenAttestor} />;

      default:
        //pretty print the raw json
        return (
          <pre>
            <code>{JSON.stringify(props.data, null, 2)}</code>
          </pre>
        );
    }
  };

  return (
    <div>
      <MuiModal open={props.open} aria-labelledby="modal-modal-title" aria-describedby="modal-modal-description">
        <Box sx={style}>
          <Typography id="modal-modal-title" variant="h6" component="h2" style={{ border: "1px solid #dbdbdb", padding: "5px" }}>
            {props.data.type}
          </Typography>
          <Typography id="modal-modal-description" sx={{ mt: 2 }}></Typography>
          <div style={{ border: "1px solid #dbdbdb", padding: "5px" }}>{getAttestation()}</div>
          <DialogActions>
            <Button onClick={props.close}>Close</Button>
          </DialogActions>
        </Box>
      </MuiModal>
    </div>
  );
}
