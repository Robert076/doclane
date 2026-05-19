"use client";
import { Suspense, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { confirmRegistration } from "@/lib/client/auth";
import ConfirmRegisterForm from "@/components/AuthComponents/ConfirmRegisterForm/ConfirmRegisterForm";
import LoadingSkeleton from "@/components/ViewComponents/LoadingSkeleton/LoadingSkeleton";

const PENDING_KEY = "pendingRegistration";

interface PendingRegistration {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  invitationCode?: string;
}

function ConfirmRegisterPage() {
  const router = useRouter();
  const [pending, setPending] = useState<PendingRegistration | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    try {
      const raw = sessionStorage.getItem(PENDING_KEY);
      if (!raw) {
        router.replace("/register");
        return;
      }
      const parsed = JSON.parse(raw) as PendingRegistration;
      if (!parsed?.email || !parsed?.password) {
        router.replace("/register");
        return;
      }
      setPending(parsed);
    } catch {
      router.replace("/register");
      return;
    }
    setIsLoading(false);
  }, [router]);

  const handleSubmit = async (code: string) => {
    if (!pending) return;
    setIsSubmitting(true);
    try {
      // Step 1 — confirm with Cognito, sign in, set cookie, return token
      const syncParams = await confirmRegistration(
        pending.email,
        pending.password,
        pending.firstName,
        pending.lastName,
        code,
        pending.invitationCode,
      );

      // Step 2 — call backend directly with the token since the cookie
      // isn't visible to Next.js server actions in the same cycle
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/auth/sync`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${syncParams.idToken}`,
          },
          body: JSON.stringify({
            email: syncParams.email,
            first_name: syncParams.firstName,
            last_name: syncParams.lastName,
            ...(syncParams.invitationCode
              ? { invitation_code: syncParams.invitationCode }
              : {}),
          }),
        },
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error ?? "Sincronizarea contului a eșuat.");
      }

      sessionStorage.removeItem(PENDING_KEY);
      toast.success("Cont confirmat cu succes!");
      router.push("/dashboard/requests");
    } catch (err: any) {
      toast.error(err?.message ?? "Confirmarea a eșuat.");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading || !pending) return <LoadingSkeleton />;

  return (
    <div className="register-page-wrapper">
      <div className="register-page">
        <ConfirmRegisterForm
          email={pending.email}
          onSubmit={handleSubmit}
          isSubmitting={isSubmitting}
        />
      </div>
    </div>
  );
}

export default function ConfirmRegisterPageWrapper() {
  return (
    <Suspense fallback={<LoadingSkeleton />}>
      <ConfirmRegisterPage />
    </Suspense>
  );
}