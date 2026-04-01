"use client";
import { useState } from "react";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import { register } from "@/lib/api/auth";
import RegisterClientForm from "@/components/AuthComponents/RegisterForm/RegisterForm";

export default function RegisterPage() {
        const [email, setEmail] = useState("");
        const [password, setPassword] = useState("");
        const [invitationCode, setInvitationCode] = useState("");
        const [firstName, setFirstName] = useState("");
        const [lastName, setLastName] = useState("");
        const router = useRouter();

        const handleRegister = async () => {
                const loadingToastID = toast.loading("Se creează contul...");
                const response = await register(
                        email,
                        password,
                        invitationCode,
                        firstName,
                        lastName,
                );
                toast.dismiss(loadingToastID);

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
                                <RegisterClientForm
                                        email={email}
                                        setEmail={setEmail}
                                        password={password}
                                        setPassword={setPassword}
                                        invitationCode={invitationCode}
                                        setInvitationCode={setInvitationCode}
                                        firstName={firstName}
                                        setFirstName={setFirstName}
                                        lastName={lastName}
                                        setLastName={setLastName}
                                        handleRegister={handleRegister}
                                />
                        </div>
                </div>
        );
}
