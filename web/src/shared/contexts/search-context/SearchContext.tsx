import React, { createContext, useContext, useEffect } from 'react';

import { ApiStatus } from '../../hooks/use-api-status/useApiStatus';
import { Dsse } from '../../../generated/graphql';
/* eslint-disable @typescript-eslint/no-empty-function */
import { navigate } from 'gatsby';
import useArchivista from '../../hooks/use-archivista/useArchivista';
import { useLocation } from '@gatsbyjs/reach-router';
import { useUiState } from '../ui-state-context/UiStateContext';

interface SearchContextProps {
  apiStatus: ApiStatus;
  searchQuery: string;
  searchResults: Dsse[];
}

export const SearchContext = createContext<[SearchContextProps, React.Dispatch<React.SetStateAction<string>>]>([
  {
    searchQuery: '',
    searchResults: [],
    apiStatus: {
      isLoading: false,
      hasError: false,
    },
  },
  () => {},
]);

/**
 * The Search Context
 * Provides Search Functionality globally
 *
 * @param {*} { children }
 * @returns
 */
export const SearchProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { uiState } = useUiState();
  const { search } = useLocation();
  const [{ apiStatus, searchQuery, searchResults }, setSearchQuery] = useArchivista();

  useEffect(() => {
    if (typeof window === undefined || !uiState.isSearchOpen) return;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    const params = new URLSearchParams(search);
    if (searchQuery && searchQuery !== '') {
      params.set('s', searchQuery);
    } else {
      params.delete('s');
    }
    const newUrl = `${window.location.pathname}?${params?.toString()}`;
    void navigate(newUrl);
  }, [search, searchQuery, uiState.isSearchOpen]);

  return (
    <SearchContext.Provider
      value={[
        {
          apiStatus,
          searchQuery,
          searchResults,
        },
        setSearchQuery,
      ]}
    >
      {children}
    </SearchContext.Provider>
  );
};

/**
 * Uses the Search Context to bring search behavior into your component.
 * @returns the Search Context
 */
export const useSearch = () => {
  return useContext(SearchContext);
};
