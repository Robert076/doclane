"use client";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import Input from "@/components/InputComponents/Input";
import TextArea from "@/components/InputComponents/TextArea";
import { useState, useEffect } from "react";
import "./CreateTemplateForm.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import { RecurrenceUnit, Department } from "@/types";
import RadioInput from "@/components/InputComponents/RadioInput";
import { createTemplate } from "@/lib/api/templates";
import { getDepartments } from "@/lib/api/departments";
import { buildCronExpression } from "@/lib/cron";
import CronInput from "@/components/InputComponents/CronInput";
import ExpectedDocumentsList, {
        ExpectedDocumentInput,
} from "@/components/InputComponents/ExpectedDocumentList";

const CreateTemplateForm = () => {
        const [title, setTitle] = useState("");
        const [description, setDescription] = useState("");
        const [departmentID, setDepartmentID] = useState<number | null>(null);
        const [departments, setDepartments] = useState<Department[]>([]);
        const [isNoneSelected, setIsNoneSelected] = useState(true);
        const [isRecurring, setIsRecurring] = useState(false);
        const [unit, setUnit] = useState<RecurrenceUnit>("month");
        const [hour, setHour] = useState("09");
        const [minute, setMinute] = useState("00");
        const [expectedDocuments, setExpectedDocuments] = useState<ExpectedDocumentInput[]>([
                { title: "", description: "" },
        ]);
        const router = useRouter();

        useEffect(() => {
                getDepartments().then((res) => {
                        if (res.success && res.data) setDepartments(res.data);
                });
        }, []);

        const validateForm = () => {
                if (!title) {
                        toast.error("Completează titlul șablonului.");
                        return false;
                }
                if (!departmentID) {
                        toast.error("Selectează un departament.");
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
                                const res = await createTemplate({
                                        title,
                                        description: description || undefined,
                                        department_id: departmentID!,
                                        is_recurring: isRecurring,
                                        recurrence_cron: isRecurring
                                                ? buildCronExpression(unit, hour, minute)
                                                : undefined,
                                        expected_documents: expectedDocuments.map((ed) => ({
                                                title: ed.title,
                                                description: ed.description,
                                                example_file: ed.exampleFile,
                                        })),
                                });
                                if (!res.success) throw new Error(res.error);
                                router.push("/dashboard/templates");
                        })(),
                        {
                                loading: "Se creează șablonul...",
                                success: "Șablon creat cu succes!",
                                error: (err) => `Eroare: ${err.message}`,
                        },
                );
        };

        return (
                <form className="add-form" onSubmit={handleSubmit}>
                        <Input
                                label="Titlul șablonului"
                                placeholder="Scrie titlul șablonului..."
                                value={title}
                                onChange={(e: any) => setTitle(e.target.value)}
                        />
                        <TextArea
                                label="Descrierea șablonului"
                                placeholder="Scrie descrierea șablonului..."
                                value={description}
                                onChange={(e: any) => setDescription(e.target.value)}
                        />
                        <div className="form-field">
                                <label className="form-label">Departament</label>
                                <select
                                        className="form-select"
                                        value={departmentID ?? ""}
                                        onChange={(e) =>
                                                setDepartmentID(Number(e.target.value))
                                        }
                                >
                                        <option value="">Selectează departamentul...</option>
                                        {departments.map((d) => (
                                                <option key={d.id} value={d.id}>
                                                        {d.name}
                                                </option>
                                        ))}
                                </select>
                        </div>
                        <div className="radio-inputs-time">
                                <RadioInput
                                        isChecked={isNoneSelected}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                                setIsNoneSelected(e.target.checked);
                                                setIsRecurring(false);
                                        }}
                                        label="Fără recurență"
                                />
                                <RadioInput
                                        isChecked={isRecurring}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
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
                                        text="Creează șablon"
                                        onClick={handleSubmit}
                                        type="button"
                                />
                        </div>
                </form>
        );
};

export default CreateTemplateForm;
