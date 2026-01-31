import React from "react";
import "./ButtonPrimary.css";

interface ButtonPrimaryProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  text: string;
  variant?: "primary" | "secondary" | "danger" | "ghost";
  fullWidth?: boolean;
}

const ButtonPrimary: React.FC<ButtonPrimaryProps> = ({
  text,
  variant = "primary",
  fullWidth = false,
  className = "",
  ...props
}) => {
  const variantClass = `button-primary--${variant}`;
  const widthClass = fullWidth ? "button-primary--full-width" : "";

  return (
    <button
      className={`button button-primary-base ${variantClass} ${widthClass} ${className}`}
      {...props}
    >
      {text}
    </button>
  );
};

export default ButtonPrimary;
