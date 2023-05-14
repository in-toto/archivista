import { ConfigConstants } from '../shared/constants';
import React from 'react';
import { UserManager } from 'oidc-client';
import { useEffect } from 'react';

// TODO: error state
const SilentRenew = (): JSX.Element => {
  useEffect(() => {
    void (() => {
      const settings = {
        authority: ConfigConstants.hydra.url,
        client_id: ConfigConstants.hydra.client_id,
        redirect_uri: `${ConfigConstants.hydra.root_url}/signin-callback`,
        silent_redirect_uri: `${ConfigConstants.hydra.root_url}/silent-renew`,
        post_logout_redirect_uri: `${ConfigConstants.hydra.root_url}/`,
        response_type: 'code',
        scope: ConfigConstants.hydra.client_scopes,
      };
      const userManager = new UserManager(settings);
      userManager.signinSilentCallback().catch(function (error) {
        console.error(error);
      });
    })();
  }, []);

  return <h1>Renewing...</h1>;
};

export default SilentRenew;
