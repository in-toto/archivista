import '@testing-library/jest-dom';

import * as React from 'react';

import NotFoundPage from './404';
import { render } from '@testing-library/react';

describe('NotFoundPage component', () => {
  test('renders "Page not found" heading', () => {
    const { getByText } = render(<NotFoundPage />);
    const headingText = getByText(/page not found/i);
    expect(headingText).toBeInTheDocument();
  });

  test('renders "Go home" link', () => {
    const { getByText } = render(<NotFoundPage />);
    const linkText = getByText(/go home/i);
    expect(linkText).toBeInTheDocument();
  });
});
