import React from 'react';
import Timeline from './Timeline';
import { render } from '@testing-library/react';

describe('Timeline', () => {
  // eslint-disable-next-line jest/expect-expect
  it('should render without error', () => {
    render(<Timeline dsseArray={[]} />);
  });
});
