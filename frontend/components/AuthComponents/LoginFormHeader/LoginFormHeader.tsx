import Logo from "@/components/Logo/Logo";
import Separator from "@/components/Separators/Separator/Separator";
import "./LoginFormHeader.css";

interface LoginFormHeaderProps {
        title: string;
        subtitle: string;
}

const LoginFormHeader: React.FC<LoginFormHeaderProps> = ({ title, subtitle }) => {
        return (
                <div className="login-form-wrapper">
                        <Logo />
                        <h2 className="login-form-title">{title}</h2>
                        <p className="login-form-subtitle">{subtitle}</p>
                        <Separator />
                </div>
        );
};

export default LoginFormHeader;
