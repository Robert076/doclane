"use client";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";
import { useUser } from "@/context/UserContext";
import "./TemplatesActions.css";

export default function TemplatesActions() {
        const router = useRouter();
        const user = useUser();
        const canManage = user.role === "admin" || user.department_id !== null;

        if (!canManage) return null;

        return (
                <div className="templates-actions has-margin-bottom">
                        <ButtonPrimary
                                text="Șablon nou"
                                onClick={() => router.push("/dashboard/templates/create")}
                        />
                </div>
        );
}
