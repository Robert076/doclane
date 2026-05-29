import "@/lib/amplify";
import { signIn, signUp, confirmSignUp, signOut, fetchAuthSession } from "@aws-amplify/auth";
import { logout as serverLogout } from "@/lib/api/auth";
import { syncUser } from "@/lib/api/auth";

async function setAuthCookieViaAPI(idToken: string, accessToken: string) {
  const res = await fetch("/api/auth/set-cookie", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ idToken, accessToken }),
    credentials: "same-origin",
  });
  if (!res.ok) throw new Error("Failed to set auth cookie");
}

export async function login(email: string, password: string) {
  // Defensive: ensure no lingering Amplify session, otherwise signIn throws
  // UserAlreadyAuthenticatedException.
  try { await signOut(); } catch (_) {}

  const result = await signIn({ username: email, password });
  if (!result.isSignedIn) throw new Error("Login failed.");

  const session = await fetchAuthSession();
  const idToken = session.tokens?.idToken?.toString();
  const accessToken = session.tokens?.accessToken?.toString();
  if (!idToken || !accessToken) throw new Error("Could not retrieve session token.");

  await setAuthCookieViaAPI(idToken, accessToken);
}

export async function register(
  email: string,
  password: string,
  firstName: string,
  lastName: string,
  invitationCode?: string
) {
  const { nextStep } = await signUp({
    username: email,
    password,
    options: { userAttributes: { email } },
  });

  if (nextStep.signUpStep !== "CONFIRM_SIGN_UP") {
    throw new Error("Unexpected registration step.");
  }

  // Return what's needed for the confirmation step
  return { email, password, firstName, lastName, invitationCode };
}
// lib/client/auth.ts
export async function confirmRegistration(
  email: string,
  password: string,
  firstName: string,
  lastName: string,
  code: string,
  invitationCode?: string
) {
  const { isSignUpComplete } = await confirmSignUp({
    username: email,
    confirmationCode: code,
  });
  if (!isSignUpComplete) throw new Error("Confirmation failed.");

  try { await signOut({ global: false }); } catch (_) {}
  await new Promise(resolve => setTimeout(resolve, 100));
  await signIn({ username: email, password });

  const session = await fetchAuthSession();
  const idToken = session.tokens?.idToken?.toString();
  const accessToken = session.tokens?.accessToken?.toString();
  if (!idToken || !accessToken) throw new Error("Could not retrieve session token.");

  await setAuthCookieViaAPI(idToken, accessToken);
  return { email, firstName, lastName, invitationCode, idToken };
}

export async function logout() {
  try { await signOut(); } catch (_) {}
  await serverLogout();
}