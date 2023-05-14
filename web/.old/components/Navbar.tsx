import * as React from "react";

import { AppBar, Box, IconButton, Link, Menu, MenuItem, Toolbar, Typography } from "@mui/material";

import { AccountCircle } from "@mui/icons-material";
import Logo from "../../assets/logo.png";
import { Link as RouterLink } from "react-router-dom";

export interface Props {
  onLoginClick: () => void;
  onLogoutClick: () => void;
  userName: string | null;
}

export function Navbar({ onLoginClick, userName, onLogoutClick }: Props): JSX.Element {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <AppBar position="sticky">
      <Toolbar>
        <img src={Logo} height="50px" />
        <Typography variant="h4" component="div">
          Judge
        </Typography>
        <Box sx={{ pl: 1, flexGrow: 3, display: "flex", justifyContent: "flex-end" }}>
          <Link component={RouterLink} to="/" color="inherit" underline="none">
            Dashboard
          </Link>
        </Box>
        <Box sx={{ pl: 1, flexGrow: 1, display: "flex", justifyContent: "flex-end" }}>
          <Link component={RouterLink} to="/search" color="inherit" underline="none">
            Search
          </Link>
        </Box>
        <Box sx={{ pl: 1, flexGrow: 1, display: "flex", justifyContent: "flex-end" }}>
          <Link component={RouterLink} to="/sources" color="inherit" underline="none">
            Sources
          </Link>
        </Box>

        <Box sx={{ pl: 1, flexGrow: 1, display: "flex", justifyContent: "flex-end" }}>
          <Link component={RouterLink} to="/endpoints" color="inherit" underline="none">
            Endpoints
          </Link>
        </Box>

        <Box sx={{ pl: 1, flexGrow: 1, display: "flex", justifyContent: "flex-end" }}>
          <Link component={RouterLink} to="/polic" color="inherit" underline="none">
            Policies
          </Link>
        </Box>

        <Box sx={{ pl: 1, flexGrow: 3, display: "flex", justifyContent: "flex-end" }}>
          {userName === null ? (
            <Link component="button" color="inherit" variant="button" onClick={onLoginClick}>
              Sign In
            </Link>
          ) : (
            <div>
              {userName}
              <IconButton
                size="large"
                aria-label="account of current user"
                aria-controls="menu-appbar"
                aria-haspopup="true"
                onClick={handleMenu}
                color="inherit"
              >
                <AccountCircle />
              </IconButton>
              <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: "top",
                  horizontal: "right",
                }}
                keepMounted
                transformOrigin={{
                  vertical: "top",
                  horizontal: "right",
                }}
                open={Boolean(anchorEl)}
                onClose={handleClose}
              >
                {/* <MenuItem onClick={handleClose}>Profile</MenuItem> */}
                <MenuItem onClick={onLogoutClick}>Logout</MenuItem>
              </Menu>
            </div>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
}
