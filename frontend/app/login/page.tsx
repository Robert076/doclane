"use client";
import { useState } from "react";
import "./style.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import LoginForm from "@/components/AuthComponents/LoginForm/LoginForm";
import { login } from "@/lib/api/auth";

const LoginPage = () => {
        const [email, setEmail] = useState("");
        const [password, setPassword] = useState("");
        const router = useRouter();

        const handleLogin = async () => {
                const loginPromise = login(email, password).then((res) => {
                        if (!res.success) {
                                throw new Error(res.message || "Login failed");
                        }
                        return res;
                });

                toast.promise(loginPromise, {
                        loading: "Logging in...",
                        success: "Login successful!",
                        error: (err) => `Login failed: ${err.message}`,
                });

                await loginPromise;
                router.push("/dashboard/requests");
        };

        return (
                <div className="login-page-wrapper">
                        <div className="login-page">
                                <LoginForm
                                        email={email}
                                        setEmail={setEmail}
                                        password={password}
                                        setPassword={setPassword}
                                        handleLogin={handleLogin}
                                />
                        </div>
                </div>
        );
};

export default LoginPage;
