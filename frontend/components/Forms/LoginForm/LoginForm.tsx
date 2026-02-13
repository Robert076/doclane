import { Dispatch, SetStateAction } from "react";
import "./LoginForm.css";

import ClickableCard from "../ClickableCard/ClickableCard";
import { MdCardGiftcard, MdLock, MdWork } from "react-icons/md";
import Separator from "@/components/Separators/Separator/Separator";
import Logo from "@/components/Logo/Logo";
import Input from "@/components/Input/Input";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import SeparatorWithText from "@/components/Separators/SeparatorWithText/SeparatorWithText";
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
      <Logo />
      <h2 className="login-form-title">Welcome to your Portal</h2>
      <p className="login-form-subtitle">
        Please enter your details to access your secure documents and requests.
      </p>

      <Separator />

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

      <div className="login-form-footer">
        <MdLock size={20} />
        <p>Securely encrypted for your privacy.</p>
      </div>
    </div>
  );
};

export default LoginForm;
