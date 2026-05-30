"use client";
import { useState } from "react";
import "./style.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import LoginForm from "@/components/AuthComponents/LoginForm/LoginForm";
import { login } from "@/lib/client/auth";

const LoginPage = () => {
        const [email, setEmail] = useState("");
        const [password, setPassword] = useState("");
        const router = useRouter();

        const handleLogin = async () => {
    const loginPromise = login(email, password);

    toast.promise(loginPromise, {
        loading: "Logging in...",
        success: "Login successful!",
        error: (err) => `Login failed: ${err.message}`,
    });

    await loginPromise;
    window.location.href = "/dashboard/requests";
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
