import * as React from 'react';

import type { HeadFC } from 'gatsby';
import Layout from '../layouts/DefaultLayout';
import RepoList from '../components/repo-list/RepoList';

const DashboardPage = () => {
  const body = <RepoList />;

  return <Layout>{body}</Layout>;
};

export default DashboardPage;

export const Head: HeadFC = () => <title>Dashboard - Judge Platform</title>;
