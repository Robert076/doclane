import React from "react";
import "./ButtonPrimary.css";
import { IconType } from "react-icons";

interface ButtonPrimaryProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  text: string;
  variant?: "primary" | "secondary" | "danger" | "ghost";
  fullWidth?: boolean;
  icon?: IconType;
  iconPosition?: "left" | "right";
}

const ButtonPrimary: React.FC<ButtonPrimaryProps> = ({
  text,
  variant = "primary",
  fullWidth = false,
  icon: Icon,
  iconPosition = "left",
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
      {Icon && iconPosition === "left" && <Icon className="button-icon" />}
      {text}
      {Icon && iconPosition === "right" && <Icon className="button-icon" />}
    </button>
  );
};

export default ButtonPrimary;
