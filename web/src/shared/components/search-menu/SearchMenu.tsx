import { Clear as ClearIcon, Search as SearchIcon } from '@mui/icons-material';
import { IconButton, InputAdornment, InputBase, Toolbar } from '@mui/material';
import React, { useCallback, useEffect, useState } from 'react';

import { navigate } from 'gatsby';
import { useLocation } from '@gatsbyjs/reach-router';
import { useSearch } from '../../contexts/search-context/SearchContext';
import { useTheme } from '../../contexts/theme-context/ThemeContext';
import { useUiState } from '../../contexts/ui-state-context/UiStateContext';

interface Props {
  position: 'left' | 'right';
}
/**
 * The Search Menu
 * Provides a Search Icon, input, and controls for a search feature
 * Drop this anywhere in the app
 *
 * @param {*} { position }
 * @returns
 */
const SearchMenu: React.FC<Props> = ({ position }) => {
  const { isDarkMode } = useTheme();
  const { uiState, setUiState } = useUiState();
  const [{ searchQuery }, setSearchQuery] = useSearch();
  const { search } = useLocation();
  const [searchValue, setSearchValue] = useState('');

  useEffect(() => {
    if (searchValue?.length >= 3) {
      const delayDebounce = setTimeout(() => {
        setSearchQuery(searchValue);
      }, 500);

      return () => clearTimeout(delayDebounce);
    }
  }, [searchValue, searchValue?.length, setSearchQuery]);

  useEffect(() => {
    setSearchValue(searchQuery);
  }, [searchQuery]);

  /**
   * This hook is responsible for picking up a incoming ?s=some-value search param and plug it in to the search context
   * it should only run once on a fresh page load
   */
  useEffect(() => {
    if (typeof window === undefined || uiState.isSearchOpen) return;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    const params = new URLSearchParams(search);
    const searchParam = params.get('s');
    if (searchParam && searchParam !== searchQuery) {
      setSearchQuery(searchParam);
      setUiState({ ...uiState, isSearchOpen: true });
    }
    // only run once per page load
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search]);

  const handleExpand = useCallback(() => {
    setUiState({ ...uiState, isSearchOpen: true });
  }, [setUiState, uiState]);

  const handleClear = useCallback(() => {
    const params = new URLSearchParams(search);
    params.delete('s');
    setUiState({ ...uiState, isSearchOpen: false });
    const newUrl = `${window.location.pathname}`;
    void navigate(newUrl);
    setSearchQuery('');
  }, [search, setSearchQuery, setUiState, uiState]);

  const iconButton = (
    <IconButton color="inherit" onClick={handleExpand}>
      <SearchIcon />
    </IconButton>
  );

  const clearButton = (
    <IconButton color="inherit" onClick={handleClear}>
      <ClearIcon />
    </IconButton>
  );

  // const withGraphSwitch = (
  //   <>
  //     <Switch checked={searchState.withGraph} onChange={() => setSearchState({ ...searchState, withGraph: !searchState.withGraph })} />
  //     <label>Graph</label>
  //   </>
  // );

  const searchInput = (
    <InputBase
      sx={{ color: isDarkMode ? 'inherit' : 'white' }}
      placeholder="Search"
      value={searchValue}
      autoFocus
      fullWidth
      onChange={(e) => setSearchValue(e.target.value)}
      endAdornment={
        <>
          <InputAdornment position="end">{clearButton}</InputAdornment>
          {/* {withGraphSwitch} */}
        </>
      }
    />
  );

  return (
    <>
      {position === 'left' ? (
        <>{uiState.isSearchOpen ? searchInput : iconButton}</>
      ) : (
        <>
          {uiState.isSearchOpen ? (
            searchInput
          ) : (
            <Toolbar>
              <>
                <div style={{ flex: 1 }} />
                {iconButton}
              </>
            </Toolbar>
          )}
        </>
      )}
    </>
  );
};

export default SearchMenu;
