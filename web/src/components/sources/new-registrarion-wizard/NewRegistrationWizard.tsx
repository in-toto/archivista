import { Box, Button, MenuItem, Step, StepLabel, Stepper, TextField, Typography } from '@mui/material';
import React, { useState } from 'react';

enum NodeType {
  TPM = 'TPM',
  GCP = 'Google Cloud',
  AWS = 'AWS',
  Azure = 'Azure',
  GitHub = 'GitHub',
}

enum NodeGroup {
  Dev = 'dev',
  Prod = 'prod',
}

enum SelectorKey {
  Repository = 'Repository',
  ProjectID = 'ProjectId',
  PublicKeyHash = 'Public Key Hash',
}

const nodeTypeSelectorKeys: Record<NodeType, SelectorKey> = {
  [NodeType.TPM]: SelectorKey.ProjectID,
  [NodeType.GCP]: SelectorKey.PublicKeyHash,
  [NodeType.AWS]: SelectorKey.Repository,
  [NodeType.Azure]: SelectorKey.Repository,
  [NodeType.GitHub]: SelectorKey.Repository,
};

const steps = ['Choose Node Type', 'Choose Node Group', 'Choose Selector Key', 'Input Selector Value'];

const NewRegistrationWizard = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [nodeType, setNodeType] = useState<NodeType>(NodeType.TPM);
  const [nodeGroup, setNodeGroup] = useState<NodeGroup>(NodeGroup.Dev);
  const [selectorKey, setSelectorKey] = useState<SelectorKey>(nodeTypeSelectorKeys[nodeType]);
  const [selectorValue, setSelectorValue] = useState('');

  const handleNext = () => setActiveStep((prevActiveStep) => prevActiveStep + 1);
  const handleBack = () => setActiveStep((prevActiveStep) => prevActiveStep - 1);
  const handleReset = () => setActiveStep(0);

  const handleNodeTypeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedNodeType = event.target.value as NodeType;
    setNodeType(selectedNodeType);
    setSelectorKey(nodeTypeSelectorKeys[selectedNodeType]);
  };

  const renderStepContent = (stepIndex: number) => {
    switch (stepIndex) {
      case 0:
        return (
          <TextField select fullWidth label="Node Type" value={nodeType} onChange={handleNodeTypeChange}>
            {Object.values(NodeType).map((type) => (
              <MenuItem key={type} value={type}>
                {type}
              </MenuItem>
            ))}
          </TextField>
        );
      case 1:
        return (
          <TextField select fullWidth label="Node Group" value={nodeGroup} onChange={(event) => setNodeGroup(event.target.value as NodeGroup)}>
            {Object.values(NodeGroup).map((group) => (
              <MenuItem key={group} value={group}>
                {group}
              </MenuItem>
            ))}
          </TextField>
        );
      case 2:
        return (
          <TextField select fullWidth label="Selector Key" value={selectorKey} onChange={(event) => setSelectorKey(event.target.value as SelectorKey)}>
            {Object.values(SelectorKey).map((key) => (
              <MenuItem key={key} value={key}>
                {key}
              </MenuItem>
            ))}
          </TextField>
        );
      case 3:
        return <TextField fullWidth label="Selector Value" value={selectorValue} onChange={(event) => setSelectorValue(event.target.value)} />;
      default:
        return null;
    }
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mt: 3 }}>
      <Stepper activeStep={activeStep}>
        {steps.map((label, i) => (
          <Step key={i}>
            <StepLabel>
              {activeStep > 0 && i === 0
                ? nodeType
                : activeStep > 1 && i === 1
                ? nodeGroup
                : activeStep > 2 && i === 2
                ? selectorKey
                : activeStep > 3 && i === 3
                ? selectorValue
                : label}
            </StepLabel>
          </Step>
        ))}
      </Stepper>
      {activeStep === steps.length ? (
        <>
          <Typography variant="h5" gutterBottom>
            All steps completed
          </Typography>
          <Typography variant="subtitle1">Review your selections and submit.</Typography>
          <Button onClick={handleReset} variant="outlined" sx={{ mt: 3 }}>
            Reset
          </Button>
        </>
      ) : (
        <>
          {renderStepContent(activeStep)}
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 3 }}>
            {activeStep !== 0 && (
              <Button onClick={handleBack} sx={{ mr: 1 }}>
                Back
              </Button>
            )}
            <Button onClick={handleNext} variant="contained">
              {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
            </Button>
          </Box>
        </>
      )}
    </Box>
  );
};

export default NewRegistrationWizard;
