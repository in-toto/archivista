import Footer from './Footer';
import React from 'react';
import { render } from '@testing-library/react';

describe('Footer', () => {
  it('should render without error', () => {
    render(<Footer />);
  });
});
