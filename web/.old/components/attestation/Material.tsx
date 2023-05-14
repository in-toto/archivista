import * as React from "react";
import { MaterialAttestor } from "../../lib/Witness";
import { NestedListItem } from "../NestedListItem";

export function Material(props: { data: MaterialAttestor }) {
  const attestation = props.data;
  return (
    <ul>
      {Object.keys(attestation).map((key) => (
        <NestedListItem key={key} k={key} v={attestation[key] ?? {}} />
      ))}
    </ul>
  );
}
