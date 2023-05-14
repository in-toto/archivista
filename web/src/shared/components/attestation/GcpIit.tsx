import * as React from 'react';

import { GcpIitAttestor } from '../../models/Witness';
import { Jwt } from './';
import { Typography } from '@mui/material';

export function GcpIit(props: { data: GcpIitAttestor }) {
  const attestation = props.data;

  return (
    <>
      <Jwt data={attestation.jwt}></Jwt>
      <Typography variant="body2">
        <b>Project ID:</b> {attestation.project_id}
      </Typography>
      <Typography variant="body2">
        <b>Project Number:</b> {attestation.project_number}
      </Typography>
      <Typography variant="body2">
        <b>Zone:</b> {attestation.zone}
      </Typography>
      <Typography variant="body2">
        <b>Instance Id:</b> {attestation.instace_id}
      </Typography>
      <Typography variant="body2">
        <b>Instance Hostname:</b> {attestation.instance_hostname}
      </Typography>
      <Typography variant="body2">
        <b>Instance Creation Timestamp:</b> {attestation.instance_creation_timestamp}
      </Typography>
      <Typography variant="body2">
        <b>Instance Confidentiality:</b> {attestation.instance_confidentiality}
      </Typography>
      <Typography variant="body2">
        <b>License Id:</b> {attestation.licence_id}
      </Typography>
      <Typography variant="body2">
        <b>Cluster Name:</b> {attestation.cluster_name}
      </Typography>
      <Typography variant="body2">
        <b>Cluster UID:</b> {attestation.cluster_uid}
      </Typography>
      <Typography variant="body2">
        <b>Cluster Location:</b> {attestation.cluster_location}
      </Typography>
    </>
  );
}
