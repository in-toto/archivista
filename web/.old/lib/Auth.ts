export interface Token {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  scope: string;
  created_at: number;
}

export async function getToken(): Promise<Token | null> {
  let token = window.sessionStorage.getItem("token");
  if (token != null && token != "") {
    return JSON.parse(token) as Token;
  }

  const code = new URLSearchParams(window.location.search).get("code");
  if (code != null) {
    token = await requestToken(code);
    window.sessionStorage.setItem("token", token);
    return JSON.parse(token) as Token;
  }

  return null;
}

export async function requestToken(code: string): Promise<string> {
  const response = await fetch(window.OAUTH_TOKEN_ENDPOINT ?? "https://gitlab.com/oauth/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: window.OAUTH_APP_ID ?? "",
      code_verifier: window.sessionStorage.getItem("verifier") ?? "",
      grant_type: "authorization_code",
      redirect_uri: window.OAUTH_REDIRECT_URL ?? "http://localhost:1234",
      code: code,
    }),
  });

  if (response.status > 300) {
    return "";
  }
  return await response.text();
}

export function generateRandomString(length: number): string {
  let text = "";
  const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

  for (let i = 0; i < length; i++) {
    text += possible.charAt(Math.floor(Math.random() * possible.length));
  }

  return text;
}

export async function generateCodeChallenge(verifier: string): Promise<string> {
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(verifier));
  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/=/g, "")
    .replace(/\+/g, "-")
    .replace(/\//g, "_");
}
