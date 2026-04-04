"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { Department, RecurrenceUnit } from "@/types";
import { createTemplate } from "@/lib/api/templates";
import { buildCronExpression } from "@/lib/cron";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import RadioInput from "@/components/InputComponents/RadioInput";
import CronInput from "@/components/InputComponents/CronInput";
import ExpectedDocumentsList, {
        ExpectedDocumentInput,
} from "@/components/InputComponents/ExpectedDocumentList";
import "./CreateTemplateForm.css";
import Select from "@/components/InputComponents/Select";

interface Props {
        departments: Department[];
}

const EMPTY_DOCUMENT: ExpectedDocumentInput = { title: "", description: "" };

export default function CreateTemplateForm({ departments }: Props) {
        const router = useRouter();

        const [title, setTitle] = useState("");
        const [description, setDescription] = useState("");
        const [departmentID, setDepartmentID] = useState<number | null>(null);
        const [isRecurring, setIsRecurring] = useState(false);
        const [unit, setUnit] = useState<RecurrenceUnit>("month");
        const [hour, setHour] = useState("09");
        const [minute, setMinute] = useState("00");
        const [expectedDocuments, setExpectedDocuments] = useState<ExpectedDocumentInput[]>([
                EMPTY_DOCUMENT,
        ]);
        const [isSubmitting, setIsSubmitting] = useState(false);

        const validate = (): boolean => {
                if (!title.trim()) {
                        toast.error("Completează titlul șablonului.");
                        return false;
                }
                if (!departmentID) {
                        toast.error("Selectează un departament.");
                        return false;
                }
                if (expectedDocuments.some((doc) => !doc.title.trim())) {
                        toast.error("Completează titlul fiecărui document solicitat.");
                        return false;
                }
                return true;
        };

        const handleSubmit = async (e: React.FormEvent) => {
                e.preventDefault();
                if (!validate() || isSubmitting) return;

                setIsSubmitting(true);
                const res = await createTemplate({
                        title: title.trim(),
                        description: description.trim() || undefined,
                        department_id: departmentID!,
                        is_recurring: isRecurring,
                        recurrence_cron: isRecurring
                                ? buildCronExpression(unit, hour, minute)
                                : undefined,
                        expected_documents: expectedDocuments.map((doc) => ({
                                title: doc.title.trim(),
                                description: doc.description.trim(),
                                example_file: doc.exampleFile,
                        })),
                });
                setIsSubmitting(false);

                if (res.success) {
                        toast.success("Șablon creat cu succes!");
                        router.push("/dashboard/templates");
                } else {
                        toast.error(res.error ?? "A apărut o eroare.");
                }
        };

        return (
                <form className="add-form" onSubmit={handleSubmit}>
                        <Input
                                label="Titlul șablonului"
                                placeholder="Scrie titlul șablonului..."
                                value={title}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                        setTitle(e.target.value)
                                }
                        />
                        <TextArea
                                label="Descrierea șablonului"
                                placeholder="Scrie descrierea șablonului..."
                                value={description}
                                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                                        setDescription(e.target.value)
                                }
                        />
                        <Select
                                label="Departament"
                                value={departmentID ?? ""}
                                onChange={(val) => setDepartmentID(Number(val) || null)}
                                placeholder="Selectează departamentul..."
                                options={departments.map((d) => ({
                                        value: d.id,
                                        label: d.name,
                                }))}
                        />
                        <div className="radio-inputs-time">
                                <RadioInput
                                        isChecked={!isRecurring}
                                        onChange={() => setIsRecurring(false)}
                                        label="Fără recurență"
                                />
                                <RadioInput
                                        isChecked={isRecurring}
                                        onChange={() => setIsRecurring(true)}
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
                                        text={
                                                isSubmitting
                                                        ? "Se creează..."
                                                        : "Creează șablon"
                                        }
                                        type="submit"
                                        disabled={isSubmitting}
                                />
                        </div>
                </form>
        );
}
