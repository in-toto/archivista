import * as React from 'react';

import { GitlabAttestor } from '../../models/Witness';
import { Jwt } from './Jwt';
import Typography from '@mui/material/Typography';

export function Gitlab(props: { data: GitlabAttestor }) {
  const attestation = props.data;

  return (
    <>
      <Jwt data={attestation.jwt}></Jwt>
      <Typography variant="body2">
        <b>CI config path:</b> {attestation.ciconfigpath}
      </Typography>
      <Typography variant="body2">
        <b>Job ID:</b> {attestation.jobid}
      </Typography>
      <Typography variant="body2">
        <b>Job image:</b> {attestation.jobimage}
      </Typography>
      <Typography variant="body2">
        <b>Job name:</b> {attestation.jobname}
      </Typography>
      <Typography variant="body2">
        <b>Job stage:</b> {attestation.jobstage}
      </Typography>
      <Typography variant="body2">
        <b>Job url:</b> {attestation.joburl}
      </Typography>
      <Typography variant="body2">
        <b>Pipeline ID:</b> {attestation.pipelineid}
      </Typography>
      <Typography variant="body2">
        <b>Pipeline url:</b> {attestation.pipelineurl}
      </Typography>
      <Typography variant="body2">
        <b>Project ID:</b> {attestation.projectid}
      </Typography>
      <Typography variant="body2">
        <b>Project url:</b> {attestation.projecturl}
      </Typography>
      <Typography variant="body2">
        <b>Runner ID:</b> {attestation.runnerid}
      </Typography>
      <Typography variant="body2">
        <b>CI host:</b> {attestation.cihost}
      </Typography>
    </>
  );
}
