/* eslint-disable @typescript-eslint/unbound-method */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-return */
import { act, renderHook } from '@testing-library/react-hooks';

import useLocalStorage from './useLocalStorage';

//TODO: fix test with updated testing-library/react-hooks see: https://github.com/testing-library/react-hooks-testing-library/issues/654
describe('useLocalStorage', () => {
  beforeEach(() => {
    Object.defineProperty(window, 'localStorage', {
      value: {
        getItem: jest.fn(),
        setItem: jest.fn(),
        removeItem: jest.fn(),
      },
      writable: true,
    });
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('should return initial value when localStorage is empty', () => {
    const key = 'test-key';
    const initialValue = 'test-value';

    (window.localStorage.getItem as jest.Mock).mockReturnValue(null);

    const { result } = renderHook(() => useLocalStorage(key, initialValue));

    expect(result.current[0]).toEqual(initialValue);
  });

  it('should return value from localStorage when it exists', () => {
    const key = 'test-key';
    const initialValue = 'test-value';
    const storedValue = 'stored-value';

    (window.localStorage.getItem as jest.Mock).mockReturnValue(JSON.stringify(storedValue));

    const { result } = renderHook(() => useLocalStorage(key, initialValue));

    expect(result.current[0]).toEqual(storedValue);
  });

  it('should store value to localStorage', () => {
    const key = 'test-key';
    const initialValue = 'test-value';
    const newValue = 'new-value';

    (window.localStorage.getItem as jest.Mock).mockReturnValue(null);

    const { result } = renderHook(() => useLocalStorage(key, initialValue));

    act(() => {
      result.current[1](newValue);
    });

    expect(result.current[0]).toEqual(newValue);
    expect(window.localStorage.setItem).toHaveBeenCalledWith(key, JSON.stringify(newValue));
  });
});
