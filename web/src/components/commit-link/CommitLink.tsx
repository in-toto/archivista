import React, { useState } from 'react';
import { Tooltip, Typography } from '@mui/material';

export type CommitLinkProps = {
  commit: string;
  copyToClipboard: (text: string) => void;
};

const CommitLink = ({ commit, copyToClipboard }: CommitLinkProps) => {
  const [showMessage, setShowMessage] = useState(false);
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

  const handleCopyToClipboard = () => {
    const exactCommit = commit.split(':').pop();
    if (exactCommit) {
      copyToClipboard(exactCommit);
      setShowMessage(true);
      setTimeout(() => {
        setShowMessage(false);
      }, 2000);
    }
  };

  const handleMouseMove = (event: React.MouseEvent<HTMLSpanElement>) => {
    const { clientX, clientY } = event;
    setMousePosition({ x: clientX, y: clientY });
  };

  return (
    <>
      <Tooltip title="📄 Copy commit to clipboard">
        <Typography
          variant="subtitle2"
          sx={{
            fontWeight: 'bold',
            cursor: 'pointer',
            '&:hover': {
              textDecoration: 'underline',
            },
          }}
          onClick={handleCopyToClipboard}
          onMouseMove={handleMouseMove}
        >
          Commit: {commit.split(':').pop()?.slice(0, 7)}
        </Typography>
      </Tooltip>
      {showMessage && (
        <div
          style={{
            position: 'fixed',
            top: mousePosition.y + 20,
            left: mousePosition.x - 40,
            padding: '8px',
            backgroundColor: 'rgba(0, 0, 0, 0.8)',
            color: '#fff',
            borderRadius: '4px',
            zIndex: 9999,
          }}
        >
          Commit copied!
        </div>
      )}
    </>
  );
};

export default CommitLink;
