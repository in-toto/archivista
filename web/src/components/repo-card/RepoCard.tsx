import { Box, Card, CardContent, Typography } from '@mui/material';

import { Dsse } from '../../generated/graphql';
import React from 'react';
import { Repo } from '../../shared/models/app-data-model';
import Timeline from '../timeline/Timeline';

export type RepoCardProps = {
  archivistaResultsByRepo: Dsse[];
  repo: Repo;
};

const RepoCard = ({ archivistaResultsByRepo, repo }: RepoCardProps) => {
  return (
    <>
      <Card sx={{ minWidth: 275, margin: '8px' }} key={repo?.id}>
        <CardContent>
          <Typography sx={{ fontSize: 20, fontWeight: 'bold', marginBottom: '8px' }}>{repo?.name}</Typography>
          <Typography sx={{ fontSize: 14, marginBottom: '8px' }}>{repo?.description}</Typography>
          <Box>
            <Timeline dsseArray={archivistaResultsByRepo} />
          </Box>
        </CardContent>
      </Card>
    </>
  );
};

export default RepoCard;
