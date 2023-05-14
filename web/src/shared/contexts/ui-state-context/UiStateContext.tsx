import React, { PropsWithChildren, createContext, useContext, useState } from 'react';

interface UiState {
  isSearchOpen: boolean;
}

interface UiStateContextProps {
  uiState: UiState;
  setUiState: (uiState: UiState) => void;
}

export const UiStateContext = createContext<UiStateContextProps>({
  uiState: { isSearchOpen: false },
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUiState: () => {},
});

/**
- A context for managing dynamic UI state, globally.
- This context allows us to put UI state that is important globally into a single place for easy use.
- It should be used sparingly and only when local state hooks are not sufficient for managing the UI state.
- Note that in a GatsbyJS app, the UiStateContext is not available at build time due to shifting-left.
- Therefore, avoid using the UiStateContext when generating static content (e.g. static HTML pages).
- @example 
- To use the UiStateContext, just import and call the provided useUiState hook. The Context has already been wrapped in our App.tsx
*/
export const UiStateProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const [uiState, setUiState] = useState<UiState>({ isSearchOpen: false });

  return <UiStateContext.Provider value={{ uiState, setUiState }}>{children}</UiStateContext.Provider>;
};

export const useUiState = () => useContext(UiStateContext);
