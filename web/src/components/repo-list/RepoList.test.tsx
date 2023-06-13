/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import '@testing-library/jest-dom';

import { MockedProvider } from '@apollo/client/testing';
import React from 'react';
import RepoList from './RepoList';
import { render } from '@testing-library/react';
import useArchivistaByRepos from '../../shared/hooks/use-archivista/useArchivistaByRepo';
import useJudgeApi from '../../shared/hooks/use-judge-api/useJudgeApi';

jest.mock('../../shared/components/skeleton-loader/SkeletonLoader');
jest.mock('../timeline/Timeline');
jest.mock('../../shared/components/error-state/ErrorState');
jest.mock('../../shared/hooks/use-judge-api/useJudgeApi');
jest.mock('../../shared/hooks/use-archivista/useArchivistaByRepo');

describe('RepoList', () => {
  beforeEach(() => {
    (useJudgeApi as jest.Mock).mockReturnValue([
      { apiStatus: { isLoading: false }, results: [] }, // Mocked judge state
      null, // Mocked setQuery function
    ]);

    (useArchivistaByRepos as jest.Mock).mockReturnValue([
      { apiStatus: { isLoading: false }, results: [] }, // Mocked archivista state
    ]);
  });
  // eslint-disable-next-line jest/expect-expect
  it('should render without error', () => {
    render(
      <MockedProvider>
        <RepoList />
      </MockedProvider>
    );
  });

  it('should show loading state when archivista loading', () => {
    (useArchivistaByRepos as jest.Mock).mockReturnValue([
      { apiStatus: { isLoading: true }, results: [] }, // Mocked archivista state
    ]);

    const { getByTestId } = render(
      <MockedProvider>
        <RepoList />
      </MockedProvider>
    );
    const skeletonLoader = getByTestId('SkeletonLoader');
    expect(skeletonLoader).toBeVisible();
  });

  it('should show loading state when judge-api loading', () => {
    (useJudgeApi as jest.Mock).mockReturnValue([
      { apiStatus: { isLoading: true }, results: [] }, // Mocked judge state
      null, // Mocked setQuery function
    ]);

    const { getByTestId } = render(
      <MockedProvider>
        <RepoList />
      </MockedProvider>
    );
    const skeletonLoader = getByTestId('SkeletonLoader');
    expect(skeletonLoader).toBeVisible();
  });

  it('should NOT show loading state when archivista and judge-api are not loading', () => {
    const { queryByTestId } = render(
      <MockedProvider>
        <RepoList />
      </MockedProvider>
    );
    const skeletonLoader = queryByTestId('SkeletonLoader');
    expect(skeletonLoader).toBeFalsy();
  });
});
