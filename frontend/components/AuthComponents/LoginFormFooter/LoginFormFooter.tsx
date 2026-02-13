import { MdLock } from "react-icons/md";
import "./LoginFormFooter.css";

const LoginFormFooter = () => {
        return (
                <div className="login-form-footer">
                        <MdLock size={20} />
                        <p>Securely encrypted for your privacy.</p>
                </div>
        );
};

export default LoginFormFooter;
