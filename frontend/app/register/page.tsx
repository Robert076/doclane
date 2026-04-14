"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { register } from "@/lib/api/auth";
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
                const response = await register(
                        data.email,
                        data.password,
                        data.firstName,
                        data.lastName,
                );
                setIsSubmitting(false);
                if (response.success) {
                        toast.success("Cont creat cu succes!");
                        router.push("/login");
                } else {
                        toast.error(response.message);
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
