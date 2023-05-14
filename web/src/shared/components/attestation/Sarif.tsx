import * as React from 'react';

import { NestedListItem } from '../nested-list-item/NestedListItem';
import { SarifAttestor } from '../../models/Witness';
import { Typography } from '@mui/material';

export function Sarif(props: { data: SarifAttestor }) {
  const attestation = props.data;
  return (
    <ul>
      <li>
        <Typography variant="body2">
          <b>Report:</b>
        </Typography>
      </li>
      <ul>
        <li>
          <Typography variant="body2">
            <b>Version:</b> {attestation.report.version}
          </Typography>
        </li>
        <li>
          <Typography variant="body2">
            <b>Schema:</b> {attestation.report.$schema}
          </Typography>
        </li>
        <li>
          <Typography variant="body2">
            <b>Runs:</b> {attestation.report.runs}
          </Typography>
        </li>
      </ul>
      <li key="reportFileName">
        <Typography variant="body2">
          <b>Report file name:</b> {attestation.reportFileName}
        </Typography>
      </li>
      <NestedListItem key="reportDigestSet" k="Report Digest Set" v={attestation.reportDigestSet} />
    </ul>
  );
}
