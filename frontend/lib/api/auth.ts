"use server";

import { cookies } from "next/headers";
import { APIResponse } from "@/types";
import { doclaneHTTPHelper } from "./core";
import { CognitoIdentityProviderClient, ChangePasswordCommand } from "@aws-sdk/client-cognito-identity-provider";

// Called by the client after Amplify confirmSignUp + signIn succeeds
export async function syncUser(
  firstName: string,
  lastName: string,
  invitationCode?: string
): Promise<APIResponse> {
  return doclaneHTTPHelper("/auth/sync", {
    method: "POST",
    body: {
      first_name: firstName,
      last_name: lastName,
      ...(invitationCode ? { invitation_code: invitationCode } : {}),
    },
  });
}

export async function changePassword(
  oldPassword: string,
  newPassword: string
): Promise<APIResponse> {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get("access_token")?.value;
  if (!accessToken) return { success: false, message: "Not authenticated." };

  try {
    const client = new CognitoIdentityProviderClient({ 
      region: process.env.NEXT_PUBLIC_AWS_REGION 
    });
    await client.send(new ChangePasswordCommand({
      AccessToken: accessToken,
      PreviousPassword: oldPassword,
      ProposedPassword: newPassword,
    }));
    return { success: true, message: "Password updated successfully." };
  } catch (error: any) {
    return { success: false, message: error.message || "Failed to update password." };
  }
}