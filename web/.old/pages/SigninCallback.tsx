import { IAuthService } from "../lib/AuthService";
import React from "react";
import { User } from "oidc-client";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export interface Props {
  authService: IAuthService;
  setUser: (user: User) => void;
}

export function SigninCallback({ authService, setUser }: Props): JSX.Element {
  const navigate = useNavigate();

  useEffect(() => {
    void (async () => {
      try {
        const user = await authService.signinCallback();
        setUser(user);
      } catch (e) {
        console.error(e);
      } finally {
        navigate("/");
      }
    })();
  });

  return <h1>Signing in...</h1>;
}
