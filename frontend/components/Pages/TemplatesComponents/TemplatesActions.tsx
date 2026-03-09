"use client";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";
import "./TemplatesActions.css";

const TemplatesActions = () => {
        const router = useRouter();
        return (
                <div className="templates-actions has-margin-bottom">
                        <ButtonPrimary
                                text="Şablon nou"
                                onClick={() => {
                                        router.push("/dashboard/templates/create");
                                }}
                        />
                </div>
        );
};

export default TemplatesActions;
