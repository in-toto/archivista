import { act, renderHook } from '@testing-library/react-hooks';

import <%= hookName %> from './<%= hookName %>';

describe('<%= hookName %> hook', () => {
  it('should initialize the api status correctly', () => {
    const { result } = renderHook(() => <%= hookName %>());

    expect(result.current[0]).toEqual({ isLoading: false, hasError: false });
  });

  it('should update isLoading correctly', () => {
    const { result } = renderHook(() => <%= hookName %>());

    act(() => {
      result.current[1](true);
    });

    expect(result.current[0]).toEqual({ isLoading: true, hasError: false });
  });

  it('should update hasError correctly', () => {
    const { result } = renderHook(() => <%= hookName %>());

    act(() => {
      result.current[2](true);
    });

    expect(result.current[0]).toEqual({ isLoading: false, hasError: true });
  });
});
