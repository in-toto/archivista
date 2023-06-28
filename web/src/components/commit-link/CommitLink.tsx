import React from 'react';
import { Typography } from '@mui/material';

export type CommitLinkProps = {
  commit: string;
};

const CommitLink = ({ commit }: CommitLinkProps) => {
  return (
    <>
      <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
        Commit: {commit.split(':').pop()?.slice(0, 7)}
      </Typography>
    </>
  );
};

export default CommitLink;
