import fetchMock from 'jest-fetch-mock';
import { renderHook } from '@testing-library/react-hooks';
import useJudgeApi from './useJudgeApi';

// Mock the fetch function
const mockedFetch = fetchMock;
global.fetch = mockedFetch as typeof fetch;

describe('useJudgeApi', () => {
  beforeEach(() => {
    fetchMock.resetMocks();
  });

  test('should handle successful API response', async () => {
    const responseData = {
      data: {
        projects: [
          { id: 1, repoID: 'repo1', name: 'Project 1', projecturl: 'https://example.com/project1' },
          { id: 2, repoID: 'repo2', name: 'Project 2', projecturl: 'https://example.com/project2' },
        ],
      },
    };

    fetchMock.mockResponseOnce(JSON.stringify(responseData));

    // Render the hook
    const { result, waitForNextUpdate } = renderHook(() => useJudgeApi());

    // Set the query
    result.current[1]('searchQuery');

    // Check the loading state
    expect(result.current[0].apiStatus.isLoading).toBe(true);

    // Wait for the next update, indicating the API call is finished
    await waitForNextUpdate();

    // Check the loading state
    expect(result.current[0].apiStatus.isLoading).toBe(false);

    // Check the results
    expect(result.current[0].results).toEqual(responseData.data.projects);

    // Check the query
    expect(result.current[0].query).toBe('searchQuery');
  });

  test('should handle API error', async () => {
    const errorMessage = 'API Error';

    fetchMock.mockRejectOnce(new Error(errorMessage));

    // Render the hook
    const { result, waitForNextUpdate } = renderHook(() => useJudgeApi());

    // Set the query
    result.current[1]('searchQuery');

    // Check the loading state
    expect(result.current[0].apiStatus.isLoading).toBe(true);

    // Wait for the next update, indicating the API call is finished
    await waitForNextUpdate();

    // Check the loading state
    expect(result.current[0].apiStatus.isLoading).toBe(false);

    // Check the error state
    expect(result.current[0].apiStatus.hasError).toBe(true);
  });
});
