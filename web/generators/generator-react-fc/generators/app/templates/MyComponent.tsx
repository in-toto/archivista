import React from 'react';

export type <%= componentName %>Props = {
  message?: string;
};

const <%= componentName %> = ({ message }: <%= componentName %>Props) => {
  return <>{message}</>;
};

export default <%= componentName %>;
