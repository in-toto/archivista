/* eslint-disable react/display-name */
import '@testing-library/jest-dom';

import LogoDesktop from './LogoDesktop';
import React from 'react';
import { render } from '@testing-library/react';
jest.mock('gatsby', () => ({
  Link: 'a',
}));

jest.mock('../logo-image/LogoImage', () => () => <div>My Logo</div>);

describe('LogoDesktop', () => {
  it('should render the component', () => {
    const { getByText } = render(<LogoDesktop />);
    expect(getByText('Judge')).toBeInTheDocument();
  });
});
