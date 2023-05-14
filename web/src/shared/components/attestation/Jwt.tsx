import * as React from 'react';

import { JwtAttestor } from '../../models/Witness';
import { NestedListItem } from '../nested-list-item/NestedListItem';
import Typography from '@mui/material/Typography';

export function Jwt(props: { data: JwtAttestor }) {
  const attestation = props.data;
  return (
    <>
      <Typography variant="body2">
        <b>JWT:</b>
      </Typography>
      <ul>
        <li key={'claims'}>
          <Typography variant="body2">
            <b>Claims:</b>
          </Typography>
          <ul>
            {Object.keys(attestation.claims).map((key) => (
              <>
                <li key={key}>
                  <>
                    <b>{key}:</b>
                    {attestation.claims[key]}
                  </>
                </li>
              </>
            ))}
          </ul>
        </li>
        <li>
          <Typography variant="body2">
            <b>Verified By:</b>
          </Typography>
          <ul>
            <li key={'jwksurl'}>
              <Typography variant="body2">
                <b>JWKS url: </b>
                {attestation.verifiedBy.jwksUrl}
              </Typography>
            </li>
            <NestedListItem key={'jwk'} k="JWK" v={attestation.verifiedBy.jwk} />
          </ul>
        </li>
      </ul>
    </>
  );
}
