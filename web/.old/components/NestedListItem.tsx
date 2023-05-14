import * as React from "react";

export interface Props<TKey extends string | number, TValue extends Record<string, unknown>> {
  k: TKey;
  v: TValue;
}

export function NestedListItem<TKey extends string | number, TValue extends Record<string, unknown>>(props: Props<TKey, TValue>) {
  return (
    <li key={props.k}>
      <b>{props.k}</b>
      <ul>
        {Object.keys(props.v).map((key) => (
          <li key={key}>
            <b>{key}:</b>
            {props.v[key]}
          </li>
        ))}
      </ul>
    </li>
  );
}
