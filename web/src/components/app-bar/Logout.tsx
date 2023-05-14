import { Button, Typography } from '@mui/material';

import { Link } from 'gatsby';
import React from 'react';
import { useUser } from '../../shared/contexts/user-context/UserContext';

export type LogoutProps = {
  message?: string;
};

const Logout = ({ message }: LogoutProps) => {
  const { logoutUrl } = useUser();
  return (
    <>
      {logoutUrl && (
        <Button sx={{ textTransform: 'inherit' }} component={Link} to={logoutUrl}>
          <Typography textAlign="center">{message}</Typography>
        </Button>
      )}
    </>
  );
};

export default Logout;
