"use client";

import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import Input from "@/components/Input/Input";
import TextArea from "@/components/Input/TextArea";
import { useState } from "react";
import "./AddRequestForClientForm.css";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";

interface AddRequestForClientFormProps {
  id: string;
}

const AddRequestForClientForm: React.FC<AddRequestForClientFormProps> = ({ id }) => {
  const [requestName, setRequestName] = useState("");
  const [requestDescription, setRequestDescription] = useState("");

  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (requestName === "") {
      return;
    }

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
        onChange={(e: any) => {
          setRequestName(e.target.value);
        }}
      />
      <TextArea
        label="Request description"
        value={requestDescription}
        onChange={(e: any) => {
          setRequestDescription(e.target.value);
        }}
      />
      <ButtonPrimary text="Create request" />
    </form>
  );
};

export default AddRequestForClientForm;
