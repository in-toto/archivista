export interface Node {
  oid: string;
  name?: string;
  subjectNames?: (string | undefined)[];
  backrefs: string[];
}

export interface Edge {
  subjectNames: string[];
  from_oid: string;
  to_oid: string;
}
