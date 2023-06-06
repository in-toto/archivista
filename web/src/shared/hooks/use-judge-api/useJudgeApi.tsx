import useApiStatus, { ApiStatus } from '../use-api-status/useApiStatus';
import { useCallback, useLayoutEffect, useState } from 'react';

import { Repo } from '../../models/app-data-model';

export type JudgeApiProps = {
  apiStatus: ApiStatus;
  results: Repo[]; // TODO: set type later
  query: string;
};

/**
 * A hook for using the Judge graphql api
 *
 */
const useJudgeApi = (): [JudgeApiProps, (s: string) => void] => {
  const [apiStatus, setIsLoading, setHasError] = useApiStatus();
  const [results, setApiResults] = useState([] as []);

  //TODO update query and use apollo like search
  //TODO make the search query seperate from the graphql query
  const [query, setQuery] = useState(`
  query {
    projects {
      id
      repoID
      name
      projecturl
    }
  }
  `);
  const judgeApiProps: JudgeApiProps = { apiStatus, results, query };

  const getQuery = useCallback(async () => {
    try {
      setIsLoading(true);
      setHasError(false);
      const response = await fetch(`judge-api/query?query=${query}`, {
        method: 'GET',
        credentials: 'include',
      });

      const json = await response.json();
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
      setApiResults(json?.data?.projects);
      console.log(JSON.stringify(json, null, 2));
    } catch (error) {
      console.error(error);
      setHasError(true);
    } finally {
      setIsLoading(false);
    }
  }, [query, setHasError, setIsLoading]);

  useLayoutEffect(() => {
    if (query?.length >= 3) {
      void getQuery();
    }
  }, [getQuery, query?.length]);

  return [judgeApiProps, setQuery];
};

export default useJudgeApi;
