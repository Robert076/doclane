"use client";
import { useRouter, useSearchParams } from "next/navigation";
import { useUser } from "@/context/UserContext";
import { Department } from "@/types";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Select from "@/components/InputComponents/Select";
import "./TemplatesActions.css";

interface Props {
        departments?: Department[];
}

export default function TemplatesActions({ departments = [] }: Props) {
        const router = useRouter();
        const searchParams = useSearchParams();
        const user = useUser();

        const selectedDepartmentId = searchParams.get("department") ?? "";

        const handleDepartmentChange = (val: string) => {
                const params = new URLSearchParams(searchParams.toString());
                if (val) {
                        params.set("department", val);
                } else {
                        params.delete("department");
                }
                router.push(`?${params.toString()}`);
        };

        if (user.role !== "admin") return null;

        return (
                <div className="templates-actions has-margin-bottom">
                        <div className="templates-actions-button">
                                <ButtonPrimary
                                        text="Șablon nou"
                                        onClick={() =>
                                                router.push("/dashboard/templates/create")
                                        }
                                />
                        </div>
                        <Select
                                value={selectedDepartmentId}
                                onChange={handleDepartmentChange}
                                placeholder="Toate departamentele"
                                options={departments.map((d) => ({
                                        value: d.id,
                                        label: d.name,
                                }))}
                        />
                </div>
        );
}
