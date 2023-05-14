/* eslint-disable @typescript-eslint/no-empty-function */
import { Avatar, Box, IconButton, Menu, MenuItem, Tooltip, Typography } from '@mui/material';

import DarkMode from '../dark-mode/DarkMode';
import Logout from './Logout';
import React from 'react';
import Search from '../../shared/components/search-menu/SearchMenu';
import { useUiState } from '../../shared/contexts/ui-state-context/UiStateContext';
import { useUser } from '../../shared/contexts/user-context/UserContext';

export type ToolbarProps = {
  anchorElUser?: HTMLElement;
  setAnchorElUser: React.Dispatch<HTMLElement | undefined>;
};
/**
 * This is the main toolbar of the app
 *
 * @param {ToolbarProps} { anchorElUser, setAnchorElUser }
 * @returns
 */
const MainToolbar = ({ anchorElUser, setAnchorElUser }: ToolbarProps) => {
  const { uiState } = useUiState();
  const user = useUser();
  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(undefined);
  };

  const settings = ['Logout'];

  const getSettingItem = (setting: string) => {
    switch (setting) {
      case 'Logout':
        return <Logout message={setting} />;

      default:
        return <Typography textAlign="center">{setting}</Typography>;
    }
  };

  const settingsMenuButton = (
    <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
      <Avatar alt="User profile picture" src="https://picsum.photos/200" sx={{ width: 32, height: 32 }} />
    </IconButton>
  );
  return user?.username ? (
    <Box sx={{ display: 'flex', alignItems: 'center', flexGrow: 0, marginLeft: 'auto' }}>
      <Tooltip title="Open search">
        <>
          <Search position="right" />
        </>
      </Tooltip>
      <Tooltip title="Open settings">
        <>{!uiState.isSearchOpen && settingsMenuButton}</>
      </Tooltip>
      <Menu
        sx={{ mt: '45px' }}
        id="menu-appbar"
        anchorEl={anchorElUser}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        keepMounted
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        open={Boolean(anchorElUser)}
        onClose={handleCloseUserMenu}
      >
        {user.username &&
          settings.map((setting) => (
            <MenuItem key={setting} onClick={handleCloseUserMenu}>
              {getSettingItem(setting)}
            </MenuItem>
          ))}
        <MenuItem>
          <DarkMode />
        </MenuItem>
      </Menu>
    </Box>
  ) : (
    <></>
  );
};

export default MainToolbar;
