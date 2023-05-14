import * as React from 'react';

import type { HeadFC, PageProps } from 'gatsby';

import Layout from '../layouts/DefaultLayout';
import { Sources } from '../components/sources/sources';

const SourcesPages: React.FC<PageProps> = () => {
  const body = <Sources />;
  return <Layout>{body}</Layout>;
};

export default SourcesPages;

export const Head: HeadFC = () => <title>Sources - Judge Platform</title>;
