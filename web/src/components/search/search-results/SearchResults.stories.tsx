/* eslint-disable @typescript-eslint/no-empty-function */
import { Meta, Story } from '@storybook/react';
import { SearchContext, SearchResultsState } from '../../../shared/contexts/search-context/SearchContext';

import React from 'react';
import SearchMenu from '../../../shared/components/search-menu/SearchMenu';
import { SearchResults } from './SearchResults';

let searchState: SearchResultsState = { edges: [], nodes: [], searchQuery: '', searchResults: [], withGraph: false };
const getBackRefs = () => {};
const onChange = () => {};
const clearSearch = () => {};
const setSearchState = (s: SearchResultsState) => {
  searchState = s;
};

export default {
  component: SearchResults,
  title: 'Search/Search Results',
  decorators: [
    (Story) => (
      <SearchContext.Provider
        value={[
          {
            searchState,
            getBackRefs,
            onChange,
            clearSearch,
          },
          setSearchState as any,
        ]}
      >
        <Story />
      </SearchContext.Provider>
    ),
  ],
} as Meta;

const Template: Story = () => {
  return (
    <>
      <SearchMenu position="left" />
      <SearchResults />
    </>
  );
};

export const Default = Template.bind({});
