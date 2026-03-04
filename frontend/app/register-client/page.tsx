"use client";
import { useState } from "react";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import RegisterClientForm from "@/components/AuthComponents/RegisterForm/RegisterClientForm";
import { signUpClient } from "@/lib/api/api";

const LoginPage = () => {
        const [email, setEmail] = useState("");
        const [password, setPassword] = useState("");
        const [invitationCode, setInvitationCode] = useState("");
        const [firstName, setFirstName] = useState("");
        const [lastName, setLastName] = useState("");
        const router = useRouter();

        const handleRegister = async () => {
                const loadingToastID = toast.loading("Signing up...");

                const response = await signUpClient(
                        email,
                        password,
                        invitationCode,
                        firstName,
                        lastName,
                );

                toast.dismiss(loadingToastID);
                if (response.success) {
                        toast.success("Signed up successfully!");
                        router.push("/dashboard");
                } else {
                        toast.error("Failed to sign up: " + response.message);
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
};

export default LoginPage;
