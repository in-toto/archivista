import React from 'react';
import { Attestation } from '../../models/Witness';
import { IconButton, Typography, Box } from '@mui/material';
import { Close, FileCopy } from '@mui/icons-material';
import MuiModal from '@mui/material/Modal';

const style = {
  position: 'absolute' as const,
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: '60%',
  height: '60%',
  display: 'block',
  overflow: 'scroll',
  bgcolor: 'background.paper',
  boxShadow: 24,
  borderRadius: '4px',
  outline: 0,
  fontSize: '0.875rem',
  ul: {
    margin: '0',
  },
};

interface AttestationModalProps {
  data: Attestation | undefined;
  open: boolean;
  close: () => void;
}

export function Modal({ data, open, close }: AttestationModalProps) {
  if (!open) {
    return null;
  }

  const getAttestation = () => {
    if (!data) {
      return null;
    }
    return (
      <pre style={{ padding: '16px', whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
        <code style={{ padding: '0px 8px' }}>{JSON.stringify(data, undefined, 2)}</code>
      </pre>
    );
  };

  return (
    <MuiModal open={open} aria-labelledby="modal-modal-title" aria-describedby="modal-modal-description">
      <Box sx={style} style={{ overflow: 'hidden' }}>
        <header
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: '16px',
            borderBottom: '1px solid #e0e0e0',
          }}
        >
          <Typography id="modal-modal-title" variant="h6" component="h2">
            {data?.type}
          </Typography>
          <div>
            <IconButton
              onClick={() => {
                void navigator.clipboard.writeText(data ? JSON.stringify(data, undefined, 2) : '');
              }}
            >
              <FileCopy />
            </IconButton>
            <IconButton onClick={close}>
              <Close />
            </IconButton>
          </div>
        </header>
        <Typography id="modal-modal-description" sx={{ mt: 2 }}></Typography>
        <Box sx={{ height: 'calc(100% - 100px)', overflowY: 'scroll', border: '1px solid #e0e0e0', padding: '0px', width: '100%', marginBottom: '0px' }}>
          {getAttestation()}
        </Box>
      </Box>
    </MuiModal>
  );
}
