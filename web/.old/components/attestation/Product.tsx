import * as React from "react";
import { HashType, Product } from "../../lib/Witness";
import { Typography } from "@mui/material";

export function Product(props: { data: { [key: string]: Product } }) {
  const attestation = props.data;

  return (
    <ul>
      {Object.keys(attestation).map((key) => (
        <li key={key}>
          <Typography variant="body2">
            <b>Product:</b> {key}
          </Typography>
          <ul>
            <li key={attestation[key]?.mime_type}>
              <Typography variant="body2">
                <b>Mime Type:</b> {attestation[key]?.mime_type}
              </Typography>
            </li>
            {Object.keys(attestation[key]?.digest ?? {}).map((digest_key) => (
              <li key={digest_key}>
                <Typography variant="body2">
                  <b>Digest: </b>
                </Typography>
                {digest_key}: {attestation[key]?.digest[digest_key as HashType]}
              </li>
            ))}
          </ul>
        </li>
      ))}
    </ul>
  );
}
