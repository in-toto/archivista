import * as React from 'react';

// import { Box, Container, Typography } from '@mui/material';
import { HeadFC, PageProps } from 'gatsby';

import Layout from '../layouts/DefaultLayout';
import RepoList from '../components/repo-list/RepoList';

const IndexPage: React.FC<PageProps> = () => {
  const body = (
    <RepoList />
    // <Container maxWidth="md">
    //   <Box sx={{ py: 8 }}>
    //     <Typography variant="h2" sx={{ mb: 2 }}>
    //       Welcome home to Judge!
    //     </Typography>
    //   </Box>
    // </Container>
  );

  return <Layout>{body}</Layout>;
};

export default IndexPage;

export const Head: HeadFC = () => <title>Home - Judge Platform</title>;
