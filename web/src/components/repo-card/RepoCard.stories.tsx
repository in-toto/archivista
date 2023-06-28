import { Dsse } from '../../generated/graphql';
import React from 'react';
import RepoCard from './RepoCard';
import { storiesOf } from '@storybook/react';

const mockArchivistaResults: Dsse[] = [];

const mockRepo = {
  id: 0,
  repoID: 'My test repo',
  description: '',
  projecturl: 'https://github.com/testifysec/judge',
  name: 'Test Repo',
};

storiesOf('RepoCard', module).add('default', () => <RepoCard archivistaResultsByRepo={mockArchivistaResults} repo={mockRepo} />);
