import * as React from 'react';

import { Box, Button } from '@mui/material';

import { Link as GatsbyLink } from 'gatsby';
import { useUiState } from '../../shared/contexts/ui-state-context/UiStateContext';
import { useUser } from '../../shared/contexts/user-context/UserContext';

export type MainMenuDesktopProps = {
  pages: { label: string; link: string }[];
  setAnchorElNav: React.Dispatch<HTMLElement | undefined>;
};

const MainMenuDesktop = ({ pages, setAnchorElNav }: MainMenuDesktopProps) => {
  const { uiState } = useUiState();
  const user = useUser();

  const handleCloseNavMenu = () => {
    setAnchorElNav(undefined);
  };

  return (
    <Box sx={{ flexGrow: 0.5, display: { xs: 'none', md: uiState?.isSearchOpen ? 'none' : 'flex' } }}>
      {user.username &&
        pages.map((page) => (
          <Button key={page.label} onClick={handleCloseNavMenu} sx={{ my: 2, color: 'white', display: 'block' }} component={GatsbyLink} to={page.link}>
            {page.label}
          </Button>
        ))}
    </Box>
  );
};

export default MainMenuDesktop;
