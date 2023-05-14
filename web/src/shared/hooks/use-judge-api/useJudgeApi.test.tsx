import { act, renderHook } from '@testing-library/react-hooks';

import useJudgeApi from './useJudgeApi';

describe('useJudgeApi', () => {
  test('returns the initial state', () => {
    const { result } = renderHook(() => useJudgeApi());

    expect(result.current[0]).toEqual({
      apiStatus: { isLoading: false, hasError: false },
      results: [],
      query: '',
    });
  });

  test('updates the query state', () => {
    const { result } = renderHook(() => useJudgeApi());

    act(() => {
      result.current[1]('new query');
    });

    expect(result.current[0].query).toBe('new query');
  });

  test('fetches results when query length is at least 3', async () => {
    const fetchMock = jest.spyOn(window, 'fetch');
    fetchMock.mockResolvedValue({
      json: () => [{ id: 1, name: 'result 1' }],
    } as unknown as Response);

    const { result, waitForNextUpdate } = renderHook(() => useJudgeApi());
    const [, setQuery] = result.current;

    act(() => {
      setQuery('abc');
    });

    await waitForNextUpdate();

    expect(result.current[0].results).toEqual([{ id: 1, name: 'result 1' }]);
    expect(fetchMock).toHaveBeenCalledWith('judge-api/query?query=abc', {
      method: 'GET',
      credentials: 'include',
    });
  });
});
