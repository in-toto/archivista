import { Dsse } from '../../generated/graphql';
import React from 'react';
import TimelineDesktop from './timeline-desktop/TimelineDesktop';
import TimelineMobile from './timeline-mobile/TimelineMobile';
import { useMediaQuery } from '@mui/material';

interface TimelineProps {
  dsseArray: Dsse[];
}

const Timeline = ({ dsseArray }: TimelineProps) => {
  const isMobile = useMediaQuery('(max-width:600px)');

  return (
    <>
      {isMobile && <TimelineMobile dsseArray={dsseArray} />}
      {!isMobile && <TimelineDesktop dsseArray={dsseArray} />}
    </>
  );
};

export default Timeline;
