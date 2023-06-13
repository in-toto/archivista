/* eslint-disable @typescript-eslint/no-unsafe-return */
import { act, renderHook } from '@testing-library/react-hooks';

import useArchivista from './useArchivista';
import { useBySubjectDigestLazyQuery } from '../../../generated/graphql';

jest.mock('../../../generated/graphql');

describe('useArchivista', () => {
  test('should handle loading state', async () => {
    const mockGetEnvelopes = jest.fn();
    (useBySubjectDigestLazyQuery as jest.Mock).mockReturnValue([mockGetEnvelopes]);

    const { result } = renderHook(() => useArchivista());

    act(() => {
      result.current[1]('validQuery');
    });

    expect(result.current[0].apiStatus.isLoading).toBe(true);

    await act(async () => {
      await mockGetEnvelopes({ variables: { digest: 'validQuery' } });
    });

    expect(result.current[0].apiStatus.isLoading).toBe(false);
  });

  test('should handle error state', async () => {
    const mockGetEnvelopes = jest.fn().mockRejectedValue(new Error('API Error'));
    (useBySubjectDigestLazyQuery as jest.Mock).mockReturnValue([mockGetEnvelopes]);

    const { result } = renderHook(() => useArchivista());

    act(() => {
      result.current[1]('validQuery');
    });

    expect(result.current[0].apiStatus.isLoading).toBe(true);

    await act(async () => {
      try {
        await mockGetEnvelopes({ variables: { digest: 'validQuery' } });
        // eslint-disable-next-line no-empty
      } catch (error) {}
    });

    expect(result.current[0].apiStatus.isLoading).toBe(false);

    expect(result.current[0].apiStatus.hasError).toBe(true);
  });
});
