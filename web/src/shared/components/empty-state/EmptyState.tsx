import { Box, Button, Typography } from '@mui/material';

import { AddCircleOutline } from '@mui/icons-material';
import React from 'react';

const EmptyState: React.FC<{
  headerText: string;
  subheaderText: string;
  addText: string;
  onAddButtonClick: () => void;
}> = ({ addText, headerText, subheaderText, onAddButtonClick }) => {
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100%',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          bgcolor: 'grey.200',
          borderRadius: '50%',
          width: '120px',
          height: '120px',
          mb: '2rem',
        }}
      >
        <AddCircleOutline sx={{ fontSize: '4rem' }} />
      </Box>
      <Typography variant="h5" align="center" gutterBottom>
        {headerText}
      </Typography>
      <Typography variant="subtitle1" align="center" gutterBottom>
        {subheaderText}
      </Typography>
      <Button variant="contained" onClick={onAddButtonClick}>
        {addText}
      </Button>
    </Box>
  );
};

export default EmptyState;
