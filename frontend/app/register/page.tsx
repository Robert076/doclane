"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { register } from "@/lib/client/auth";
import RegisterForm from "@/components/AuthComponents/RegisterForm/RegisterForm";

export default function RegisterPage() {
        const router = useRouter();
        const [isSubmitting, setIsSubmitting] = useState(false);

        const handleSubmit = async (data: {
                email: string;
                password: string;
                firstName: string;
                lastName: string;
        }) => {
                setIsSubmitting(true);
                try {
                        const pending = await register(
                                data.email,
                                data.password,
                                data.firstName,
                                data.lastName,
                        );
                        sessionStorage.setItem(
                                "pendingRegistration",
                                JSON.stringify(pending),
                        );
                        toast.success("Verifică-ți emailul pentru codul de confirmare.");
                        router.push("/register/confirm");
                } catch (err: any) {
                        toast.error(err?.message ?? "Înregistrarea a eșuat.");
                } finally {
                        setIsSubmitting(false);
                }
        };

        return (
                <div className="register-page-wrapper">
                        <div className="register-page">
                                <RegisterForm
                                        onSubmit={handleSubmit}
                                        isSubmitting={isSubmitting}
                                />
                        </div>
                </div>
        );
}
