/* eslint-disable @typescript-eslint/no-unsafe-member-access */
import { Log, User, UserManager } from 'oidc-client';

import { Config } from '../../models/app-data-model';

export interface IAuthService {
  getUser(): Promise<User | null>;
  login(): Promise<void>;
  logout(): Promise<void>;
  renewToken(): Promise<User>;
  signinCallback(): Promise<User>;
  storeUser(user: User): Promise<void>;
  getGroups(): Promise<string[]>;
  getIdToken(): Promise<string>;
}

// TODO: refactor to be functional, deprecate classes
export default class AuthService implements IAuthService {
  static getGroups() {
    throw new Error('Method not implemented.');
  }
  public static apiRoot = '';
  _userManager: UserManager;

  constructor(config: Config) {
    const settings = {
      authority: config.hydra.url,
      client_id: config.hydra.client_id,
      redirect_uri: `${config.hydra.root_url}/signin-callback`,
      silent_redirect_uri: `${config.hydra.root_url}/silent-renew`,
      // tslint:disable-next-line:object-literal-sort-keys
      post_logout_redirect_uri: `${config.hydra.root_url}/`,
      response_type: 'code',
      scope: config.hydra.client_scopes,
    };
    this._userManager = new UserManager(settings);

    Log.logger = console;
    Log.level = Log.INFO;
  }

  async getUser(): Promise<User | null> {
    return this._userManager.getUser();
  }

  async login(): Promise<void> {
    return this._userManager.signinRedirect();
  }

  public renewToken(): Promise<User> {
    return this._userManager.signinSilent();
  }

  public logout(): Promise<void> {
    return this._userManager.signoutRedirect();
  }

  async signinCallback(): Promise<User> {
    return this._userManager.signinRedirectCallback();
  }

  async storeUser(user: User): Promise<void> {
    return this._userManager.storeUser(user);
  }

  async getGroups(): Promise<string[]> {
    const groups = this._userManager.getUser().then((user) => {
      if (user) {
        const groups = user.profile.groups as string;
        return groups.split(' ');
      }
      return [];
    });
    return groups;
  }

  async getIdToken(): Promise<string> {
    const id_token = this._userManager.getUser().then((user) => {
      if (user) {
        return user.id_token;
      }
      return '';
    });
    return id_token;
  }
}
