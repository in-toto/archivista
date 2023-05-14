/* eslint-disable @typescript-eslint/no-empty-function */
import * as fm from '../../services/registrar-service/fetch.pb';

import { useState } from 'react';
import { useUser } from '../../contexts/user-context/UserContext';

export enum NodeType {
  TPM = 'TPM',
  GCP = 'GCP',
  AWS = 'AWS',
  AZURE = 'AZURE',
  GITHUB = 'GITHUB',
}

export enum WorkloadType {
  OCI = 'OCI',
  BIN = 'BIN',
}

export type RegisterNodeRequest = {
  nodeType?: NodeType;
  selectors?: SelectorList;
  nodeGroup?: string;
};

export type DecativateNodeRequest = {
  nodeId?: string;
};

export type GetNodeRegistrationsResponse = {
  nodeRegistrations?: NodeRegistration[];
};

export type NodeRegistration = {
  nodeId?: string;
  nodeType?: NodeType;
  selectors?: SelectorList;
  nodeGroup?: string;
  registeredBy?: string;
  accountUUID?: string;
  registeredAt?: string;
};

export type RegisterWorkloadRequest = {
  workloadType?: WorkloadType;
  selectors?: SelectorList;
  nodeGroup?: string;
  workloadName?: string;
};

export type RegisterResponse = {
  spiffeID?: string;
};

export type GetAccountRequest = {
  selectors?: SelectorList;
  nodeType?: NodeType;
};

export type GetAccountResponse = {
  account?: string;
  nodeGroup?: string;
};

export type Selector = {
  key?: string;
  value?: string;
};

export type SelectorList = {
  selectors?: Selector[];
};

export interface IRegistrarService {
  backend: string;
  baseURL: string;
  token?: string;
  setToken?: React.Dispatch<string | undefined>;

  registerNode(request: RegisterNodeRequest): Promise<RegisterResponse | void>;
  getNodeRegistrations(): Promise<GetNodeRegistrationsResponse | void>;
  deactivateNode(request: DecativateNodeRequest): Promise<never>;
}
/**
 * A custom React hook for interacting with the Registrar service.
 *
 * @param {string} registrarAddress - The URL for the Registrar service backend.
 * @param {string} token - An optional authentication token for making authorized requests.
 *
 * @returns {Object} - An object containing functions for interacting with the Registrar service:
 *   - registerNode: A function for registering a node with the service.
 *   - getNodeRegistrations: A function for getting a list of node registrations from the service.
 *   - deactivateNode: A function for deactivating a node with the service.
 * @param {string} registrarAddress
 * @param {string} [token]
 * @returns {IRegistrarService}
 */
// TODO: refactor state!!
const useRegistrar = (registrarAddress: string): IRegistrarService => {
  const [backend] = useState(registrarAddress);
  const [baseURL] = useState(window.location.origin);
  const user = useUser();
  const [serviceToken, setToken] = useState(user?.identity?.id);

  function registerNode(request: RegisterNodeRequest): Promise<RegisterResponse | void> {
    if (serviceToken) {
      return fm.fetchReq<RegisterNodeRequest, RegisterResponse>(`${backend}/v1/registernode`, {
        method: 'POST',
        body: JSON.stringify(request, fm.replacer),
        headers: { Authorization: serviceToken },
      });
    }

    return new Promise(() => {});
  }

  function getNodeRegistrations(): Promise<GetNodeRegistrationsResponse | void> {
    if (serviceToken) {
      return fm.fetchReq<never, GetNodeRegistrationsResponse>(`${backend}/v1/getnoderegistrations?${fm.renderURLSearchParams({}, [])}`, {
        method: 'GET',
        headers: { Authorization: serviceToken },
      });
    }

    return new Promise(() => {});
  }

  function deactivateNode(request: DecativateNodeRequest): Promise<never> {
    if (serviceToken) {
      return fm.fetchReq<DecativateNodeRequest, never>(`${backend}/v1/deactivatenode`, {
        method: 'POST',
        body: JSON.stringify(request, fm.replacer),
        headers: { Authorization: serviceToken },
      });
    }

    return new Promise(() => {});
  }

  return {
    backend,
    baseURL,
    token: serviceToken,
    setToken,
    registerNode,
    getNodeRegistrations,
    deactivateNode,
  };
};

export default useRegistrar;
