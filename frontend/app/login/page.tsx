"use client";
import { useState } from "react";
import LoginForm from "../components/LoginForm/LoginForm";
import "./style.css";
import toast from "react-hot-toast";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    const loginPromise = fetch("http://localhost:8080/api/auth/login", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    }).then(async (res) => {
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || "Login failed");
      }
      return res.json();
    });

    toast.promise(loginPromise, {
      loading: "Logging in...",
      success: "Login successful!",
      error: (err) => `Login failed: ${err.message}`,
    });

    loginPromise.then((data) => {
      console.log("User logged in:", data);
    });
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
