import '@testing-library/jest-dom'; // Import testing-library/jest-dom

import { fireEvent, render } from '@testing-library/react';

import CommitLink from './CommitLink';
import React from 'react';

const mockCopyToClipboard = jest.fn();

describe('CommitLink', () => {
  test('renders the component with the provided commit', () => {
    const commit = 'abcd1234efgh5678';

    const { getByText } = render(<CommitLink commit={commit} copyToClipboard={mockCopyToClipboard} />);

    const commitText = getByText('Commit: abcd123');

    expect(commitText).toBeInTheDocument();
  });

  test('copies the commit text to the clipboard when clicked', () => {
    const commit = 'https://witness.dev/attestations/git/v0.1/commithash:733c37a31f76c20b9d2c706237a210891ea57930';

    const { getByText } = render(<CommitLink commit={commit} copyToClipboard={mockCopyToClipboard} />);
    const commitText = getByText('Commit: 733c37a');

    fireEvent.click(commitText);

    expect(mockCopyToClipboard).toHaveBeenCalledWith('733c37a31f76c20b9d2c706237a210891ea57930');
  });

  test('displays "Commit copied!" message when clicked', () => {
    const commit = 'https://witness.dev/attestations/git/v0.1/commithash:733c37a31f76c20b9d2c706237a210891ea57930';

    const { getByText } = render(<CommitLink commit={commit} copyToClipboard={mockCopyToClipboard} />);
    const commitText = getByText('Commit: 733c37a');

    fireEvent.click(commitText);

    const message = getByText('Commit copied!');
    expect(message).toBeInTheDocument();
  });
});
