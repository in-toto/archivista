import * as React from "react";

import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";
import AuthService, { IAuthService } from "./lib/AuthService";
import { Box, Container, CssBaseline, ThemeProvider, styled } from "@mui/material";
import { Route, Routes } from "react-router-dom";

import { Dashboard } from "./pages/Dashboard";
import DateAdapter from "@mui/lab/AdapterDateFns";
import { FeatureToggleProvider } from "./shared/contexts/FeatureToggleContext/FeatureToggleContext";
import { LocalizationProvider } from "@mui/lab";
import { Navbar } from "./components/Navbar";
import { Registrar } from "./lib/registrar/registrar.pb";
import { RegistrarService } from "./lib/registrar/registrar.pb";
import { Search as SearchPage } from "./pages/Search";
import { SigninCallback } from "./pages/SigninCallback";
import { Sources as SourcesPage } from "./pages/Sources/Sources";
import { Theme } from "./theme";
import { User } from "oidc-client";
import YAML from "yaml";

// import { WitnessApiService } from "./lib/Witness";

const Root = styled(Box)(() => ({
  display: "flex",
  flexFlow: "column",
  height: "100vh",
}));

const Content = styled(Box)(({ theme }) => ({
  background: theme.palette.grey[100],
  flex: "1 1 auto",
}));

export type Config = {
  hydra: {
    url: string;
    client_id: string;
    client_scopes: string;
    root_url: string;
  };
  registrar_svc: string;
  auth_svc: string;
  archivist_svc: string;
};

export const ConfigContext = React.createContext({} as Config);
export const AuthContext = React.createContext({} as IAuthService);
export const RegistrarContext = React.createContext({} as RegistrarService);

const fetchConfig = () => {
  return fetch("/config.yaml")
    .then((response) => response.text())
    .then((text) => YAML.parse(text) as Config);
};

export const App = () => {
  const [graphQLClient, setGraphQLClient] = React.useState<ApolloClient<unknown>>();
  const [config, setConfig] = React.useState(null as Config | null);
  const [authService, setAuthService] = React.useState(null as IAuthService | null);
  const [userName, setUserName] = React.useState<string | null>(null);
  const [user, setUser] = React.useState<User | null>(null);

  React.useEffect(() => {
    const client = new ApolloClient({
      cache: new InMemoryCache(),
      uri: `${config?.archivist_svc || ""}/query`,
    });
    setGraphQLClient(client);
  }, [config]);

  React.useEffect(() => {
    void fetchConfig().then((config) => {
      setConfig(config);
      setAuthService(new AuthService(config));
    });
  }, []);

  React.useEffect(() => {
    if (user) {
      setUserName(user.profile.name ?? null);
      return;
    } else {
      if (authService) {
        void authService.getUser().then((user) => {
          setUser(user);
        });
      }
    }
  }, [user, authService]);

  if (!config || !authService) {
    return null;
  }

  return (
    // note: if graphQlClient fails, nothing will load
    // TODO: implement error state here
    // TODO: if dev kube isn't available, error st
    <>
      {graphQLClient && (
        <FeatureToggleProvider>
          <ApolloProvider client={graphQLClient}>
            <ConfigContext.Provider value={config}>
              <AuthContext.Provider value={authService}>
                <LocalizationProvider dateAdapter={DateAdapter}>
                  <ThemeProvider theme={Theme}>
                    <CssBaseline>
                      <Root>
                        <Navbar onLoginClick={() => void authService.login()} userName={userName} onLogoutClick={() => void authService.logout()} />
                        <Content component="main">
                          <Container sx={{ mt: 2, mb: 2 }}>
                            {user && (
                              <RegistrarContext.Provider value={new Registrar(config.registrar_svc, user.id_token)}>
                                <Routes>
                                  <Route path="/" element={<Dashboard />} />
                                  <Route path="/dashboard" element={<Dashboard />} />
                                  <Route path="/search" element={<SearchPage />} />
                                  <Route path="/sources" element={<SourcesPage />} />
                                  <Route path="/signin-callback" element={<SigninCallback authService={authService} setUser={setUser} />} />
                                </Routes>
                              </RegistrarContext.Provider>
                            )}

                            {!user && (
                              <Routes>
                                <Route path="/" element={<SearchPage />} />
                                <Route path="/signin-callback" element={<SigninCallback authService={authService} setUser={setUser} />} />
                              </Routes>
                            )}
                          </Container>
                        </Content>
                      </Root>
                    </CssBaseline>
                  </ThemeProvider>
                </LocalizationProvider>
              </AuthContext.Provider>
            </ConfigContext.Provider>
          </ApolloProvider>
        </FeatureToggleProvider>
      )}
    </>
  );
};
