import React from 'react';
import { Typography } from '@mui/material';
import { styled } from '@mui/material/styles';

const ellipsis = () => ({
  '0%': {
    content: '""',
  },
  '33%': {
    content: '"."',
  },
  '66%': {
    content: '".."',
  },
  '100%': {
    content: '"..."',
  },
});

const AnimatedEllipsis = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '200px',
  '& span': {
    // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
    animation: `${ellipsis} 1s infinite`,
  },
}));

const SearchEmptyState = () => {
  return (
    <AnimatedEllipsis>
      <Typography variant="h5" align="center">
        Start typing to search<span></span>
      </Typography>
    </AnimatedEllipsis>
  );
};

export default SearchEmptyState;
