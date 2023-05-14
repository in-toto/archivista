import { ApolloClient, ApolloProvider, InMemoryCache } from '@apollo/client';
import React, { PropsWithChildren, useEffect } from 'react';

import { ConfigConstants } from './shared/constants';
import { LocationProvider } from '@gatsbyjs/reach-router';
import { SearchProvider } from './shared/contexts/search-context/SearchContext';
import { ThemeProvider } from './shared/contexts/theme-context/ThemeContext';
import { UiStateProvider } from './shared/contexts/ui-state-context/UiStateContext';
import { UserProvider } from './shared/contexts/user-context/UserContext';

/**
 * This is the HighOrder app component fot bootstrapping JS things into our app globally, above the layout and pages.
 * There is one more level above this in gatsby-browser.ts for static things that we bootstrap globally.
 *
 * @param {*} { children }
 * @returns
 */
const App: React.FC<PropsWithChildren> = ({ children }) => {
  const [graphQLClient, setGraphQLClient] = React.useState<ApolloClient<unknown>>();
  useEffect(() => {
    const client = new ApolloClient({
      cache: new InMemoryCache(),
      uri: `${ConfigConstants?.archivista_svc || ''}/query`,
    });
    setGraphQLClient(client);
  }, []);

  return (
    <>
      <ThemeProvider>
        <LocationProvider>
          <UiStateProvider>
            <UserProvider>
              {graphQLClient ? (
                <ApolloProvider client={graphQLClient}>
                  <SearchProvider>{children}</SearchProvider>
                </ApolloProvider>
              ) : (
                children
              )}
            </UserProvider>
          </UiStateProvider>
        </LocationProvider>
      </ThemeProvider>
    </>
  );
};

export default App;
