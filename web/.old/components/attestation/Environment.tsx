import * as React from "react";
import { EnvironmentAttestor } from "../../lib/Witness";
import { Typography } from "@mui/material";

export function Environment(props: { data: EnvironmentAttestor }) {
  const attestation = props.data;
  return (
    <>
      <Typography variant="body2">
        <b>OS:</b> {attestation.os}
      </Typography>
      <Typography variant="body2">
        <b>Hostname:</b> {attestation.hostname}
      </Typography>
      <Typography variant="body2">
        <b>Username:</b> {attestation.username}
      </Typography>
      <Typography variant="body2">
        <b>Variables:</b>
      </Typography>
      <ul>
        {Object.keys(attestation.variables).map((key) => (
          <li key={key}>
            <b>{key}:</b>
            {attestation.variables[key]}
          </li>
        ))}
      </ul>
    </>
  );
}
