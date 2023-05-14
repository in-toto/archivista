import { Paper, Stack, styled } from "@mui/material";
import React, { useCallback, useContext, useEffect, useState } from "react";

import { NodeRegistration } from "../lib/registrar/registrar.pb";
import { Project } from "../lib/Gitlab";
import { RegistrarContext } from "../App";
import { useFeatureToggles } from "../shared/contexts/FeatureToggleContext/FeatureToggleContext";

export const Dashboard = (): JSX.Element => {
  const featureToggles = useFeatureToggles();
  const dashboardFeature = featureToggles?.features?.find?.((f) => f.name === "Judge Dashboard");
  const [sources, setSources] = useState<NodeRegistration[]>([]);
  const registrar = useContext(RegistrarContext);

  const fetchRegistrations = useCallback(async () => {
    await registrar.GetNodeRegistrations().then((response) => {
      setSources(response?.nodeRegistrations || []);
    });
  }, [setSources, registrar]);

  useEffect(() => {
    void fetchRegistrations();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const dashboard = dashboardFeature?.enabled && (
    <Stack spacing={2}>
      {sources.map((project, i) => (
        <ProjectItem key={i} project={{ id: i, name: project.nodeId || "", name_with_namespace: project.nodeGroup || "" }} />
      ))}
    </Stack>
  );
  return <>{dashboard}</>;
};

interface ProjectItemProps {
  project: Project;
  onClick?: (project: Project) => void;
}

const StyledItem = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
  cursor: "pointer",
}));

function ProjectItem({ project, onClick }: ProjectItemProps): JSX.Element {
  const clickHandler = () => {
    if (onClick === undefined) {
      return;
    }

    onClick(project);
  };

  return <StyledItem onClick={clickHandler}>{project.name_with_namespace}</StyledItem>;
}
