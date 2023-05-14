# Feature Toggles

## What are Feature Toggles, anyway?

Feature toggles, also known as feature flags, are a software development technique that allows developers to turn on or off certain features or functionalities within an application. This can provide numerous benefits, including:

- Continuous delivery: With feature toggles, developers can release new features or updates without disrupting existing functionality. This allows for more frequent and smaller releases, which can result in faster feedback, quicker bug fixes, and a smoother development process overall.

- A/B testing: Feature toggles allow developers to test new features with a subset of users before rolling them out to everyone. This can help identify and fix issues early on, as well as gather feedback from users before committing to a full release.

- Risk mitigation: By using feature toggles, developers can minimize the risk of introducing bugs or breaking changes by gradually rolling out new features to a subset of users before releasing them to everyone. This can also help avoid downtime and other issues that may arise from unexpected changes.

- Customization: Feature toggles allow developers to create custom experiences for different users or user groups based on their needs or preferences. This can help improve user engagement and satisfaction by providing tailored experiences.

Overall, feature toggles provide developers with greater control, flexibility, and agility in software development, ultimately leading to better user experiences and more successful applications.

## FeatureTogglesContext

The `FeatureTogglesContext` provides access to an array of feature toggles.

Someday we should be using FeatureToggles from a cms or personalization platform,
so rather than just use a simple TS const, we're abstracting it with this Context so we can easily connect it to a platform someday.
This will also easily allow us to connect these toggles to a UI internally for now.

The context value includes the following properties:

- `featureToggles`: An array of feature toggle objects, where each object has the following properties:
- `name`: A string representing the name of the feature toggle.
- `enabled`: A boolean indicating whether the feature toggle is currently enabled.
- `notes` (optional): A string containing additional notes about the feature toggle.

- `setFeatureToggles`: A function that takes an array of feature toggle objects and updates the context's `featureToggles` property with the new array.

Today, the `FeatureTogglesProvider` component loads a typescript consts file containing feature toggles and provides the `FeatureTogglesContext` to its children.
The provider must wrap any components that consume the `FeatureTogglesContext`.

To use the `FeatureTogglesContext` and `FeatureTogglesProvider`, follow these steps:

1. Use the `useContext` hook to access the `featureToggles` array from the `FeatureTogglesContext`.
1. For a shortcut, you can also just use the `useFeatureToggles` hook we've provided.

```tsx
import React, { useContext } from "react";
import { FeatureTogglesProvider, FeatureTogglesContext } from "./FeatureTogglesProvider";

const App: React.FC = () => {
  return (
    <FeatureTogglesProvider>
      <MyComponent />
    </FeatureTogglesProvider>
  );
};

const MyComponent: React.FC = () => {
  // Access the `featureToggles` array from the `FeatureTogglesContext`.
  const { featureToggles } = useContext(FeatureTogglesContext);

  return (
    <div>
      {featureToggles.map((featureToggle) => (
        <div key={featureToggle.name}>
          <h2>{featureToggle.name}</h2>
          <p>{`Enabled: ${featureToggle.enabled}`}</p>
          {featureToggle.notes && <p>{featureToggle.notes}</p>}
        </div>
      ))}
    </div>
  );
};
```

The `feature-toggles.ts` file should contain an array of feature toggles.
Each feature toggle should be represented as an object with `name`, `enabled`, and `notes` properties.

Here's an example json file:

```ts
export const features = [
  {
    name: "Feature A",
    enabled: true,
    notes: "Notes are optional",
  },
  {
    name: "Feature B",
    enabled: false,
  },
];
```
