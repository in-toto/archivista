import { Box, Typography } from '@mui/material';

import { Link as GatsbyLink } from 'gatsby';
import LogoImage from '../logo-image/LogoImage';
import React from 'react';

const LogoDesktop: React.FC = () => {
  return (
    <Box
      component={GatsbyLink}
      to="/"
      sx={{
        mr: 2,
        display: { xs: 'none', md: 'flex' },
        fontWeight: 700,
        letterSpacing: '.3rem',
        color: 'inherit',
        textDecoration: 'none',
      }}
    >
      <LogoImage />
      <Typography variant="h6" component="div">
        Judge
      </Typography>
    </Box>
  );
};

export default LogoDesktop;
