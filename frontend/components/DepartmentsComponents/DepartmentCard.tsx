"use client";
import { Department } from "@/types";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import { formatDate } from "@/lib/client/formatDate";

interface DepartmentCardProps {
        department: Department;
}

export default function DepartmentCard({ department }: DepartmentCardProps) {
        return (
                <BaseDashboardCard title={department.name}>
                        <div className="template-info">
                                <div className="template-info-item">
                                        <span className="template-label">Creat la:</span>
                                        <span className="template-value">
                                                {formatDate(department.created_at)}
                                        </span>
                                </div>
                        </div>
                </BaseDashboardCard>
        );
}
