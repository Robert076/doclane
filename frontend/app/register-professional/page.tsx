"use client";
import { useState } from "react";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import RegisterProfessionalForm from "@/components/Forms/RegisterForm/RegisterProfessionalForm";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const router = useRouter();

  const handleRegister = async () => {
    const registerPromise = fetch("/api/backend/auth/register/professional", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email,
        password,
        first_name: firstName,
        last_name: lastName,
      }),
    }).then(async (res) => {
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || "Sign up failed");
      }
      return res.json();
    });

    toast.promise(registerPromise, {
      loading: "Signing up...",
      success: "Sign up successful!",
      error: (err) => `Sign up failed: ${err.message}`,
    });

    registerPromise.then((_) => {
      router.push("/login");
    });
  };

  return (
    <div className="register-page-wrapper">
      <div className="register-page">
        <RegisterProfessionalForm
          email={email}
          setEmail={setEmail}
          password={password}
          setPassword={setPassword}
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
