import RepoList from './RepoList';
import React from 'react';
import { render } from '@testing-library/react';

describe('RepoList', () => {
  it('should render without error', () => {
    render(<RepoList />);
  });
});
