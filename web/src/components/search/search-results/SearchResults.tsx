import * as React from 'react';

import { DsseEnvelope } from './DsseEnvelope';
import ErrorState from '../../../shared/components/error-state/ErrorState';
import SearchEmptyState from '../search-empty-state/SearchEmptyState';
import SearchSkeletonLoader from '../../../shared/components/skeleton-loader/SkeletonLoader';
import { useSearch } from '../../../shared/contexts/search-context/SearchContext';

// graph can't be used at build time, this prevents an error
// const { DisplayGraph } = typeof window !== 'undefined' ? require('./Graph') : { DisplayGraph: null };
/**
 * The Search Results
 * Provides the Search Results behavior for the search feature
 * Drop this anywhere in the app
 *
 * @returns
 */
const SearchResults = () => {
  const [search] = useSearch();

  const showErrorState = search.apiStatus.hasError;
  const showResults = search?.searchResults?.length > 0;
  const isLoading = search.apiStatus.isLoading && !showResults;
  const showEmptyState = !showErrorState && !isLoading && !search.searchQuery;
  return (
    <>
      {showErrorState && <ErrorState />}
      {isLoading && <SearchSkeletonLoader />}
      {showEmptyState && <SearchEmptyState />}
      {/* {showGraph && <DisplayGraph nodes={search.searchState.nodes} edges={search.searchState.edges} />} */}
      {showResults && search.searchResults.map((result) => <DsseEnvelope key={result.gitoidSha256} result={result} />)}
    </>
  );
};

export default SearchResults;
