import { useState } from 'react';

export type ApiStatus = {
  isLoading: boolean;
  hasError: boolean;
};
/**
 * A custom hook that provides state variables and functions for managing
 * the different states when calling an api
 *
 * @template T - The type of the data returned by the API request.
 * @param {() => Promise<T>} apiCall - A function that returns a Promise for
 * the API call. The function should not take any arguments.
 * @returns {[
 *  T | undefined,
 *  () => Promise<void>,
 *  { isLoading: boolean, setIsLoading: (isLoading: boolean) => void, hasError: boolean, setHasError: (hasError: boolean) => void }
 * ]} - A tuple containing the following items:
 *  - An object containing the following properties and methods:
 *    - `isLoading`: A boolean that indicates whether the API call is currently
 *    in progress.
 *    - `hasError`: A boolean that indicates whether an error occurred during
 *  - `setIsLoading`: A function that sets the `isLoading` value.
 *    the API call.
 *  - `setHasError`: A function that sets the `hasError` value.
 */

const useApiStatus = (): [ApiStatus, (isLoading: boolean) => void, (hasError: boolean) => void] => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [hasError, setHasError] = useState<boolean>(false);
  const apiStatus: ApiStatus = { isLoading, hasError };

  return [apiStatus, setIsLoading, setHasError];
};

export default useApiStatus;
