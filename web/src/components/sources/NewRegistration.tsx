import * as React from 'react';

import { Button, TextField } from '@mui/material';
import { IRegistrarService, RegisterNodeRequest } from '../../shared/hooks/use-registrar/useRegistrar';

import FormControl from '@mui/material/FormControl';
import FormHelperText from '@mui/material/FormHelperText';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';

//props for the NewRegistration component
interface NewRegistrationProps {
  registrar?: IRegistrarService;
  fetchRegistrations: () => void;
}

export const NewRegistration = (props: NewRegistrationProps) => {
  const [NodeType, setNodeType] = React.useState('');
  const [NodeGroup, setNodeGroup] = React.useState('');
  const [SelectorKey, setSelectorKey] = React.useState('');
  const [SelectorValue, setSelectorValue] = React.useState('');
  const [Groups, setGroups] = React.useState([] as string[]);
  // const auth = useAuth();
  // const ac = auth?.authService;

  // React.useEffect(() => {
  //   void ac?.getGroups?.().then((groups) => {
  //     setGroups(groups);
  //   });
  // }, [ac]);

  const handleSubmit = () => {
    const req = {
      nodeType: NodeType,
      nodeGroup: NodeGroup,
      selectors: {
        selectors: [
          {
            key: SelectorKey,
            value: SelectorValue,
          },
        ],
      },
    } as RegisterNodeRequest;
    props.registrar
      ?.registerNode?.(req)
      ?.then((response) => {
        console.log(response);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        props.fetchRegistrations();
      });
  };

  return (
    <div>
      <div>
        New Source<br></br>
        <FormControl sx={{ m: 1, minWidth: 150 }}>
          <Select id="node-type" value={NodeType} onChange={(e) => setNodeType(e.target.value.toUpperCase())}>
            <MenuItem value={'TPM'}>TPM</MenuItem>
            <MenuItem value={'GCP'}>Google Cloud</MenuItem>
            <MenuItem value={'AWS'}>AWS</MenuItem>
            <MenuItem value={'AZURE'}>Azure</MenuItem>
            <MenuItem value={'GITHUB'}>GitHub</MenuItem>
          </Select>
          <FormHelperText>Node Type</FormHelperText>
        </FormControl>
        <FormControl sx={{ m: 1, minWidth: 150 }}>
          <Select label="Node Group" id="node-group" value={NodeGroup} onChange={(e) => setNodeGroup(e.target.value)}>
            {Groups?.map((group) => {
              return (
                <MenuItem key={group} value={group}>
                  {group}
                </MenuItem>
              );
            })}
          </Select>
          <FormHelperText>Node Group</FormHelperText>
        </FormControl>
        <FormControl sx={{ m: 1, minWidth: 150 }}>
          <Select label="Selector Key" id="selector-key" value={SelectorKey} onChange={(e) => setSelectorKey(e.target.value)}>
            {NodeType === 'TPM' && <MenuItem value={'pub_hash'}>Public Key Hash</MenuItem>}
            {NodeType === 'GCP' && <MenuItem value={'project-id'}>Project ID</MenuItem>}
            {NodeType === 'GITHUB' && <MenuItem value={'repo'}>Repository</MenuItem>}
          </Select>
          <FormHelperText>Selector Key</FormHelperText>
        </FormControl>
        <FormControl sx={{ m: 1, minWidth: 400 }}>
          <TextField id="selector-value" value={SelectorValue} onChange={(e) => setSelectorValue(e.target.value)} />
          <FormHelperText>Selector Value</FormHelperText>
        </FormControl>
        <FormControl sx={{ mt: 2, minWidth: 150 }}>
          <Button variant="contained" onClick={handleSubmit} disabled={NodeType === '' || NodeGroup === '' || SelectorKey === '' || SelectorValue === ''}>
            Submit
          </Button>
        </FormControl>
      </div>
    </div>
  );
};
