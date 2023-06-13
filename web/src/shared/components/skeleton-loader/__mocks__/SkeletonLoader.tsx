import React from 'react';

const SkeletonLoader: React.FC = ({ message }: any) => <div data-testid="SkeletonLoader">{message}</div>;

export default SkeletonLoader;
