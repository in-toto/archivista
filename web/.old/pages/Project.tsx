import * as React from "react";
import { Accordion, AccordionDetails, AccordionSummary, Paper, Stack, Typography, styled } from "@mui/material";
import { Link, useNavigate } from "react-router-dom";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

import { ApiService as GitlabApiService, Pipeline, Project } from "../lib/Gitlab";
import { LogEntry, ApiService as RekorApiService } from "../lib/Rekor";
import { AttestationCollection } from "../lib/Witness";
import { Statement } from "../lib/InToto";

export interface Props {
  project: Project | null;
  gitlabService: GitlabApiService | null;
  rekorService: RekorApiService;
}

const SectionPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
  display: "flex",
  flexDirection: "column",
  gap: theme.spacing(2),
}));

export function Project({ project, gitlabService, rekorService }: Props): JSX.Element {
  const navigate = useNavigate();
  if (project === null) {
    navigate("/");
  }

  const [pipelines, setPipelines] = React.useState<Pipeline[]>([]);
  React.useEffect(() => {
    if (project === null || gitlabService === null) {
      return;
    }

    void (async () => {
      setPipelines(await gitlabService.fetchPipelines(project.id));
    })();
  }, [setPipelines, project, gitlabService]);

  return project === null ? (
    <></>
  ) : (
    <>
      <Typography variant="h6" sx={{ pb: 2 }}>
        {project.name_with_namespace}
      </Typography>
      <SectionPaper>
        <Typography variant="body1">Pipelines</Typography>
        <Pipelines rekorService={rekorService} pipelines={pipelines} />
      </SectionPaper>
    </>
  );
}

interface PipelineAccordionProps {
  pipelines: Pipeline[];
  rekorService: RekorApiService;
}

function Pipelines({ pipelines, rekorService }: PipelineAccordionProps): JSX.Element {
  return (
    <>
      {pipelines?.map((pipeline) => (
        <Accordion key={pipeline.id}>
          <AccordionSummary key={pipeline.id} expandIcon={<ExpandMoreIcon />}>
            <Typography>{pipeline.ref}</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <PipelineDetails pipeline={pipeline} rekorService={rekorService} />{" "}
          </AccordionDetails>
        </Accordion>
      ))}
    </>
  );
}

interface PipelineDetailProps {
  pipeline: Pipeline;
  rekorService: RekorApiService;
}

function PipelineDetails({ pipeline, rekorService }: PipelineDetailProps): JSX.Element {
  const [logEntries, setLogEntries] = React.useState<Record<string, LogEntry>>({});
  React.useEffect(() => {
    void (async () => {
      const foundEntries = await rekorService.searchIndexByString(pipeline.web_url);
      const entries: Record<string, LogEntry> = {};
      for (const entryId of foundEntries) {
        const entry = await rekorService.getLogEntryByUUID(entryId);
        entries[entryId] = entry;
      }

      setLogEntries(entries);
    })();
  }, [pipeline, rekorService, setLogEntries]);

  if (Object.keys(logEntries).length <= 0) {
    return <Typography>No attestations found for this pipeline</Typography>;
  }

  return (
    <Stack>
      {Object.entries(logEntries).map((objEntry) => {
        const [uuid, entry] = objEntry;
        const decodedAttestation = atob(entry.attestation.data);
        const statement = JSON.parse(decodedAttestation) as Statement<AttestationCollection>;
        return (
          <Link to={`/attestation/${uuid}`} key={uuid}>
            Attestation found for {statement.predicate.name}
          </Link>
        );
      })}
    </Stack>
  );
}
