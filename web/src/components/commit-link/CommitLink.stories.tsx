import CommitLink from './CommitLink';
import React from 'react';
import { storiesOf } from '@storybook/react';

const commit = 'asdfbu892jlkj23';

// eslint-disable-next-line @typescript-eslint/no-empty-function
storiesOf('CommitLink', module).add('default', () => <CommitLink commit={commit} copyToClipboard={() => {}} />);
