import { renderHook } from '@testing-library/react-hooks';
import useArchivistaByRepos from './useArchivistaByRepo';
import { useQuery } from '@apollo/client';

// Mock the useQuery hook
jest.mock('@apollo/client');

describe('useArchivistaByRepos', () => {
  test('should handle loading state', () => {
    (useQuery as jest.Mock).mockReturnValue({
      loading: true,
      data: null,
      error: undefined,
    } as ReturnType<typeof useQuery>);

    // Render the hook
    const { result, rerender } = renderHook(() => useArchivistaByRepos(['repo1', 'repo2']));

    expect(result.current[0].apiStatus.isLoading).toBe(true);

    rerender();

    expect(result.current[0].apiStatus.isLoading).toBe(true);
  });

  test('should handle error state', () => {
    (useQuery as jest.Mock).mockReturnValue({
      loading: false,
      data: null,
      error: new Error('API Error'),
    } as ReturnType<typeof useQuery>);

    const { result, rerender } = renderHook(() => useArchivistaByRepos(['repo1', 'repo2']));

    expect(result.current[0].apiStatus.isLoading).toBe(false);

    rerender();

    expect(result.current[0].apiStatus.isLoading).toBe(false);

    expect(result.current[0].apiStatus.hasError).toBe(true);
  });
});
