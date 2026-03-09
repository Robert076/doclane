"use client";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import { useState } from "react";
import "./CreateTemplateForm.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import CronInput from "@/components/Pages/ClientsComponents/CronInput";
import { RecurrenceUnit } from "@/types";
import RadioInput from "@/components/InputComponents/RadioInput";
import ExpectedDocumentsList, {
        ExpectedDocumentInput,
} from "@/components/Pages/ClientsComponents/ExpectedDocumentsList";
import { createTemplate, addExpectedDocumentTemplate } from "@/lib/api/api";

const CreateTemplateForm = () => {
        const [title, setTitle] = useState("");
        const [description, setDescription] = useState("");
        const [isNoneSelected, setIsNoneSelected] = useState(true);
        const [isRecurring, setIsRecurring] = useState(false);
        const [unit, setUnit] = useState<RecurrenceUnit>("month");
        const [hour, setHour] = useState("09");
        const [minute, setMinute] = useState("00");
        const [expectedDocuments, setExpectedDocuments] = useState<ExpectedDocumentInput[]>([
                { title: "", description: "" },
        ]);
        const router = useRouter();

        const validateForm = () => {
                if (!title) {
                        toast.error("Completează titlul şablonului.");
                        return false;
                }
                if (expectedDocuments.some((doc) => !doc.title)) {
                        toast.error("Completează titlul fiecărui document solicitat.");
                        return false;
                }
                return true;
        };

        const handleSubmit = async (e: React.FormEvent) => {
                e.preventDefault();
                if (!validateForm()) return;

                toast.promise(
                        (async () => {
                                console.log(title);
                                const res = await createTemplate({
                                        title,
                                        description: description || undefined,
                                        is_recurring: isRecurring,
                                        recurrence_cron: isRecurring
                                                ? buildCron(hour, minute, unit)
                                                : undefined,
                                });

                                if (!res.success || !res.data) throw new Error(res.error);

                                const templateID = res.data;

                                for (const ed of expectedDocuments) {
                                        const addRes = await addExpectedDocumentTemplate(
                                                templateID,
                                                ed.title,
                                                ed.description,
                                                ed.exampleFile,
                                        );
                                        if (!addRes.success) throw new Error(addRes.error);
                                }

                                router.push("/dashboard/templates");
                        })(),
                        {
                                loading: "Se creează şablonul...",
                                success: "Şablon creat cu succes!",
                                error: (err) => `Eroare: ${err.message}`,
                        },
                );
        };

        return (
                <>
                        <form className="add-form" onSubmit={handleSubmit}>
                                <Input
                                        label="Titlul şablonului"
                                        placeholder="Scrie titlul şablonului..."
                                        value={title}
                                        onChange={(e: any) => setTitle(e.target.value)}
                                />
                                <TextArea
                                        label="Descrierea şablonului"
                                        placeholder="Scrie descrierea şablonului..."
                                        value={description}
                                        onChange={(e: any) => setDescription(e.target.value)}
                                />
                                <div className="radio-inputs-time">
                                        <RadioInput
                                                isChecked={isNoneSelected}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) => {
                                                        setIsNoneSelected(e.target.checked);
                                                        setIsRecurring(false);
                                                }}
                                                label="Fără recurenţă"
                                        />
                                        <RadioInput
                                                isChecked={isRecurring}
                                                onChange={(
                                                        e: React.ChangeEvent<HTMLInputElement>,
                                                ) => {
                                                        setIsRecurring(e.target.checked);
                                                        setIsNoneSelected(false);
                                                }}
                                                label="Recurent"
                                        />
                                </div>
                                {isRecurring && (
                                        <CronInput
                                                unit={unit}
                                                setUnit={setUnit}
                                                hour={hour}
                                                minute={minute}
                                                setHour={setHour}
                                                setMinute={setMinute}
                                        />
                                )}
                                <ExpectedDocumentsList
                                        documents={expectedDocuments}
                                        onChange={setExpectedDocuments}
                                />
                                <div className="button-group">
                                        <ButtonPrimary
                                                text="Crează şablon"
                                                onClick={handleSubmit}
                                                type="button"
                                        />
                                </div>
                        </form>
                </>
        );
};

export default CreateTemplateForm;

const buildCron = (hour: string, minute: string, unit: RecurrenceUnit) => {
        const h = hour || "0";
        const m = minute || "0";
        switch (unit) {
                case "day":
                        return `${m} ${h} * * *`;
                case "week":
                        return `${m} ${h} * * 1`;
                case "month":
                        return `${m} ${h} 1 * *`;
                case "year":
                        return `${m} ${h} 1 1 *`;
        }
};
