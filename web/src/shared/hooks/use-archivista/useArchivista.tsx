import { Dsse, useBySubjectDigestLazyQuery } from '../../../generated/graphql';
import useApiStatus, { ApiStatus } from '../use-api-status/useApiStatus';
import { useCallback, useEffect, useState } from 'react';

export type ArchivistaProps = {
  apiStatus: ApiStatus;
  searchQuery: string;
  searchResults: Dsse[];
};

/**
 * A hook for calling Archivista via graphql
 *
 */
const useArchivista = (): [ArchivistaProps, React.Dispatch<React.SetStateAction<string>>] => {
  const [searchQuery, setSearchQuery] = useState('');
  const [apiStatus, setIsLoading, setHasError] = useApiStatus();
  const [getEnvelopes] = useBySubjectDigestLazyQuery();
  const [searchResults, setSearchResults] = useState([] as Dsse[]);
  const archivistaProps: ArchivistaProps = { apiStatus, searchResults, searchQuery };

  async function sha256(s: string) {
    const utf8 = new TextEncoder().encode(s);
    const hashBuffer = await crypto.subtle.digest('SHA-256', utf8);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map((bytes) => bytes.toString(16).padStart(2, '0')).join('');
    return hashHex;
  }

  const isHash = (text: string) => {
    // Regular expression to check if string is a SHA256 hash
    const regexExpSHA256 = /^[a-f0-9]{64}$/gi;
    const regexExpSHA1 = /^[a-f0-9]{40}$/gi;

    return regexExpSHA256.test(text) || regexExpSHA1.test(text);
  };

  const backRefs = [
    'https://witness.dev/attestations/git/v0.1/commithash',
    'https://witness.dev/attestations/gitlab/v0.1/pipelineurl',
    'https://witness.dev/attestations/github/v0.1/pipelineurl',
  ];

  const executeSearch = useCallback(() => {
    let digest = '';
    setHasError(false);
    setIsLoading(true);
    if (isHash(searchQuery)) {
      getEnvelopes({ variables: { digest: searchQuery } })
        .then((res) => {
          const freshResults = res.data?.dsses?.edges?.map((edge) => edge?.node) || [];
          setSearchResults(freshResults as Dsse[]);
          console.log(JSON.stringify(freshResults, null, 2));
        })
        .finally(() => {
          setIsLoading(false);
        })
        .catch((e) => {
          setHasError(true);
          setIsLoading(false);
        });
    } else {
      sha256(searchQuery)
        .then((hash) => {
          digest = hash;
        })
        .then(async () => {
          await getEnvelopes({ variables: { digest: digest } }).then((res) => {
            const freshResults = res.data?.dsses?.edges?.map((edge) => edge?.node) || [];
            setSearchResults(freshResults as Dsse[]);
            console.log(JSON.stringify(freshResults, null, 2));
          });
        })
        .finally(() => {
          setIsLoading(false);
        })
        .catch((e) => {
          setHasError(true);
          setIsLoading(false);
        });
    }
  }, [getEnvelopes, searchQuery, setHasError, setIsLoading]);

  useEffect(() => {
    if (searchQuery === '') {
      // unset search results when query is set back to default
      setSearchResults([]);
    } else if (searchQuery?.length >= 3) {
      // console.log('exec search!');
      executeSearch();
    }
  }, [executeSearch, searchQuery]);

  return [archivistaProps, setSearchQuery];
};

export default useArchivista;
