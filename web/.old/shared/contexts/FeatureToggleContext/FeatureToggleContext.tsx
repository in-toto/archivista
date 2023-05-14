import React, { createContext, useContext, useState } from "react";
// TODO: Load the feature toggles from a platform, not a const...
import FeatureToggles from "./feature-toggles";

export type FeatureToggle = {
  name: string;
  enabled: boolean;
  notes?: string;
};

interface FeatureToggleContextProps {
  features: FeatureToggle[];
  setFeatures: React.Dispatch<React.SetStateAction<FeatureToggle[]>>;
}

export const FeatureTogglesContext = createContext<FeatureToggleContextProps>({
  features: [],
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setFeatures: () => {},
});

interface FeatureToggleProviderProps {
  children: React.ReactNode;
}

/**
 * Provides a context for FeatureToggles
 *
 * @param {*} { children }
 * @returns
 */
export const FeatureToggleProvider: React.FC<FeatureToggleProviderProps> = ({ children }) => {
  // TODO: load the toggles from a platform
  const [features, setFeatures] = useState<FeatureToggle[]>(FeatureToggles);

  return <FeatureTogglesContext.Provider value={{ features, setFeatures }}>{children}</FeatureTogglesContext.Provider>;
};

export const useFeatureToggles = () => {
  return useContext(FeatureTogglesContext);
};
