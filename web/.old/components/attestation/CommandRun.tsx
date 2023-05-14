import * as React from "react";
import { CommandRunAttestor, HashType, ProcessInfo } from "../../lib/Witness";
import { Typography } from "@mui/material";

import { Material } from "./";

export interface CommandRunProps {
  attestation: CommandRunAttestor;
}

export function CommandRun({ attestation }: CommandRunProps) {
  return (
    <>
      <Typography variant="body2">
        <b>Command:</b> {attestation.cmd.join(" ")}
      </Typography>
      <Typography variant="body2">
        <b>Exit code:</b> {attestation.exitcode}
      </Typography>
      <Typography variant="body2">
        <b>Standard out:</b> {attestation.stdout ?? ""}
      </Typography>
      <Typography variant="body2">
        <b>Standard Error:</b> {attestation.stderr ?? ""}
      </Typography>
      {attestation.processes === undefined || attestation.processes.length === 0 ? null : (
        <>
          <Typography variant="body2">
            <b>Processes:</b>
          </Typography>
          <ul>
            {attestation.processes.map((pi) => (
              <ProcessDetails key={pi.processid} processInfo={pi} />
            ))}
          </ul>
        </>
      )}
    </>
  );
}

export interface ProcessInfoProps {
  processInfo: ProcessInfo;
}

function ProcessDetails({ processInfo }: ProcessInfoProps) {
  return (
    <>
      <Typography variant="body2">
        <b>Process ID:</b> {processInfo.processid}
      </Typography>
      <Typography variant="body2">
        <b>Parent Process ID:</b> {processInfo.parentpid}
      </Typography>
      {processInfo.program === undefined ? null : (
        <Typography variant="body2">
          <b>Program:</b> {processInfo.program}
        </Typography>
      )}
      {processInfo.programdigest === undefined ? null : (
        <>
          <Typography variant="body2">
            <b>Program Digest:</b>
          </Typography>
          <ul>
            {Object.keys(processInfo.programdigest).map((digest) => {
              return processInfo.programdigest === undefined ? null : (
                <li key={digest}>
                  <b>{digest}</b>: {processInfo.programdigest[digest as HashType]}
                </li>
              );
            })}
          </ul>
        </>
      )}
      {processInfo.comm === undefined ? null : (
        <Typography variant="body2">
          <b>Comm:</b> {processInfo.comm}
        </Typography>
      )}
      {processInfo.cmdline === undefined ? null : (
        <Typography variant="body2">
          <b>Cmdline:</b> {processInfo.cmdline}
        </Typography>
      )}
      {processInfo.environ === undefined ? null : (
        <Typography variant="body2">
          <b>Environ:</b> {processInfo.environ}
        </Typography>
      )}
      {processInfo.specbypassisvuln === undefined ? null : (
        <Typography variant="body2">
          <b>Spec Bypass Is Vuln:</b> {processInfo.specbypassisvuln ? "TRUE" : "FALSE"}
        </Typography>
      )}
      {processInfo.openedfiles === undefined ? null : (
        <>
          <Typography variant="body2">
            <b>Opened Files:</b>
          </Typography>
          <Material data={processInfo.openedfiles} />
        </>
      )}
    </>
  );
}
