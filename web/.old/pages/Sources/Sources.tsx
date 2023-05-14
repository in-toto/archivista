import * as React from "react";

import { Box, Button, Fab, Modal, Skeleton, Typography } from "@mui/material";
import { DecativateNodeRequest, NodeRegistration } from "../../lib/registrar/registrar.pb";

import AddIcon from "@mui/icons-material/Add";
import { DateTime } from "luxon";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import EmptyState from "../../shared/components/states/EmptyState/EmptyState";
import { NewRegistration } from "../../components/sources/NewRegistration";
import NewRegistrationWizard from "../../components/sources/NewRegistrationWizard/NewRegistrationWizard";
import { RegistrarContext } from "../../App";
import ResponsiveSourcesTable from "./ResponsiveSourcesTable/ResponsiveSourcesTable";
import { useFeatureToggles } from "../../shared/contexts/FeatureToggleContext/FeatureToggleContext";

//TODO: refactor, two components.
export const Sources = () => {
  const featureToggles = useFeatureToggles();
  const newSourceExperienceToggle = featureToggles?.features?.find?.((f) => f.name === "New Source Experience");
  const [isLoading, setIsLoading] = React.useState(true);
  const [rows, setRows] = React.useState<NodeRegistration[]>([]);
  const registrar = React.useContext(RegistrarContext);
  const [isWizardOpen, setIsWizardOpen] = React.useState(false);

  const deactiveNode = (nodeId?: string) => {
    registrar
      .DeactivateNode({
        nodeId: nodeId,
      } as DecativateNodeRequest)
      .then((response) => {
        console.log(response);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        void fetchRegistrations();
      });
  };

  const fetchRegistrations = async () => {
    setIsLoading(true); // set isLoading to true when starting to fetch data
    await registrar
      .GetNodeRegistrations()
      .then((response) => {
        setRows(response?.nodeRegistrations || []);
      })
      .finally(() => {
        setIsLoading(false); // set isLoading to false when done fetching data
      });
  };

  React.useEffect(() => {
    void fetchRegistrations();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const oldAddSourceExperience = <NewRegistration registrar={registrar} fetchRegistrations={() => void fetchRegistrations()} />;

  const newAddSourceExperience = (
    <>
      <Fab
        color="primary"
        aria-label="add source"
        onClick={() => setIsWizardOpen(true)}
        style={{
          position: "fixed",
          bottom: "16px",
          right: "16px",
          zIndex: 1,
        }}
      >
        <AddIcon />
      </Fab>
      <Modal
        open={isWizardOpen}
        onClose={() => setIsWizardOpen(false)}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Box
          sx={{
            position: "absolute" as const,
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            minWidth: 400,
            maxWidth: 900,
            bgcolor: "background.paper",
            border: "2px solid #000",
            boxShadow: 24,
            p: 4,
          }}
        >
          <Typography id="modal-modal-title" variant="h6" component="h2">
            Add a Source
          </Typography>
          <Typography id="modal-modal-description" sx={{ mt: 2 }}>
            Configure and add a new Source to Judge.
          </Typography>
          <NewRegistrationWizard
          //  registrar={registrar}
          //fetchRegistrations={() => void fetchRegistrations()}
          />
        </Box>
      </Modal>
    </>
  );

  const emptyState = (
    <EmptyState
      headerText="Your sources are empty"
      subheaderText="Please add a source to get started!"
      addText="Add a source"
      onAddButtonClick={() => setIsWizardOpen(newSourceExperienceToggle?.enabled || false)}
    />
  );

  const oldTableExperience = (
    <table style={{ width: "100%" }}>
      <thead style={{ backgroundColor: "lightgray" }}>
        <tr>
          <th style={{ width: "30px", textAlign: "center" }}></th>
          <th style={{ width: "80px", textAlign: "center" }}>Source Type</th>
          <th style={{ width: "80px", textAlign: "center" }}>Source Group</th>
          <th style={{ width: "80px", textAlign: "center" }}>Registered Date</th>
          <th style={{ width: "80px", textAlign: "center" }}>User</th>
          <th style={{ width: "250px", textAlign: "center" }}>Selector</th>
        </tr>
      </thead>
      <tbody>
        {isLoading ? (
          <tr>
            <td colSpan={6}>
              <div
                style={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                  padding: "16px",
                }}
              >
                <Skeleton animation="wave" width={120} height={20} style={{ marginRight: 8 }} />
                <Skeleton animation="wave" width={80} height={20} style={{ marginRight: 8 }} />
                <Skeleton animation="wave" width={80} height={20} style={{ marginRight: 8 }} />
                <Skeleton animation="wave" width={120} height={20} style={{ marginRight: 8 }} />
                <Skeleton animation="wave" width={60} height={20} style={{ marginRight: 8 }} />
                <Skeleton animation="wave" width={200} height={20} style={{ marginRight: 8 }} />
              </div>
            </td>
          </tr>
        ) : rows.length > 0 ? (
          rows.map((row) => (
            <tr key={row.nodeId}>
              <td style={{ textAlign: "center" }}>
                <Button onClick={() => deactiveNode(row.nodeId as string)}>
                  <DeleteOutlineIcon />
                </Button>
              </td>
              <td style={{ textAlign: "center" }}>{row.nodeType}</td>
              <td style={{ textAlign: "center" }}>{row.nodeGroup}</td>
              <td style={{ textAlign: "center" }}>{DateTime.fromISO(row.registeredAt as string).toLocaleString()}</td>
              <td style={{ textAlign: "center" }}>{row.registeredBy}</td>
              {row?.selectors?.selectors && row?.selectors?.selectors.length > 0 ? (
                <td>
                  {row.selectors.selectors.map((selector, index) => (
                    <React.Fragment key={selector.key}>
                      {index > 0 && <br />}
                      <code>{selector.value}</code>
                    </React.Fragment>
                  ))}
                </td>
              ) : (
                <td style={{ textAlign: "center" }}>-</td>
              )}
            </tr>
          ))
        ) : (
          <tr>
            <td colSpan={6}>{emptyState}</td>
          </tr>
        )}
      </tbody>
    </table>
  );

  const newTableExperience = rows.length > 0 ? <ResponsiveSourcesTable rows={rows} deactiveNode={deactiveNode} isLoading={isLoading} /> : emptyState;

  return (
    <div>
      {newSourceExperienceToggle?.enabled ? newAddSourceExperience : oldAddSourceExperience}
      {newSourceExperienceToggle?.enabled ? newTableExperience : oldTableExperience}
    </div>
  );
};
