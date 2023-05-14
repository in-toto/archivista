import * as React from "react";

import { Button, FormControl, Switch, TextField } from "@mui/material";
import { Dsse, useBySubjectDigestLazyQuery } from "../generated/graphql";

import { DisplayGraph } from "./Graph";
import { DsseEnvelope } from "../components/search/DsseEnvelope";
import { useState } from "react";

//make search take n hashes.  we just add backrefs to the query
const backRefs = [
  "https://witness.dev/attestations/git/v0.1/commithash",
  "https://witness.dev/attestations/gitlab/v0.1/pipelineurl",
  "https://witness.dev/attestations/github/v0.1/pipelineurl",
];

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

export const Search = () => {
  const [search, setSearch] = useState("");
  const [results, setResults] = useState<Dsse[]>([]);
  const [nodesState, setNodesState] = useState<Node[]>([]);
  const [edgesState, setEdgesState] = useState<Edge[]>([]);
  const [useGraph, setUseGraph] = useState(false);

  const [getEnvelopes] = useBySubjectDigestLazyQuery();

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(event.target.value);
  };

  const populateGraph = (envelopes: Dsse[]) => {
    const nodes: Node[] = [];
    const edges: Edge[] = [];

    envelopes.forEach((envelope) => {
      const backRefsSubjectNames: string[] = [];

      backRefs.forEach((backref) => {
        envelope.statement?.subjects.edges?.forEach((subject) => {
          if (subject?.node?.name.includes(backref)) {
            backRefsSubjectNames.push(subject?.node?.name);
          }
        });
      });

      const node = {
        oid: envelope.gitoidSha256,
        name: envelope.statement?.attestationCollections?.name,
        subjectNames: envelope?.statement?.subjects?.edges?.map((edge) => edge?.node?.name),
        backrefs: backRefsSubjectNames,
      };

      nodes.push(node);
    });

    //loop through all nodes and make an edge where two nodes share a backref
    nodes.forEach((node) => {
      nodes.forEach((node2) => {
        if (node.oid !== node2.oid) {
          const sharedBackRefs = node.backrefs.filter((backref) => node2.backrefs.includes(backref));

          if (sharedBackRefs.length > 0) {
            const edge = {
              subjectNames: sharedBackRefs,
              from_oid: node.oid,
              to_oid: node2.oid,
            };

            edges.push(edge);
          }
        }
      });
    });

    setNodesState(nodes);
    setEdgesState(edges);

    console.log(nodes);
    console.log(edges);
  };

  const combineArrays = (arr1: Dsse[], arr2: Dsse[]) => {
    //dedup on key gitoid_sha256 and combine
    const map = new Map();
    arr1.forEach((item) => {
      map.set(item.gitoidSha256, item);
    });
    arr2.forEach((item) => {
      map.set(item.gitoidSha256, item);
    });

    return Array.from(map.values()) as Dsse[];
  };

  const handleClick = () => {
    let digest = "";

    if (search == "") {
      return;
    }

    if (isHash(search)) {
      void getEnvelopes({ variables: { digest: search } }).then((res) => {
        const newResults = res.data?.dsses?.edges?.map((edge) => edge?.node) || [];
        const combined = combineArrays(results, newResults as Dsse[]);
        setResults(combined);
        populateGraph(combined);
      });
    } else {
      void sha256(search)
        .then((hash) => {
          digest = hash;
        })
        .then(() => {
          void getEnvelopes({ variables: { digest: digest } }).then((res) => {
            const newResults = res.data?.dsses?.edges?.map((edge) => edge?.node) || [];
            const combined = combineArrays(results, newResults as Dsse[]);
            setResults(combined);
            populateGraph(combined);
          });
        });
    }
  };

  const getBackRefs = () => {
    backRefs.forEach((backRef) => {
      results.forEach((result) => {
        result.statement?.subjects?.edges?.forEach((subject) => {
          if (subject?.node?.name.includes(backRef)) {
            if (subject?.node?.subjectDigests != null) {
              const digest = subject?.node?.subjectDigests[0]?.value as string;
              void getEnvelopes({ variables: { digest: digest } }).then((res) => {
                const newResults = res.data?.dsses?.edges?.map((edge) => edge?.node) || [];
                const combined = combineArrays(results, newResults as Dsse[]);
                setResults(combined);
              });
            }
          }
        });
      });
    });
  };

  const handleClear = () => {
    setResults([]);
  };

  return (
    <div>
      <div style={{ display: "flex", flexDirection: "row", justifyContent: "center", alignItems: "center" }}>
        <div style={{ backgroundColor: "rgba(255, 255, 255, 0.5)" }}>
          <div style={{ display: "flex", flexDirection: "row", justifyContent: "center", width: "100%" }}>
            <FormControl sx={{ m: 1, width: 500, minWidth: 200, maxWidth: "80%" }}>
              <TextField id="outlined-basic" label="Search" variant="outlined" value={search} onChange={handleChange} />
            </FormControl>
            <FormControl sx={{ mt: 2, paddingRight: 1 }}>
              <Button variant="contained" onClick={handleClick}>
                Search
              </Button>
            </FormControl>
            <FormControl sx={{ mt: 2, paddingRight: 1 }}>
              <Button variant="contained" onClick={handleClear}>
                Clear
              </Button>
            </FormControl>
            <FormControl sx={{ mt: 2 }}>
              <Button variant="contained" onClick={getBackRefs}>
                Get BackRefs
              </Button>
            </FormControl>
            <FormControl sx={{ mt: 0, paddingLeft: 1 }}>
              <Switch checked={useGraph} onChange={() => setUseGraph(!useGraph)} />
              <label>Use Graph</label>
            </FormControl>
          </div>
        </div>
      </div>
      <div>
        <div style={{ display: useGraph ? "block" : "none" }}>
          <DisplayGraph nodes={nodesState} edges={edgesState} />
        </div>
        <div style={{ display: useGraph ? "none" : "block" }}>
          <div>
            {results.map((result) => (
              <div key={result.gitoidSha256} style={{ backgroundColor: "rgba(255, 255, 255, 0.5)" }}>
                <DsseEnvelope key={result.gitoidSha256} result={result} />
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

async function sha256(s: string) {
  const utf8 = new TextEncoder().encode(s);
  const hashBuffer = await crypto.subtle.digest("SHA-256", utf8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map((bytes) => bytes.toString(16).padStart(2, "0")).join("");
  return hashHex;
}

const isHash = (text: string) => {
  // Regular expression to check if string is a SHA256 hash
  const regexExpSHA256 = /^[a-f0-9]{64}$/gi;
  const regexExpSHA1 = /^[a-f0-9]{40}$/gi;

  return regexExpSHA256.test(text) || regexExpSHA1.test(text);
};

export default Search;
