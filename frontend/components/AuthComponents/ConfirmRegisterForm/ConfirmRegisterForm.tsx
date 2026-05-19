"use client";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { MdLogin } from "react-icons/md";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/OtherComponents/Separators/SeparatorWithText/SeparatorWithText";
import ClickableCard from "../ClickableCard/ClickableCard";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";
import "./ConfirmRegisterForm.css";

interface ConfirmRegisterFormProps {
        email: string;
        onSubmit: (code: string) => Promise<void>;
        isSubmitting?: boolean;
}

const ConfirmRegisterForm: React.FC<ConfirmRegisterFormProps> = ({
        email,
        onSubmit,
        isSubmitting,
}) => {
        const router = useRouter();
        const [code, setCode] = useState("");
        const [showErrors, setShowErrors] = useState(false);

        const isValidCode = (value: string) => /^\d{4,8}$/.test(value.trim());

        const isFormValid = useMemo(() => isValidCode(code), [code]);

        const handleSubmit = async () => {
                if (!isFormValid) {
                        setShowErrors(true);
                        return;
                }
                await onSubmit(code.trim());
        };

        return (
                <div className="confirm-register-form">
                        <LoginFormHeader
                                title="Confirmă-ți contul"
                                subtitle={`Am trimis un cod de verificare la adresa ${email}. Introdu-l mai jos pentru a finaliza înregistrarea.`}
                        />
                        <Input
                                label="Cod de verificare:"
                                value={code}
                                onChange={(e) => setCode(e.target.value)}
                                placeholder="Codul primit pe email"
                        />
                        {showErrors && !isValidCode(code) && (
                                <p className="confirm-register-form-error">
                                        Introdu codul primit pe email (doar cifre).
                                </p>
                        )}
                        <ButtonPrimary
                                text={
                                        isSubmitting
                                                ? "Se confirmă contul..."
                                                : "Confirmă contul"
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

export default ConfirmRegisterForm;
