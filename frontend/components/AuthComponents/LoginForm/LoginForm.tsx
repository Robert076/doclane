"use client";
import { Dispatch, SetStateAction } from "react";
import "./LoginForm.css";
import ClickableCard from "../ClickableCard/ClickableCard";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/OtherComponents/Separators/SeparatorWithText/SeparatorWithText";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";
import { MdCardGiftcard } from "react-icons/md";
import { useRouter } from "next/navigation";

interface LoginFormProps {
        email: string;
        setEmail: Dispatch<SetStateAction<string>>;
        password: string;
        setPassword: Dispatch<SetStateAction<string>>;
        handleLogin: () => void;
}

const LoginForm: React.FC<LoginFormProps> = ({
        email,
        setEmail,
        password,
        setPassword,
        handleLogin,
}) => {
        const router = useRouter();

        return (
                <div className="login-form">
                        <LoginFormHeader
                                title="Bun venit pe Portal"
                                subtitle="Introduceți datele pentru a accesa documentele și cererile dvs."
                        />
                        <Input
                                label="Email:"
                                value={email}
                                onChange={(e: any) => setEmail(e.target.value)}
                                placeholder="Adresa ta de email"
                        />
                        <Input
                                label="Parolă:"
                                placeholder="Parola ta"
                                value={password}
                                onChange={(e: any) => setPassword(e.target.value)}
                                isPassword={true}
                        />
                        <ButtonPrimary text="Autentificare" onClick={handleLogin} />
                        <SeparatorWithText text="Nou pe Doclane?" />
                        <ClickableCard
                                text="Înregistrează-te"
                                icon={<MdCardGiftcard size={20} />}
                                onClick={() => router.push("/register")}
                        />
                </div>
        );
};

export default LoginForm;
