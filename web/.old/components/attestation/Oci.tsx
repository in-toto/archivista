import * as React from "react";
import { HashType, OciAttestor } from "../../lib/Witness";
import { Typography } from "@mui/material";

export function Oci(props: { data: OciAttestor }) {
  const attestation = props.data;
  return (
    <>
      <Typography variant="body2">
        <b>Tar Digest:</b>
      </Typography>
      <ul>
        {Object.keys(attestation.tardigest).map((key) => (
          <li key={key}>
            <b>{key}:</b>
            {attestation.tardigest[key as HashType]}
          </li>
        ))}
      </ul>

      <Typography variant="body2">
        <b>Manifest:</b>
      </Typography>
      <ul>
        {attestation.manifest.map((key) => (
          <>
            <li>
              <Typography variant="body2">
                <b>Config:</b> {key.Config}
              </Typography>
            </li>
            <li>
              <Typography variant="body2">
                <b>Repo tags:</b>
              </Typography>
              <ul>
                {key.RepoTags.map((repoTagKey) => (
                  <li key={repoTagKey}>
                    <Typography variant="body2">{repoTagKey}</Typography>
                  </li>
                ))}
              </ul>
            </li>
            <li>
              <Typography variant="body2">
                <b>Layers:</b>
              </Typography>
              <ul>
                {key.Layers.map((layersKey) => (
                  <li key={layersKey}>
                    <Typography variant="body2">{layersKey}</Typography>
                  </li>
                ))}
              </ul>
            </li>
          </>
        ))}
      </ul>

      <Typography variant="body2">
        <b>Image tags:</b>
      </Typography>
      <ul>
        {attestation.imagetags.map((tag) => (
          <li key={tag}>
            <b>{tag}</b>
          </li>
        ))}
      </ul>
      <Typography variant="body2">
        <b>Diffids:</b>
      </Typography>
      <ul>
        {attestation.diffids.map((key) =>
          Object.keys(key).map((key) => (
            <li key={key}>
              <b>{key}:</b>
              {attestation.imageid[key as HashType]}
            </li>
          ))
        )}
      </ul>
      <Typography variant="body2">
        <b>Image id:</b>
      </Typography>
      <ul>
        {Object.keys(attestation.imageid).map((key) => (
          <li key={key}>
            <b>{key}:</b>
            {attestation.imageid[key as HashType]}
          </li>
        ))}
      </ul>
      <Typography variant="body2">
        <b>Manifest raw:</b> {attestation.manifestraw}
      </Typography>
    </>
  );
}
