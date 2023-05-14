import * as React from "react";
import { Theme, Typography } from "@mui/material";
import { SxProps } from "@mui/system";

export interface Props {
  sx?: SxProps<Theme>;
}

export function Title({ children, sx }: React.PropsWithChildren<Props>) {
  return (
    <Typography variant="h5" sx={sx ?? {}}>
      {children}
    </Typography>
  );
}
