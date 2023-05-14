import { Box, Link, Typography } from '@mui/material';

import ArchivistaSvg from '../../svgs/archivista.inline.svg';
import { Link as GatsbyLink } from 'gatsby';
import JudgeSvg from '../../svgs/judge.inline.svg';
import React from 'react';
import WitnessSvg from '../../svgs/witness.inline.svg';

const Footer = () => {
  return (
    <Box sx={{ py: 3, bgcolor: 'grey.900', color: 'white', position: 'relative' }}>
      <Box
        sx={{
          position: 'absolute',
          top: '-69px',
          left: '16px',
          maxWidth: '200px',
          overflow: 'hidden',
          fontSize: '16px',
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        <ArchivistaSvg />
        <JudgeSvg />
        <WitnessSvg />
      </Box>
      <Box sx={{ display: 'flex', justifyContent: 'center' }}>
        <Link component={GatsbyLink} to="/privacy-policy" sx={{ color: '#BDBDBD', mr: 2 }}>
          Privacy Policy
        </Link>
        {/* <Link href="/about" sx={{ mr: 2 }}>
          About
        </Link>
        <Link href="/contact">Contact</Link> */}
      </Box>
      <Typography variant="body2" align="center" gutterBottom>
        Copyright © {new Date().getFullYear()} TestfySec LLC. All rights reserved.
      </Typography>
    </Box>
  );
};

export default Footer;
