import * as React from 'react';

import { MavenAttestor } from '../../models/Witness';
import { Typography } from '@mui/material';

export function Maven(props: { data: MavenAttestor }) {
  const attestation = props.data;

  return (
    <>
      <Typography variant="body2">
        <b>GroupId:</b> {attestation.groupid}
      </Typography>
      <Typography variant="body2">
        <b>ArtifactId:</b> {attestation.artifactid}
      </Typography>
      <Typography variant="body2">
        <b>Version:</b> {attestation.version}
      </Typography>
      <Typography variant="body2">
        <b>Project Name:</b> {attestation.projectname}
      </Typography>
      <Typography variant="body2">
        <b>Dependencies:</b>
      </Typography>
      <ul>
        {attestation.dependencies.map((dependency, index) => (
          <li key={dependency.groupid + index.toString()}>
            <Typography variant="body2">
              <b>GroupId</b> {dependency.groupid}
            </Typography>
            <Typography variant="body2">
              <b>ArtifactId</b> {dependency.artifactid}
            </Typography>
            <Typography variant="body2">
              <b>Version</b> {dependency.version}
            </Typography>
            <Typography variant="body2">
              <b>Scope</b> {dependency.scope}
            </Typography>
            <br></br>
          </li>
        ))}
      </ul>
    </>
  );
}
