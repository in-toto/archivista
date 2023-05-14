/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/react';

import NoAuth from './NoAuth';
import React from 'react';
import useContent from '../../../shared/hooks/use-content/useContent';

jest.mock('../../../shared/hooks/use-content/useContent');

describe('NoAuth', () => {
  beforeEach(() => {
    (useContent as jest.Mock).mockReturnValue({
      contentYaml: {
        signInButtonText: 'Sign In',
      },
      noAuthYaml: {
        noAuthHeading: 'No Auth Heading',
        noAuthSubheading: 'No Auth Subheading',
      },
    });
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should render no auth heading and subheading', () => {
    render(<NoAuth />);
    expect(screen.getByText('No Auth Heading')).toBeInTheDocument();
    expect(screen.getByText('No Auth Subheading')).toBeInTheDocument();
  });
});
