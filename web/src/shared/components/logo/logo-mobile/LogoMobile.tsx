import { Box, Typography } from '@mui/material';

import { Link as GatsbyLink } from 'gatsby';
import LogoImage from '../logo-image/LogoImage';
import React from 'react';
import { useUiState } from '../../../contexts/ui-state-context/UiStateContext';

const LogoMobile: React.FC = () => {
  const { uiState } = useUiState();

  return (
    <>
      {!uiState.isSearchOpen && (
        <Box
          data-testid="logo"
          component={GatsbyLink}
          to="/"
          sx={{
            mr: 2,
            display: { xs: 'flex', md: 'none' },
            flexGrow: 1,
            fontWeight: 700,
            letterSpacing: '.3rem',
            color: 'inherit',
            textDecoration: 'none',
          }}
        >
          <LogoImage />
          <Typography variant="h5" component="div">
            Judge
          </Typography>
        </Box>
      )}
    </>
  );
};

export default LogoMobile;
