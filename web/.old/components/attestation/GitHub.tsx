import * as React from "react";

import { GithubAttestor } from "../../lib/Witness";
import Typography from "@mui/material/Typography";

export function Github(props: { data: GithubAttestor }) {
  const attestation = props.data;

  return (
    <>
      <Typography variant="body2">
        <b>Actor:</b> {attestation.claims.actor}
      </Typography>
      <Typography variant="body2">
        <b>Actor ID:</b> {attestation.claims.actor_id}
      </Typography>
      <Typography variant="body2">
        <b>Audience:</b> {attestation.claims.aud}
      </Typography>
      <Typography variant="body2">
        <b>Base ref:</b> {attestation.claims.base_ref}
      </Typography>
      <Typography variant="body2">
        <b>Event name:</b> {attestation.claims.event_name}
      </Typography>
      <Typography variant="body2">
        <b>Expiration:</b> {attestation.claims.exp}
      </Typography>
      <Typography variant="body2">
        <b>Head ref:</b> {attestation.claims.head_ref}
      </Typography>
      <Typography variant="body2">
        <b>Issued at:</b> {attestation.claims.iat}
      </Typography>
      <Typography variant="body2">
        <b>Issuer:</b> {attestation.claims.iss}
      </Typography>
      <Typography variant="body2">
        <b>Job workflow ref:</b> {attestation.claims.job_workflow_ref}
      </Typography>
      <Typography variant="body2">
        <b>JWT ID:</b> {attestation.claims.jti}
      </Typography>
      <Typography variant="body2">
        <b>Not before:</b> {attestation.claims.nbf}
      </Typography>
      <Typography variant="body2">
        <b>Ref:</b> {attestation.claims.ref}
      </Typography>
      <Typography variant="body2">
        <b>Ref type:</b> {attestation.claims.ref_type}
      </Typography>
      <Typography variant="body2">
        <b>Repository:</b> {attestation.claims.repository}
      </Typography>
      <Typography variant="body2">
        <b>Repository ID:</b> {attestation.claims.repository_id}
      </Typography>
      <Typography variant="body2">
        <b>Repository owner:</b> {attestation.claims.repository_owner}
      </Typography>
      <Typography variant="body2">
        <b>Repository owner ID:</b> {attestation.claims.repository_owner_id}
      </Typography>
      <Typography variant="body2">
        <b>Repository visibility:</b> {attestation.claims.repository_visibility}
      </Typography>
      <Typography variant="body2">
        <b>Run attempt:</b> {attestation.claims.run_attempt}
      </Typography>
      <Typography variant="body2">
        <b>Run ID:</b> {attestation.claims.run_id}
      </Typography>
      <Typography variant="body2">
        <b>Run number:</b> {attestation.claims.run_number}
      </Typography>
      <Typography variant="body2">
        <b>SHA:</b> {attestation.claims.sha}
      </Typography>
      <Typography variant="body2">
        <b>Subject:</b> {attestation.claims.sub}
      </Typography>
      <Typography variant="body2">
        <b>Workflow:</b> {attestation.claims.workflow}
      </Typography>
      <Typography variant="body2">
        <b>Verified by:</b> {attestation.verifiedBy?.jwksUrl}
      </Typography>
      <Typography variant="body2">
        <b>CI config path:</b> {attestation.ciconfigpath}
      </Typography>
      <Typography variant="body2">
        <b>Pipeline ID:</b> {attestation.pipelineid}
      </Typography>
      <Typography variant="body2">
        <b>Pipeline name:</b> {attestation.pipelinename}
      </Typography>
      <Typography variant="body2">
        <b>Pipeline url:</b> {attestation.pipelineurl}
      </Typography>

      <Typography variant="body2">
        <b>Project url:</b> {attestation.projecturl}
      </Typography>
      <Typography variant="body2">
        <b>Runner ID:</b> {attestation.runnerid}
      </Typography>
      <Typography variant="body2">
        <b>CI host:</b> {attestation.cihost}
      </Typography>
      <Typography variant="body2">
        <b>CI server url:</b> {attestation.ciserverurl}
      </Typography>
      <Typography variant="body2">
        <b>Runner arch:</b> {attestation.runnerarch}
      </Typography>
      <Typography variant="body2">
        <b>Runner os:</b> {attestation.runneros}
      </Typography>
    </>
  );
}
