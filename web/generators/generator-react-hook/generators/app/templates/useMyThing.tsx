import { useState } from 'react';

export type <%= propsName %> = {
  isLoading: boolean;
};
/**
 * <%= desc %>
 *
 */
const <%= hookName %> = (): [<%= propsName %>, (isLoading: boolean) => void] => {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const <%= propsName %>: <%= propsName %> = { isLoading };

  return [<%= propsName %>, setIsLoading];
};

export default <%= hookName %>;
