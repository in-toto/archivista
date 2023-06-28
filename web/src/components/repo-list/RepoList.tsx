/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
import { Box, Card, CardContent, TextField, Typography } from '@mui/material';
import React, { useState } from 'react';

import ErrorState from '../../shared/components/error-state/ErrorState';
import RepoCard from '../repo-card/RepoCard';
import SkeletonLoader from '../../shared/components/skeleton-loader/SkeletonLoader';
import Timeline from '../timeline/Timeline';
import useArchivistaByRepos from '../../shared/hooks/use-archivista/useArchivistaByRepo';
import useJudgeApi from '../../shared/hooks/use-judge-api/useJudgeApi';

const RepoList = () => {
  const [searchText, setSearchText] = useState('');
  const [judge, setQuery] = useJudgeApi();
  const [archivista] = useArchivistaByRepos(judge?.results?.map((r) => r.projecturl));

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(event.target.value);
  };

  const filteredRepos = judge.results?.filter((repo) => repo.name.toLowerCase().includes(searchText.toLowerCase()));

  // Sort the filtered repos based on archivista results and timestamps
  // TODO: this should not live here, and when we put all of htis in judge-api, we should handle it there.
  filteredRepos.sort((a, b) => {
    const aArchivistaResults = archivista.results.filter((result) =>
      result?.statement?.subjects?.edges?.some((edge) => edge?.node?.name?.includes(a.projecturl))
    );
    const bArchivistaResults = archivista.results.filter((result) =>
      result?.statement?.subjects?.edges?.some((edge) => edge?.node?.name?.includes(b.projecturl))
    );

    if (aArchivistaResults.length && bArchivistaResults.length) {
      const aLatestTimestamp = Math.max(...aArchivistaResults.map((result) => new Date(result?.signatures?.[0]?.timestamps?.[0]?.timestamp || '').getTime()));
      const bLatestTimestamp = Math.max(...bArchivistaResults.map((result) => new Date(result?.signatures?.[0]?.timestamps?.[0]?.timestamp || '').getTime()));

      return bLatestTimestamp - aLatestTimestamp;
    }

    if (aArchivistaResults.length) {
      return -1;
    }

    if (bArchivistaResults.length) {
      return 1;
    }

    return 0;
  });

  const repos = filteredRepos.map((repo) => {
    // Filter archivista results by the current repo's project URL
    const archivistaResultsByRepo = archivista.results.filter((result) =>
      result?.statement?.subjects?.edges?.some((edge) => edge?.node?.name.endsWith(repo?.projecturl))
    );

    // only show repos with data
    if (archivistaResultsByRepo.length > 1) {
      return <RepoCard key={repo.id} repo={repo} archivistaResultsByRepo={archivistaResultsByRepo} />;
    } else return null;
  });

  const results = (
    <>
      <TextField sx={{ marginBottom: '16px' }} label="Search repos" value={searchText} onChange={handleSearchChange} variant="outlined" fullWidth />
      {repos}
    </>
  );

  return (
    <>
      {judge.apiStatus.isLoading || archivista.apiStatus.isLoading ? (
        <SkeletonLoader />
      ) : judge.apiStatus.hasError || archivista.apiStatus.hasError ? (
        <ErrorState />
      ) : (
        <>{repos.length > 0 && results}</>
      )}
    </>
  );
};

export default RepoList;
