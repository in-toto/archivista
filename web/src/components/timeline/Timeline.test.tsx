import Timeline from './Timeline';
import React from 'react';
import { render } from '@testing-library/react';

describe('Timeline', () => {
  it('should render without error', () => {
    render(<Timeline />);
  });
});
