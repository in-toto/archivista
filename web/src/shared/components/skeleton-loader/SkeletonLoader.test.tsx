import React from 'react';
import SkeletonLoader from './SkeletonLoader';
import { render } from '@testing-library/react';

describe('SkeletonLoader', () => {
  it('should render without error', () => {
    render(<SkeletonLoader />);
  });
});
