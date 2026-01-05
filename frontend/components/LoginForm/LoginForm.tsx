import { Dispatch, SetStateAction } from "react";
import Input from "../Input/Input";
import Logo from "../Logo/Logo";
import "./LoginForm.css";

import SeparatorWithText from "../Separators/SeparatorWithText/SeparatorWithText";
import Separator from "../Separators/Separator/Separator";
import ClickableCard from "./ClickableCard/ClickableCard";
import { MdCardGiftcard, MdLock, MdWork } from "react-icons/md";
import ButtonPrimary from "../Buttons/ButtonPrimary/ButtonPrimary";

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
        onClick={() => console.log("Clicked!")}
      />

      <ClickableCard
        text="Join as a Professional"
        icon={<MdWork size={20} />}
        onClick={() => console.log("Sign up as Professional clicked")}
      />

      <div className="login-form-footer">
        <MdLock size={20} />
        <p>Securely encrypted for your privacy.</p>
      </div>
    </div>
  );
};

export default LoginForm;
