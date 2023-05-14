import TimelineDesktop from './TimelineDesktop';
import React from 'react';
import { render } from '@testing-library/react';

describe('TimelineDesktop', () => {
  it('should render without error', () => {
    render(<TimelineDesktop />);
  });
});
