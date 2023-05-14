import '@testing-library/jest-dom';

import Auth from './Auth';
import React from 'react';
import { render } from '@testing-library/react';
import { useUser } from './shared/contexts/user-context/UserContext';

jest.mock('./shared/contexts/user-context/UserContext');

// Mock child components
// eslint-disable-next-line react/display-name
jest.mock('./components/auth/no-auth/NoAuth', () => () => <div>NoAuth</div>);

describe('Auth component', () => {
  it('should render children when user is signed in', () => {
    // Mock the user context hook to return a signed in user
    (useUser as jest.Mock).mockReturnValue({ username: 'testUser' });

    const { getByText, queryByText } = render(
      <Auth>
        <div>Test Child</div>
      </Auth>
    );

    expect(getByText('Test Child')).toBeInTheDocument();
    expect(queryByText('NoAuth')).toBeNull();
  });

  it('should render NoAuth component when user is not signed in', () => {
    // Mock the user context hook to return a signed out user
    (useUser as jest.Mock).mockReturnValue({ username: undefined });

    const { getByText, queryByText } = render(
      <Auth>
        <div>Test Child</div>
      </Auth>
    );

    expect(getByText('NoAuth')).toBeInTheDocument();
    expect(queryByText('Test Child')).toBeNull();
  });
});
