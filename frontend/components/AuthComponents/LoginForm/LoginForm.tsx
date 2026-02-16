import { Dispatch, SetStateAction } from "react";
import "./LoginForm.css";

import ClickableCard from "../ClickableCard/ClickableCard";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/OtherComponents/Separators/SeparatorWithText/SeparatorWithText";
import LoginFormFooter from "../LoginFormFooter/LoginFormFooter";
import LoginFormHeader from "../LoginFormHeader/LoginFormHeader";
import { MdCardGiftcard, MdWork } from "react-icons/md";
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
                                title="Welcome to your Portal"
                                subtitle="Please enter your details to access your secure documents and
                                requests."
                        />

                        <Input
                                label="Email:"
                                value={email}
                                onChange={(e: any) => setEmail(e.target.value)}
                                placeholder="Your email address here"
                        />

                        <Input
                                label="Password:"
                                placeholder="Your password here"
                                value={password}
                                onChange={(e: any) => setPassword(e.target.value)}
                                isPassword={true}
                        />

                        <ButtonPrimary text="Log In" onClick={handleLogin} />

                        <SeparatorWithText text="New to Doclane?" />

                        <ClickableCard
                                text="Redeem Invitation"
                                icon={<MdCardGiftcard size={20} />}
                                onClick={() => {
                                        router.push("/register-client");
                                }}
                        />

                        <ClickableCard
                                text="Join as a Professional"
                                icon={<MdWork size={20} />}
                                onClick={() => {
                                        router.push("/register-professional");
                                }}
                        />

                        <LoginFormFooter />
                </div>
        );
};

export default LoginForm;
