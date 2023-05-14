/*
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";

export enum NodeType {
  TPM = "TPM",
  GCP = "GCP",
  AWS = "AWS",
  AZURE = "AZURE",
  GITHUB = "GITHUB",
}

export enum WorkloadType {
  OCI = "OCI",
  BIN = "BIN",
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

export interface RegistrarService {
  backend: string;
  baseURL: string;
  token: string;

  RegisterNode(request: RegisterNodeRequest): Promise<RegisterResponse>;
  GetNodeRegistrations(): Promise<GetNodeRegistrationsResponse>;
  DeactivateNode(request: DecativateNodeRequest): Promise<never>;
}

export class Registrar implements RegistrarService {
  backend: string;
  baseURL: string;
  token: string;

  constructor(registrarAddress: string, token: string) {
    this.backend = registrarAddress;
    this.baseURL = window.location.origin;
    this.token = token;
  }

  GetNodeRegistrations(initReq?: fm.InitReq): Promise<GetNodeRegistrationsResponse> {
    return fm.fetchReq<never, GetNodeRegistrationsResponse>(`${this.backend}/v1/getnoderegistrations?${fm.renderURLSearchParams({}, [])}`, {
      ...initReq,
      method: "GET",
      headers: { Authorization: this.token },
    });
  }

  RegisterNode(req: RegisterNodeRequest, initReq?: fm.InitReq): Promise<RegisterResponse> {
    return fm.fetchReq<RegisterNodeRequest, RegisterResponse>(`${this.backend}/v1/registernode`, {
      ...initReq,
      method: "POST",
      body: JSON.stringify(req, fm.replacer),
      headers: { Authorization: this.token },
    });
  }

  DeactivateNode(req: DecativateNodeRequest, initReq?: fm.InitReq): Promise<never> {
    return fm.fetchReq<DecativateNodeRequest, never>(`${this.backend}/v1/deactivatenode`, {
      ...initReq,
      method: "POST",
      body: JSON.stringify(req, fm.replacer),
      headers: { Authorization: this.token },
    });
  }
}
