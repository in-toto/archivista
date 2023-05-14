export type Config = {
  hydra: {
    url: string;
    client_id: string;
    client_scopes: string;
    root_url: string;
  };
  registrar_svc: string;
  auth_svc: string;
  archivista_svc: string;
};

export type Repo = {
  id: number;
  repoID: string;
  name: string;
  description: string;
  projecturl: string;
};
