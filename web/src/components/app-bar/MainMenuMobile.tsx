import * as React from 'react';

import { Box, IconButton, Menu, MenuItem, Typography } from '@mui/material';

import { Link as GatsbyLink } from 'gatsby';
import MenuIcon from '@mui/icons-material/Menu';
import { useUiState } from '../../shared/contexts/ui-state-context/UiStateContext';
import { useUser } from '../../shared/contexts/user-context/UserContext';

export type MainMenuMobileProps = {
  anchorElNav?: HTMLElement;
  pages: { label: string; link: string }[];
  setAnchorElNav: React.Dispatch<HTMLElement | undefined>;
};

const MainMenuMobile = ({ anchorElNav, pages, setAnchorElNav }: MainMenuMobileProps) => {
  const { uiState } = useUiState();
  const user = useUser();

  const handleOpenNavMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElNav(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(undefined);
  };

  return (
    <Box sx={{ flexGrow: 1, display: { xs: uiState?.isSearchOpen ? 'none' : 'flex', md: 'none' } }}>
      {user.username && pages?.length > 0 && (
        <>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleOpenNavMenu}
            color="inherit"
          >
            <MenuIcon />
          </IconButton>
        </>
      )}
      <Menu
        id="menu-appbar"
        anchorEl={anchorElNav}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        keepMounted
        transformOrigin={{
          vertical: 'top',
          horizontal: 'left',
        }}
        open={Boolean(anchorElNav)}
        onClose={handleCloseNavMenu}
        sx={{
          display: { xs: 'block', md: 'none' },
        }}
      >
        {user.username &&
          pages.map((page) => (
            <MenuItem component={GatsbyLink} to={page.link} key={page.label} onClick={handleCloseNavMenu}>
              <Typography textAlign="center">{page.label}</Typography>
            </MenuItem>
          ))}
      </Menu>
    </Box>
  );
};

export default MainMenuMobile;
