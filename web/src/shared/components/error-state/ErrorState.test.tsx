import MyComponent from './ErrorState';
import React from 'react';
import { render } from '@testing-library/react';

describe('MyComponent', () => {
  it('should render without error', () => {
    render(<MyComponent />);
  });
});
