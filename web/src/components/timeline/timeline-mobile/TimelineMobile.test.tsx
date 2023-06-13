import React from 'react';
import TimelineMobile from './TimelineMobile';
import { render } from '@testing-library/react';

describe('TimelineMobile', () => {
  // eslint-disable-next-line jest/expect-expect
  it('should render without error', () => {
    render(<TimelineMobile dsseArray={[]} />);
  });
});
