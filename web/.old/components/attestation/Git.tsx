import * as React from "react";
import { GitAttestor } from "../../lib/Witness";
import { NestedListItem } from "../NestedListItem";
import { Typography } from "@mui/material";

export function Git(props: { data: GitAttestor }) {
  const attestation = props.data;
  const statusText = attestation.status === undefined ? "No untracked or modified files" : "";
  const listItems: JSX.Element[] = [];
  if (attestation.status !== undefined) {
    for (const key in attestation.status) {
      listItems.push(<NestedListItem key={key} k={key} v={attestation.status[key] ?? {}} />);
    }
  }

  return (
    <>
      <Typography variant="body2">
        <b>Commithash:</b> {attestation.commithash}
      </Typography>
      <Typography variant="body2">
        <b>Status:</b> {statusText}
      </Typography>
      {listItems.length > 0 ? <ul>{listItems}</ul> : null}
    </>
  );
}
