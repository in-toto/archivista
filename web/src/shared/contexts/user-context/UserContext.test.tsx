import '@testing-library/jest-dom';

import { UserProvider, useUser } from './UserContext';
import { render, screen } from '@testing-library/react';

import React from 'react';
import { Session } from '@ory/client';
import userEvent from '@testing-library/user-event';

describe('UserContext', () => {
  it('should allow setting and getting the user profile', async () => {
    const TestComponent: React.FC = () => {
      const { username, setSession } = useUser();

      const handleSetProfile = () => {
        setSession({ identity: { traits: { username: 'John Doe' } } } as unknown as Session);
      };

      return (
        <div>
          <p>Name: {username}</p>
          <button onClick={handleSetProfile}>Set Profile</button>
        </div>
      );
    };

    render(
      <UserProvider>
        <TestComponent />
      </UserProvider>
    );

    const setProfileButton = screen.getByRole('button', { name: /set profile/i });
    await userEvent.click(setProfileButton);

    expect(screen.getByText('Name: John Doe')).toBeInTheDocument();
  });
});
