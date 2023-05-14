import * as React from 'react';

import { GitAttestor } from '../../models/Witness';
import { NestedListItem } from '../nested-list-item/NestedListItem';
import { Typography } from '@mui/material';

export function Git(props: { data: GitAttestor }) {
  const attestation = props.data;
  const statusText = attestation.status === undefined ? 'No untracked or modified files' : '';
  const listItems: JSX.Element[] = [];
  if (attestation.status !== undefined) {
    for (const key in attestation.status) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      listItems.push(<NestedListItem key={key} k={key} v={(attestation.status[key] as any) ?? {}} />);
    }
  }

  return (
    <>
      <Typography variant="body2">
        <b>Commithash:</b> {attestation.commithash}
      </Typography>
      <Typography variant="body2">
        <b>Status:</b> {statusText}
      </Typography>
      {listItems.length > 0 ? <ul>{listItems}</ul> : null}
    </>
  );
}
