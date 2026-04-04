"use client";
import { Department } from "@/types";
import BaseDashboardCard from "@/components/CardComponents/BaseDashboardCard/BaseDashboardCard";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { formatDate } from "@/lib/client/formatDate";
import { useRouter } from "next/navigation";

interface DepartmentCardProps {
        department: Department;
}

export default function DepartmentCard({ department }: DepartmentCardProps) {
        const router = useRouter();

        return (
                <BaseDashboardCard
                        title={department.name}
                        footer={
                                <>
                                        <ButtonPrimary
                                                text="Vezi membri"
                                                variant="ghost"
                                                fullWidth
                                                onClick={() =>
                                                        router.push(
                                                                `/dashboard/departments/${department.id}`,
                                                        )
                                                }
                                        />
                                        <ButtonPrimary
                                                text="Vezi șabloane"
                                                variant="ghost"
                                                fullWidth
                                                onClick={() =>
                                                        router.push(
                                                                `/dashboard/templates?department=${department.id}`,
                                                        )
                                                }
                                        />
                                </>
                        }
                >
                        <InfoList>
                                <InfoItem
                                        label="Creat la"
                                        value={formatDate(department.created_at)}
                                />
                        </InfoList>
                </BaseDashboardCard>
        );
}
