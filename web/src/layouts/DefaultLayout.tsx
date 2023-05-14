import { Box, Container } from '@mui/material';

import AppBarComponent from '../components/app-bar/AppBarComponent';
import Auth from '../Auth';
import Footer from '../components/footer/Footer';
import React from 'react';
import SearchResults from '../components/search/search-results/SearchResults';
import { useUiState } from '../shared/contexts/ui-state-context/UiStateContext';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { uiState } = useUiState();

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBarComponent />
      <Container sx={{ marginTop: 8, marginBottom: 8, flex: 1 }}>
        <Auth>
          {uiState?.isSearchOpen && <SearchResults />}
          <Box sx={{ display: !uiState.isSearchOpen ? 'inherit' : 'none' }}>{children}</Box>
        </Auth>
      </Container>
      <footer>
        <Footer />
      </footer>
    </Box>
  );
};

export default Layout;
