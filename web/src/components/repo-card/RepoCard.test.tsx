import '@testing-library/jest-dom';

import { Dsse } from '../../generated/graphql';
import React from 'react';
import RepoCard from './RepoCard';
import { render } from '@testing-library/react';

const mockArchivistaResults: Dsse[] = [];

const mockRepo = {
  id: 0,
  repoID: 'My test repo',
  description: '',
  projecturl: 'https://github.com/testifysec/judge',
  name: 'Test Repo',
};

describe('RepoCard', () => {
  test('renders the component with the provided data', () => {
    const { getByText } = render(<RepoCard archivistaResultsByRepo={mockArchivistaResults} repo={mockRepo} />);

    const repoName = getByText('Test Repo');

    expect(repoName).toBeInTheDocument();
  });
});
