import * as React from 'react';

import { AppBar, Container, Toolbar } from '@mui/material';

import LogoDesktop from '../../shared/components/logo/logo-desktop/LogoDesktop';
import LogoMobile from '../../shared/components/logo/logo-mobile/LogoMobile';
import MainMenuDesktop from './MainMenuDesktop';
import MainMenuMobile from './MainMenuMobile';
import MainToolbar from './MainToolbar';
import { useState } from 'react';

/**
 * This is the global AppBarComponent for the app
 * It is displayed in the header area of the application.
 *
 * @returns
 */
const AppBarComponent: React.FC = () => {
  const [anchorElNav, setAnchorElNav] = useState<undefined | HTMLElement>(undefined);
  const [anchorElUser, setAnchorElUser] = useState<undefined | HTMLElement>(undefined);

  const pages = [
    // { label: 'Onboarding', link: '/onboarding' },
    // { label: 'Sources', link: '/sources' },
    // { label: 'Endpoints', link: '/endpoints' },
    // { label: 'Policies', link: '/policies' },
  ] as any;

  return (
    <AppBar position="static">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <LogoDesktop />
          <MainMenuMobile pages={pages} anchorElNav={anchorElNav} setAnchorElNav={setAnchorElNav} />
          <LogoMobile />
          <MainMenuDesktop pages={pages} setAnchorElNav={setAnchorElNav} />
          <MainToolbar setAnchorElUser={setAnchorElUser} anchorElUser={anchorElUser} />
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default AppBarComponent;
