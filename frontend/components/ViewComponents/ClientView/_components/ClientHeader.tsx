import { User } from "@/types";
import React from "react";
import "./ClientHeader.css";

interface ClientHeaderProps {
  user: User;
  length: number;
}

const ClientHeader: React.FC<ClientHeaderProps> = ({ user, length }) => {
  return (
    <header className="client-header">
      <h1 className="overview-h1">
        Welcome back, {user.first_name} {user.last_name}
      </h1>
      <p className="overview-p">You have {length} active document requests.</p>
    </header>
  );
};

export default ClientHeader;
