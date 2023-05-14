/* eslint-disable jest/expect-expect */
import AppBarComponent from './AppBarComponent';
import React from 'react';
import { render } from '@testing-library/react';

jest.mock('../../shared/components/logo/logo-desktop/LogoDesktop');
jest.mock('../../shared/components/logo/logo-mobile/LogoMobile');
jest.mock('./MainMenuDesktop');
jest.mock('./MainMenuMobile');
jest.mock('./MainToolbar');

describe('AppBarComponent', () => {
  test('renders the component without errors', () => {
    render(<AppBarComponent />);
  });
});
