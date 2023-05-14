/* eslint-disable @typescript-eslint/unbound-method */
import AuthService from './AuthService';
import { Config } from '../../models/app-data-model';

describe('AuthService', () => {
  const config = {
    hydra: {
      url: 'https://example.com',
      client_id: 'clientId',
      root_url: 'https://example.com/app',
      client_scopes: 'openid profile email',
    },
  };
  const authService = new AuthService(config as Config);

  describe('login', () => {
    it('should call signinRedirect on the userManager', async () => {
      // TODO one reason we should deprecate classes is they aren't compositional, and we have to access private methods instead of more intuitively mocking imported modules
      authService._userManager.signinRedirect = jest.fn();
      await authService.login();
      expect(authService._userManager.signinRedirect).toHaveBeenCalled();
    });
  });

  describe('logout', () => {
    it('should call signoutRedirect on the userManager', async () => {
      authService._userManager.signoutRedirect = jest.fn();
      await authService.logout();
      expect(authService._userManager.signoutRedirect).toHaveBeenCalled();
    });
  });

  describe('renewToken', () => {
    it('should call signinSilent on the userManager', async () => {
      authService._userManager.signinSilent = jest.fn();
      await authService.renewToken();
      expect(authService._userManager.signinSilent).toHaveBeenCalled();
    });
  });

  describe('getUser', () => {
    it('should call getUser on the userManager', async () => {
      authService._userManager.getUser = jest.fn().mockResolvedValueOnce(null);
      const user = await authService.getUser();
      expect(authService._userManager.getUser).toHaveBeenCalled();
      expect(user).toBeNull();
    });
  });

  describe('getGroups', () => {
    it('should return an empty array if user is null', async () => {
      authService._userManager.getUser = jest.fn().mockResolvedValueOnce(null);
      const groups = await authService.getGroups();
      expect(groups).toEqual([]);
    });

    it('should return an array of groups if user has groups', async () => {
      const user = {
        profile: {
          groups: 'group1 group2',
        },
      };
      authService._userManager.getUser = jest.fn().mockResolvedValueOnce(user);
      const groups = await authService.getGroups();
      expect(groups).toEqual(['group1', 'group2']);
    });
  });

  describe('getIdToken', () => {
    it('should return an empty string if user is null', async () => {
      authService._userManager.getUser = jest.fn().mockResolvedValueOnce(null);
      const idToken = await authService.getIdToken();
      expect(idToken).toBe('');
    });

    it('should return the user id_token if user is not null', async () => {
      const user = {
        id_token: 'someToken',
      };
      authService._userManager.getUser = jest.fn().mockResolvedValueOnce(user);
      const idToken = await authService.getIdToken();
      expect(idToken).toBe('someToken');
    });
  });
});
