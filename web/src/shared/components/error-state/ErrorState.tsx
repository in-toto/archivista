import { Box, Link, Typography } from '@mui/material';

import React from 'react';
import { styled } from '@mui/material/styles';

const RootBox = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '200px',
  flexDirection: 'column',
  padding: theme.spacing(2),
  backgroundColor: theme.palette.error.light,
}));

const RobotPre = styled('pre')(({ theme }) => ({
  fontSize: '3rem',
  marginBottom: theme.spacing(2),
}));

const ErrorState = () => {
  return (
    <RootBox>
      <Typography variant="h5" align="center" gutterBottom>
        Oops! Something went wrong.
      </Typography>
      <Typography variant="subtitle1" align="center" gutterBottom>
        We apologize for the inconvenience. Please try again later or contact <Link href="mailto:support@testifysec.com">support@testifysec.com</Link> for
        assistance.
      </Typography>
      <RobotPre>
        {`
       /\\_/\\  
      / o o \\ 
     (   "   )
      \\~(*)~/ 
       - ^ -  
`}
      </RobotPre>
    </RootBox>
  );
};

export default ErrorState;
