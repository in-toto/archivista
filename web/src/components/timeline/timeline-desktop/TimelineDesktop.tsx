/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
import { Avatar, Box, IconButton, Tooltip, Typography, styled } from '@mui/material';

import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CommitLink from '../../commit-link/CommitLink';
import { Dsse } from '../../../generated/graphql';
import LinkIcon from '@mui/icons-material/Link';
import React from 'react';
import { useSearch } from '../../../shared/contexts/search-context/SearchContext';
import { useTheme } from '@mui/material/styles';
import { useUiState } from '../../../shared/contexts/ui-state-context/UiStateContext';

interface DesktopTimelineProps {
  dsseArray: Dsse[];
}

const StyledBox = styled(Box)(({ theme }) => ({
  border: `1px solid ${theme.palette.divider}`,
  padding: '8px',
  borderRadius: '4px',
  cursor: 'pointer',
  position: 'relative',
  '&:hover': {
    borderColor: theme.palette.primary.main,
  },
}));

const TimelineDesktop: React.FC<DesktopTimelineProps> = ({ dsseArray }: DesktopTimelineProps) => {
  const theme = useTheme();

  const [{ searchQuery }, setSearchQuery] = useSearch();
  const handleSearch = (searchValue: string) => {
    setSearchQuery(searchValue);
    setUiState({ ...uiState, isSearchOpen: true });
  };

  const { uiState, setUiState } = useUiState();

  const groupByCommit = dsseArray.reduce((acc: { [key: string]: any[] }, dsse) => {
    const commit = dsse?.statement?.subjects?.edges?.find((edge) => edge?.node?.name?.startsWith('https://witness.dev/attestations/git/v0.1/commithash:'))?.node
      ?.name;

    if (commit) {
      if (!acc[commit]) {
        acc[commit] = [];
      }
      acc[commit].push(dsse);
      acc[commit].sort((a, b) => {
        const aTimestamp = new Date(a.signatures[0].timestamps[0].timestamp).getTime();
        const bTimestamp = new Date(b.signatures[0].timestamps[0].timestamp).getTime();
        return aTimestamp - bTimestamp;
      });
    }

    return acc;
  }, {});

  const sortedCommits = Object.keys(groupByCommit).sort((a, b) => {
    const aTimestamp = new Date(groupByCommit[a][0].signatures[0].timestamps[0].timestamp).getTime();
    const bTimestamp = new Date(groupByCommit[b][0].signatures[0].timestamps[0].timestamp).getTime();
    return bTimestamp - aTimestamp;
  });

  return (
    <Box sx={{ overflow: 'auto', display: 'flex', justifyContent: 'flex-start', alignItems: 'center' }}>
      {sortedCommits.map((commit, index) => {
        const firstStepTimestamp = new Date(groupByCommit[commit][0].signatures[0].timestamps[0].timestamp).getTime();
        const lastStepTimestamp = new Date(groupByCommit[commit][groupByCommit[commit].length - 1].signatures[0].timestamps[0].timestamp).getTime();
        const timeElapsed = lastStepTimestamp - firstStepTimestamp;
        // const previousCommitTimestamp = index === 0 ? null : new Date(groupByCommit[sortedCommits[index - 1]][0].signatures[0].timestamps[0].timestamp);

        return (
          <StyledBox key={commit} sx={{ minWidth: 250, marginRight: '16px' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
              <Avatar sx={{ backgroundColor: 'success.main', marginRight: '8px' }}>
                <CheckCircleIcon />
              </Avatar>
              <CommitLink commit={commit} />
            </Box>
            {groupByCommit[commit].map((dsse, index) => {
              const currentTimestamp = new Date(dsse.signatures[0].timestamps[0].timestamp);
              const currentDate = currentTimestamp.toLocaleDateString('en-US', { timeZone: 'GMT' });
              const currentTime = currentTimestamp.toLocaleTimeString('en-US', { timeZone: 'GMT', hour: '2-digit', minute: '2-digit' });

              const previousTimestamp = index > 0 ? new Date(groupByCommit[commit][index - 1].signatures[0].timestamps[0].timestamp) : null;
              const previousDate = previousTimestamp ? previousTimestamp.toLocaleDateString('en-US', { timeZone: 'GMT' }) : null;

              return (
                <Box key={dsse.id}>
                  {(!previousDate || currentDate !== previousDate) && (
                    <Typography variant="subtitle2" sx={{ fontWeight: 'bold', marginBottom: '4px' }}>
                      {currentDate}
                    </Typography>
                  )}
                  <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
                    <Typography variant="body2">Step: {dsse.statement.attestationCollections.name}</Typography>
                    <Typography variant="caption" sx={{ marginLeft: '8px' }}>
                      {currentTime}
                    </Typography>
                  </Box>
                </Box>
              );
            })}
            <Typography variant="caption" sx={{ marginTop: '4px' }}>
              Duration: {timeElapsed / 1000}s
            </Typography>

            <Box
              sx={{
                position: 'absolute',
                top: 0,
                bottom: 0,
                right: 0,
                width: '48px',
                display: 'flex',
                flexDirection: 'column',
                backgroundColor: theme.palette.grey[100],
                alignItems: 'center',
                justifyContent: 'center',
                borderLeft: `1px solid ${theme.palette.divider}`,
              }}
            >
              <Tooltip title="Find linked attestations">
                <IconButton onClick={() => handleSearch(commit.split(':').pop() as string)}>
                  <LinkIcon />
                </IconButton>
              </Tooltip>
            </Box>
          </StyledBox>
        );
      })}
    </Box>
  );
};

export default TimelineDesktop;
