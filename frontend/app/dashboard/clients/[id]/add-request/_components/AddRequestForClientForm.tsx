"use client";

import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import Input from "@/components/Input/Input";
import TextArea from "@/components/Input/TextArea";
import { useState } from "react";
import "./AddRequestForClientForm.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import CronInput from "./CronInput";
import { RecurrenceUnit } from "@/types";
import CheckboxInput from "@/components/Input/CheckboxInput";

interface AddRequestForClientFormProps {
  id: string;
}

const AddRequestForClientForm: React.FC<AddRequestForClientFormProps> = ({ id }) => {
  const [requestName, setRequestName] = useState("");
  const [requestDescription, setRequestDescription] = useState("");

  const [isRecurring, setIsRecurring] = useState(false);
  const [unit, setUnit] = useState<RecurrenceUnit>("month");
  const [hour, setHour] = useState("09");
  const [minute, setMinute] = useState("00");

  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!requestName) return;

    const recurrenceCron = isRecurring ? buildCron(hour, minute, unit) : undefined;

    const createRequestPromise = fetch("/api/backend/document-requests", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: requestName,
        description: requestDescription,
        client_id: +id,
        ...(recurrenceCron ? { recurrence_cron: recurrenceCron } : {}),
      }),
    }).then(async (res) => {
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || "Failed to create request");
      }
      return res.json();
    });

    toast.promise(createRequestPromise, {
      loading: "Creating request...",
      success: "Request created successfully!",
      error: (err) => `Failed: ${err.message}`,
    });

    createRequestPromise.then(() => {
      setRequestName("");
      setRequestDescription("");
      router.push("/dashboard/clients");
    });
  };

  return (
    <form className="add-request-for-client-form" onSubmit={handleSubmit}>
      <Input
        label="Request title"
        value={requestName}
        onChange={(e: any) => setRequestName(e.target.value)}
      />

      <TextArea
        label="Request description"
        value={requestDescription}
        onChange={(e: any) => setRequestDescription(e.target.value)}
      />

      <CheckboxInput
        isChecked={isRecurring}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          setIsRecurring(e.target.checked);
        }}
        label="Is recurring"
      />

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

      <ButtonPrimary text="Create request" />
    </form>
  );
};

export default AddRequestForClientForm;

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
