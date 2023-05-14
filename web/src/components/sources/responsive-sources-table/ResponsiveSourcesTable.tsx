/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { IconButton, Skeleton, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';

import { DateTime } from 'luxon';
import { DeleteOutline as DeleteOutlineIcon } from '@mui/icons-material';
import { NodeRegistration } from '../../../shared/hooks/use-registrar/useRegistrar';
import React from 'react';

interface TableProps {
  rows: NodeRegistration[];
  deactiveNode: (nodeId?: string) => void;
  isLoading: boolean; // New prop for loading state
}

const ResponsiveSourcesTable: React.FC<TableProps> = ({ rows, deactiveNode, isLoading }) => {
  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow sx={{ backgroundColor: 'lightgray' }}>
            <TableCell sx={{ width: '30px', textAlign: 'center' }}></TableCell>
            <TableCell sx={{ width: '80px', textAlign: 'center' }}>Source Type</TableCell>
            <TableCell sx={{ width: '80px', textAlign: 'center' }}>Source Group</TableCell>
            <TableCell sx={{ width: '80px', textAlign: 'center' }}>Registered Date</TableCell>
            <TableCell sx={{ width: '80px', textAlign: 'center' }}>User</TableCell>
            <TableCell sx={{ width: '250px', textAlign: 'center' }}>Selector</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {isLoading ? (
            // Render a skeleton when isLoading is true
            <TableRow>
              <TableCell colSpan={6} sx={{ textAlign: 'center' }}>
                <Skeleton variant="rectangular" width="100%" height={50} />
              </TableCell>
            </TableRow>
          ) : (
            // Render the table rows when isLoading is false
            rows.map((row, i) => (
              <TableRow key={i}>
                <TableCell sx={{ textAlign: 'center' }}>
                  <IconButton title="Delete" onClick={() => deactiveNode(row?.nodeId || '')}>
                    <DeleteOutlineIcon />
                  </IconButton>
                </TableCell>
                <TableCell sx={{ textAlign: 'center' }}>{row.nodeType}</TableCell>
                <TableCell sx={{ textAlign: 'center' }}>{row.nodeGroup}</TableCell>
                {row?.registeredAt && <TableCell sx={{ textAlign: 'center' }}>{DateTime.fromISO(row?.registeredAt).toLocaleString()}</TableCell>}
                <TableCell sx={{ textAlign: 'center' }}>{row.registeredBy}</TableCell>
                <TableCell sx={{ textAlign: 'center' }}>
                  {row.selectors?.selectors && row.selectors.selectors.length > 0
                    ? // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
                      `${row?.selectors?.selectors[0]?.key}:${row?.selectors?.selectors[0]?.value}`
                    : '-'}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default ResponsiveSourcesTable;
