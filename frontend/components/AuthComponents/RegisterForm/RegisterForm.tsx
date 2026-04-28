"use client";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { MdLogin } from "react-icons/md";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/OtherComponents/Separators/SeparatorWithText/SeparatorWithText";
import ClickableCard from "../ClickableCard/ClickableCard";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";
import "./RegisterForm.css";

interface RegisterFormProps {
        onSubmit: (data: {
                email: string;
                password: string;
                firstName: string;
                lastName: string;
        }) => Promise<void>;
        invitationCode?: string;
        departmentName?: string;
        isSubmitting?: boolean;
}

const RegisterForm: React.FC<RegisterFormProps> = ({
        onSubmit,
        invitationCode,
        departmentName,
        isSubmitting,
}) => {
        const router = useRouter();
        const [email, setEmail] = useState("");
        const [password, setPassword] = useState("");
        const [firstName, setFirstName] = useState("");
        const [lastName, setLastName] = useState("");
        const [showErrors, setShowErrors] = useState(false);

        const isValidEmail = (email: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

        const isFormValid = useMemo(() => {
                return (
                        isValidEmail(email) &&
                        firstName.trim().length >= 3 &&
                        lastName.trim().length >= 3 &&
                        password.length > 0
                );
        }, [email, firstName, lastName, password]);

        const handleSubmit = async () => {
                if (!isFormValid) {
                        setShowErrors(true);
                        return;
                }
                await onSubmit({ email, password, firstName, lastName });
        };

        return (
                <div className="register-form">
                        <LoginFormHeader
                                title="Bun venit pe Portal"
                                subtitle={
                                        departmentName
                                                ? `Te înregistrezi ca membru al departamentului „${departmentName}".`
                                                : "Introduceți datele pentru a vă crea un cont pe Doclane."
                                }
                        />
                        {departmentName && (
                                <div className="register-form-department-banner">
                                        Departament: <strong>{departmentName}</strong>
                                </div>
                        )}
                        <Input
                                label="Email:"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Adresa ta de email"
                        />
                        {showErrors && !isValidEmail(email) && (
                                <p className="register-form-error">Introdu un email valid.</p>
                        )}
                        <Input
                                label="Parolă:"
                                placeholder="Parola ta"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                isPassword
                        />
                        {showErrors && password.length === 0 && (
                                <p className="register-form-error">Parola este obligatorie.</p>
                        )}
                        <Input
                                label="Prenume:"
                                placeholder="Prenumele tău"
                                value={firstName}
                                onChange={(e) => setFirstName(e.target.value)}
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
                                onChange={(e) => setLastName(e.target.value)}
                        />
                        {showErrors && lastName.trim().length < 3 && (
                                <p className="register-form-error">
                                        Numele trebuie să aibă cel puțin 3 caractere.
                                </p>
                        )}
                        <ButtonPrimary
                                text={
                                        isSubmitting
                                                ? "Se creează contul..."
                                                : "Înregistrează-te"
                                }
                                onClick={handleSubmit}
                                disabled={isSubmitting}
                        />
                        <SeparatorWithText text="Ai deja un cont?" />
                        <ClickableCard
                                text="Autentifică-te"
                                icon={<MdLogin size={20} />}
                                onClick={() => router.push("/login")}
                        />
                </div>
        );
};

export default RegisterForm;
