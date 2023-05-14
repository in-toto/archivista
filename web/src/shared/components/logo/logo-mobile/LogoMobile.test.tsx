/* eslint-disable react/display-name */
/* eslint-disable @typescript-eslint/no-empty-function */
import '@testing-library/jest-dom';

import { Link as GatsbyLink } from 'gatsby';
import LogoImage from '../logo-image/LogoImage';
import LogoMobile from './LogoMobile';
import React from 'react';
import { UiStateContext } from '../../../contexts/ui-state-context/UiStateContext';
import { render } from '@testing-library/react';

jest.mock('../logo-image/LogoImage', () => () => <div>My Logo</div>);

describe('LogoMobile', () => {
  it('renders the mobile logo', () => {
    const { getByText, getByTestId } = render(
      <UiStateContext.Provider value={{ uiState: { isSearchOpen: false }, setUiState: () => {} }}>
        <LogoMobile />
      </UiStateContext.Provider>
    );

    const linkElement = getByTestId('logo');
    const titleElement = getByText(/judge/i);

    expect(linkElement).toBeInTheDocument();
    expect(titleElement).toBeInTheDocument();

    expect(linkElement).toHaveAttribute('href', '/');
    expect(titleElement).toHaveTextContent('Judge');
  });

  it('hides the mobile logo when search is open', () => {
    const { queryByTestId } = render(
      <UiStateContext.Provider value={{ uiState: { isSearchOpen: true }, setUiState: () => {} }}>
        <LogoMobile />
      </UiStateContext.Provider>
    );

    const linkElement = queryByTestId('logo');

    expect(linkElement).toBeFalsy();
  });
});
