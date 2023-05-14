/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import * as fm from '../../services/registrar-service/fetch.pb';

import { GetNodeRegistrationsResponse, NodeRegistration } from '../use-registrar/useRegistrar';
import useApiStatus, { ApiStatus } from '../use-api-status/useApiStatus';
import { useCallback, useLayoutEffect, useState } from 'react';

import { useUser } from '../../contexts/user-context/UserContext';

/**
 * A custom hook that fetches node registrations from a registrar address and provides the registrations data, a function to fetch the data, and the API status.
 *
 * @param {string} registrarAddress - The address of the registrar service to fetch node registrations from.
 * @returns {[NodeRegistration[] | undefined, () => Promise<void>, ApiStatus]} An array that contains node registrations data, a function to fetch the data, and the API status.
 *
 * @typedef {object} NodeRegistration - Represents a registered node.
 * @property {string} id - The node ID.
 * @property {string} publicKey - The node public key.
 * @property {string} signingKey - The node signing key.
 * @property {string} address - The node address.
 *
 * @typedef {object} ApiStatus - Represents the status of the API request.
 * @property {boolean} isLoading - Indicates whether the API request is in progress.
 * @property {boolean} hasError - Indicates whether the API request has encountered an error.
 */
const useRegistrations = (registrarAddress: string): [NodeRegistration[] | undefined, () => Promise<void>, ApiStatus] => {
  const userContext = useUser();
  const [token, setToken] = useState<string | undefined>(undefined);
  const [registrations, setRegistrations] = useState<NodeRegistration[] | undefined>(undefined);
  const [apiStatus, setIsLoading, setHasError] = useApiStatus();

  const fetchRegistrations = useCallback(async () => {
    setHasError(false);
    setIsLoading(true);
    if (token) {
      const nodes =
        (
          await fm
            .fetchReq<never, GetNodeRegistrationsResponse>(`${registrarAddress}/v1/getnoderegistrations?${fm.renderURLSearchParams({}, [])}`, {
              method: 'GET',
              headers: { Authorization: token },
            })
            .catch((err) => {
              console.log(err);
              setHasError(true);
            })
        )?.nodeRegistrations || [];
      console.log('nodes: ' + JSON.stringify(nodes, null, 2));
      setRegistrations(nodes);
      setIsLoading(false);
    }
  }, [registrarAddress, setHasError, setIsLoading, token]);

  useLayoutEffect(() => {
    void fetchRegistrations();
  }, [fetchRegistrations]);

  useLayoutEffect(() => {
    if (userContext.username) {
      setToken(userContext.identity?.id);
    }
  }, [userContext.identity?.id, userContext.username]);

  return [registrations, fetchRegistrations, apiStatus];
};

export default useRegistrations;
