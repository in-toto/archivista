import React, { PropsWithChildren, createContext, useContext, useState } from 'react';
import useRegistrar, { IRegistrarService } from '../../hooks/use-registrar/useRegistrar';

import { ConfigConstants } from '../../constants';

interface RegistrarContextProps {
  registrar?: IRegistrarService;
  setToken?: React.Dispatch<string | undefined>;
}

const RegistrarContext = createContext<RegistrarContextProps>({
  registrar: undefined,
});
/**
 * The Auth Context
 * Provides the Auth Service globally
 *
 * @param {*} { children }
 * @returns
 */
export const RegistrarProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const [registrar] = useState<IRegistrarService | undefined>(useRegistrar(ConfigConstants.registrar_svc));

  return <RegistrarContext.Provider value={{ registrar }}>{children}</RegistrarContext.Provider>;
};

/**
 * A shortcut for importing the UserContext global state
 * @returns the UserContext global state
 */
export const useRegistrarContext = () => useContext(RegistrarContext);
