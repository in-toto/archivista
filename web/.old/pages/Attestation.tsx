import * as React from "react";
import { Link, useParams } from "react-router-dom";
import { useEffect, useState } from "react";

import { default as Tab, TabProps } from "@mui/material/Tab";
import Box from "@mui/material/Box";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Tabs from "@mui/material/Tabs";
import Typography from "@mui/material/Typography";

import { Attestation, AttestationCollection } from "../lib/Witness";
import { Modal as AttestationModal } from "../components/attestation";
import { ApiService as RekorApiService } from "../lib/Rekor";
import { Statement } from "../lib/InToto";

export interface Props {
  rekorService: RekorApiService;
}

export function Attestation({ rekorService }: Props): JSX.Element {
  const { entryId } = useParams<string>();
  const decode = (str: string): string => Buffer.from(str, "base64").toString("binary");
  const [statement, setAttestation] = useState<Statement<AttestationCollection>>();
  const [showSpinner, setSpinner] = useState<boolean>(true);
  const breadcrumbs = [
    <Link to={"/"} key="1" color="inherit">
      Dashboard
    </Link>,
    <Link to={"/project"} key="2" color="inherit">
      Pipelines
    </Link>,
    <Typography key="3" color="text.primary">
      {statement?.predicate.name}
    </Typography>,
  ];

  // Tabs
  const [value, setValue] = React.useState(0);
  const handleChange = (_: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  // Download all attestations
  const exportToJson = () => {
    if (entryId === undefined) {
      return;
    }

    const jsonString = `data:text/json;chatset=utf-8,${encodeURIComponent(JSON.stringify(statement))}`;
    const link = document.createElement("a");
    link.href = jsonString;
    link.download = `${entryId}.json`;

    link.click();
  };

  useEffect(() => {
    if (entryId === undefined) {
      return;
    }

    void (async () => {
      const data = decode((await rekorService.getLogEntryByUUID(entryId)).attestation.data);
      const json = JSON.parse(data) as Statement<AttestationCollection>;
      setAttestation(json);
      setSpinner(false);
    })();
  }, [setAttestation, setSpinner, entryId, rekorService]);

  return showSpinner === true ? (
    <Box sx={{ height: "500px", margin: 0, display: "flex", justifyContent: "center", alignItems: "center" }}>
      <CircularProgress size="100px" />
    </Box>
  ) : (
    <>
      <Stack spacing={2}>
        <Breadcrumbs separator="›" aria-label="breadcrumb">
          {breadcrumbs}
        </Breadcrumbs>
      </Stack>
      <hr />
      <Typography variant="h5" mt={2} gutterBottom component="div">
        Attestions for {statement?.predicate.name}
      </Typography>

      <Grid container spacing={2} sx={{ borderBottom: "1px solid #dbdbdb" }}>
        <Grid item xs={10}>
          <Tabs value={value} onChange={handleChange} aria-label="nav tabs example">
            <LinkTab label="Attestations" href="" />
          </Tabs>
        </Grid>
        <Grid item xs={2}>
          <Button variant="outlined" onClick={exportToJson} sx={{ textTransform: "none" }}>
            Download all
          </Button>
        </Grid>
      </Grid>

      <TableContainer component={Paper} sx={{ boxShadow: "none" }}>
        <Table size="small" aria-label="collapsible table" sx={{ border: "none" }}>
          <TableHead>
            <TableRow>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>
            {statement?.predicate.attestations.map((row: Attestation) => (
              <Row key={row.type} row={row} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
}

function Row(props: { row: Attestation }) {
  const { row } = props;

  return (
    <React.Fragment>
      <TableRow
        hover
        sx={{
          "&.MuiTableRow-root:hover": {
            backgroundColor: "#E7FBD3",
          },
        }}
      >
        <TableCell sx={{ borderBottom: "none" }}>
          <AttestationModal data={row} />
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

function LinkTab(props: TabProps<"a">) {
  return (
    <Tab
      component="a"
      onClick={(event: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
        event.preventDefault();
      }}
      {...props}
    />
  );
}
