import React, { PropsWithChildren } from 'react';

import NoAuth from './components/auth/no-auth/NoAuth';
import { useUser } from './shared/contexts/user-context/UserContext';

/**
 * This logic makes sure that we don't render any of the children passed in unless the user is signed-in
 */

const Auth: React.FC<PropsWithChildren> = ({ children }) => {
  const { username } = useUser();

  return <> {username ? children : <NoAuth />} </>;
};

export default Auth;
