import { Card, CardContent, CardMedia, Skeleton, Typography } from '@mui/material';

import React from 'react';

const SkeletonLoader = () => {
  return (
    <>
      {[1, 2, 3].map((n) => (
        <Card key={n}>
          <CardMedia>
            <Skeleton variant="rectangular" width={151} height={151} />
          </CardMedia>
          <>
            <CardContent>
              <Typography gutterBottom variant="h5">
                <Skeleton />
              </Typography>
              <Typography>
                <Skeleton />
              </Typography>
              <Typography>
                <Skeleton />
              </Typography>
            </CardContent>
          </>
        </Card>
      ))}
    </>
  );
};

export default SkeletonLoader;
