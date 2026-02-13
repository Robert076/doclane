"use client";
import { Dispatch, SetStateAction, useMemo, useState } from "react";
import "./RegisterForm.css";
import { MdLogin, MdLock } from "react-icons/md";
import Separator from "@/components/Separators/Separator/Separator";
import Logo from "@/components/Logo/Logo";
import Input from "@/components/Input/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/Separators/SeparatorWithText/SeparatorWithText";
import ClickableCard from "../ClickableCard/ClickableCard";
import { useRouter } from "next/navigation";
import LoginFormFooter from "../LoginFormFooter/LoginFormFooter";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";

interface RegisterClientFormProps {
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

const RegisterClientForm: React.FC<RegisterClientFormProps> = ({
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

        const isValidEmail = (email: string) => {
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                return emailRegex.test(email);
        };

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
                                title="Welcome to your Portal"
                                subtitle="Please enter your details to sign up to Doclane."
                        />
                        <Input
                                label="Email:"
                                value={email}
                                onChange={(e: any) => setEmail(e.target.value)}
                                placeholder="Your email address here"
                        />
                        {showErrors && !isValidEmail(email) && (
                                <p className="register-form-error">
                                        Please enter a valid email
                                </p>
                        )}
                        <Input
                                label="Password:"
                                placeholder="Your password here"
                                value={password}
                                onChange={(e: any) => setPassword(e.target.value)}
                                isPassword={true}
                        />
                        {showErrors && password.length === 0 && (
                                <p className="register-form-error">Password is required</p>
                        )}
                        <Input
                                label="First name:"
                                placeholder="Your first name here"
                                value={firstName}
                                onChange={(e: any) => setFirstName(e.target.value)}
                        />
                        {showErrors && firstName.trim().length < 3 && (
                                <p className="register-form-error">
                                        First name must be at least 3 characters
                                </p>
                        )}
                        <Input
                                label="Last name:"
                                placeholder="Your last name here"
                                value={lastName}
                                onChange={(e: any) => setLastName(e.target.value)}
                        />
                        {showErrors && lastName.trim().length < 3 && (
                                <p className="register-form-error">
                                        Last name must be at least 3 characters
                                </p>
                        )}
                        <Input
                                label="Invitation code:"
                                placeholder="Your invitation code here"
                                value={invitationCode}
                                onChange={(e: any) => setInvitationCode(e.target.value)}
                        />
                        {showErrors && invitationCode.trim().length === 0 && (
                                <p className="register-form-error">
                                        Invitation code is required
                                </p>
                        )}
                        <ButtonPrimary text="Sign up" onClick={handleSubmit} />
                        <SeparatorWithText text="Already have an account?" />
                        <ClickableCard
                                text="Log in"
                                icon={<MdLogin size={20} />}
                                onClick={() => {
                                        router.push("/login");
                                }}
                        />
                        <LoginFormFooter />
                </div>
        );
};

export default RegisterClientForm;
