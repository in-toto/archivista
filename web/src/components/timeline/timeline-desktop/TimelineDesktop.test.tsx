import React from 'react';
import TimelineDesktop from './TimelineDesktop';
import { render } from '@testing-library/react';

describe('TimelineDesktop', () => {
  // eslint-disable-next-line jest/expect-expect
  it('should render without error', () => {
    render(<TimelineDesktop dsseArray={[]} />);
  });
});
