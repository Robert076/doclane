import "./Input.css";

interface InputProps {
        value: string;
        onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
        label?: string;
        isPassword?: boolean;
        icon?: React.ReactNode;
        placeholder?: string;
        fullWidth?: boolean;
}

const Input: React.FC<InputProps> = ({
        value,
        onChange,
        label,
        isPassword,
        icon,
        placeholder,
        fullWidth,
}) => {
        return (
                <div
                        className={`input-wrapper ${fullWidth ? "input-wrapper--full-width" : ""}`}
                >
                        {label && <label>{label}</label>}
                        <div className="input-with-icon">
                                {icon && <span className="input-icon">{icon}</span>}
                                <input
                                        type={isPassword ? "password" : "text"}
                                        value={value}
                                        onChange={onChange}
                                        placeholder={placeholder}
                                />
                        </div>
                </div>
        );
};

export default Input;
