/* eslint-disable @typescript-eslint/no-empty-function */
import { Configuration, FrontendApi, Identity, Session } from '@ory/client';
import React, { PropsWithChildren, createContext, useContext, useEffect, useState } from 'react';

interface UserContextProps {
  identity?: Identity;
  logoutUrl?: string;
  session?: Session;
  setSession: (session?: Session) => void;
  username?: string;
}

const UserContext = createContext<UserContextProps>({
  setSession: () => {},
});

// Get your Ory url from .env
// Or localhost for local development
const basePath = process.env.REACT_APP_ORY_URL || '/kratos';
const ory = new FrontendApi(
  new Configuration({
    basePath,
    baseOptions: {
      withCredentials: true,
    },
  })
);

// eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
const getUserName = (identity: Identity) => (identity?.traits?.email || identity?.traits?.username) as string;

/**
 * The User Context
 * Provides the user globally
 * Should be undefined if the user is not signed in
 *
 * @param {*} { children }
 * @returns
 */
export const UserProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const [session, setSession] = useState<Session | undefined>(undefined);
  const [username, setUsername] = useState<string>();
  const [identity, setIdentity] = useState<Identity>();
  const [logoutUrl, setLogoutUrl] = useState<string | undefined>();

  /**
   * This hook makes sure to set the token for the registrar if our user id changes
   */
  // Second, gather session data, if the user is not logged in, redirect to login
  useEffect(() => {
    const getSession = async () => {
      const { data } = await ory.toSession();
      if (data) {
        // User has a session!
        setSession(data);
        const { data: logoutFlow } = await ory.createBrowserLogoutFlow();
        // Get also the logout url
        setLogoutUrl(logoutFlow?.logout_url);
      }
    };

    getSession().catch((err) => {
      // NOTE: This is where we handle login errors gracefully. Ory docs suggests to go to the login page. See: https://www.ory.sh/docs/kratos/quickstart
      console.error(err);
      // Redirect to login page
      window.location.replace('/login/login');
    });
  }, [setLogoutUrl, setSession]);

  useEffect(() => {
    if (session) {
      console.log(session);
      setIdentity(session?.identity);
    }
  }, [session]);

  useEffect(() => {
    if (identity) {
      console.log(identity);
      setUsername(getUserName(identity));
    }
  }, [identity]);

  return <UserContext.Provider value={{ session, setSession, username, identity, logoutUrl }}>{children}</UserContext.Provider>;
};

/**
 * A shortcut for importing the UserContext global state
 * @returns the UserContext global state
 */
export const useUser = () => useContext(UserContext);
