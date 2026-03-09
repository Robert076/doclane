import React from "react";
import "./StatusBadge.css";
import { RequestStatus } from "@/types";

interface StatusBadgeProps {
        status: RequestStatus;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status }) => {
        return (
                <span className={`status-badge ${status.toLowerCase()}`}>
                        {status.toUpperCase()}
                </span>
        );
};

export default StatusBadge;
