"use client";
import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import toast from "react-hot-toast";
import { register } from "@/lib/api/auth";
import { getInvitationCodeInfo } from "@/lib/api/invitation_codes";
import RegisterForm from "@/components/AuthComponents/RegisterForm/RegisterForm";
import LoadingSkeleton from "@/components/ViewComponents/LoadingSkeleton/LoadingSkeleton";

export default function InviteRegisterPage() {
        const router = useRouter();
        const searchParams = useSearchParams();
        const code = searchParams.get("code");

        const [departmentName, setDepartmentName] = useState<string | null>(null);
        const [isValidating, setIsValidating] = useState(true);
        const [isSubmitting, setIsSubmitting] = useState(false);

        useEffect(() => {
                if (!code) {
                        router.replace("/register");
                        return;
                }

                getInvitationCodeInfo(code).then((res) => {
                        if (res.success && res.data) {
                                setDepartmentName(res.data.department_name);
                        } else {
                                toast.error(res.message ?? "Cod de invitație invalid.");
                                router.replace("/register");
                        }
                        setIsValidating(false);
                });
        }, [code]);

        const handleSubmit = async (data: {
                email: string;
                password: string;
                firstName: string;
                lastName: string;
        }) => {
                if (!code) return;
                setIsSubmitting(true);
                const response = await register(
                        data.email,
                        data.password,
                        data.firstName,
                        data.lastName,
                        code,
                );
                setIsSubmitting(false);
                if (response.success) {
                        toast.success("Cont creat cu succes!");
                        router.push("/login");
                } else {
                        toast.error(response.message);
                }
        };

        if (isValidating) return <LoadingSkeleton />;

        return (
                <div className="register-page-wrapper">
                        <div className="register-page">
                                <RegisterForm
                                        onSubmit={handleSubmit}
                                        invitationCode={code ?? undefined}
                                        departmentName={departmentName ?? undefined}
                                        isSubmitting={isSubmitting}
                                />
                        </div>
                </div>
        );
}
