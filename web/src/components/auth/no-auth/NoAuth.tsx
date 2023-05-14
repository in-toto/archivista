import * as React from 'react';

import { Box, Container, Typography } from '@mui/material';

import useContent from '../../../shared/hooks/use-content/useContent';

const NoAuth: React.FC = () => {
  const { noAuthYaml } = useContent();
  return (
    <Container maxWidth="md">
      <Box sx={{ py: 8 }}>
        <Typography variant="h2" sx={{ mb: 2 }}>
          {noAuthYaml.noAuthHeading}
        </Typography>
        <Typography variant="body1" sx={{ mb: 4 }}>
          {noAuthYaml.noAuthSubheading}
        </Typography>
      </Box>
    </Container>
  );
};
export default NoAuth;
