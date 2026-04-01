"use client";
import { Dispatch, SetStateAction, useMemo, useState } from "react";
import "./RegisterForm.css";
import { MdLogin } from "react-icons/md";
import Logo from "@/components/OtherComponents/Logo/Logo";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/OtherComponents/Separators/SeparatorWithText/SeparatorWithText";
import ClickableCard from "../ClickableCard/ClickableCard";
import { useRouter } from "next/navigation";
import LoginFormFooter from "../LoginFormFooter/LoginFormFooter";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";

interface RegisterFormProps {
        email: string;
        setEmail: Dispatch<SetStateAction<string>>;
        password: string;
        setPassword: Dispatch<SetStateAction<string>>;
        invitationCode: string;
        setInvitationCode: Dispatch<SetStateAction<string>>;
        firstName: string;
        setFirstName: Dispatch<SetStateAction<string>>;
        lastName: string;
        setLastName: Dispatch<SetStateAction<string>>;
        handleRegister: () => void;
}

const RegisterForm: React.FC<RegisterFormProps> = ({
        email,
        setEmail,
        password,
        setPassword,
        invitationCode,
        setInvitationCode,
        firstName,
        setFirstName,
        lastName,
        setLastName,
        handleRegister,
}) => {
        const router = useRouter();
        const [showErrors, setShowErrors] = useState(false);

        const isValidEmail = (email: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

        const isFormValid = useMemo(() => {
                return (
                        isValidEmail(email) &&
                        firstName.trim().length >= 3 &&
                        lastName.trim().length >= 3 &&
                        password.length > 0 &&
                        invitationCode.trim().length > 0
                );
        }, [email, firstName, lastName, password, invitationCode]);

        const handleSubmit = () => {
                if (!isFormValid) {
                        setShowErrors(true);
                        return;
                }
                handleRegister();
        };

        return (
                <div className="register-form">
                        <LoginFormHeader
                                title="Bun venit pe Portal"
                                subtitle="Introduceți datele pentru a vă crea un cont pe Doclane."
                        />
                        <Input
                                label="Email:"
                                value={email}
                                onChange={(e: any) => setEmail(e.target.value)}
                                placeholder="Adresa ta de email"
                        />
                        {showErrors && !isValidEmail(email) && (
                                <p className="register-form-error">Introdu un email valid.</p>
                        )}
                        <Input
                                label="Parolă:"
                                placeholder="Parola ta"
                                value={password}
                                onChange={(e: any) => setPassword(e.target.value)}
                                isPassword={true}
                        />
                        {showErrors && password.length === 0 && (
                                <p className="register-form-error">Parola este obligatorie.</p>
                        )}
                        <Input
                                label="Prenume:"
                                placeholder="Prenumele tău"
                                value={firstName}
                                onChange={(e: any) => setFirstName(e.target.value)}
                        />
                        {showErrors && firstName.trim().length < 3 && (
                                <p className="register-form-error">
                                        Prenumele trebuie să aibă cel puțin 3 caractere.
                                </p>
                        )}
                        <Input
                                label="Nume:"
                                placeholder="Numele tău"
                                value={lastName}
                                onChange={(e: any) => setLastName(e.target.value)}
                        />
                        {showErrors && lastName.trim().length < 3 && (
                                <p className="register-form-error">
                                        Numele trebuie să aibă cel puțin 3 caractere.
                                </p>
                        )}
                        <Input
                                label="Cod de invitație:"
                                placeholder="Codul tău de invitație"
                                value={invitationCode}
                                onChange={(e: any) => setInvitationCode(e.target.value)}
                        />
                        {showErrors && invitationCode.trim().length === 0 && (
                                <p className="register-form-error">
                                        Codul de invitație este obligatoriu.
                                </p>
                        )}
                        <ButtonPrimary text="Înregistrează-te" onClick={handleSubmit} />
                        <SeparatorWithText text="Ai deja un cont?" />
                        <ClickableCard
                                text="Autentifică-te"
                                icon={<MdLogin size={20} />}
                                onClick={() => router.push("/login")}
                        />
                        <LoginFormFooter />
                </div>
        );
};

export default RegisterForm;
