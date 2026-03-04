"use client";

import { useEffect } from "react";
import "./error.css";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";

export default function GlobalError({
        error,
        reset,
}: {
        error: Error & { digest?: string };
        reset: () => void;
}) {
        useEffect(() => {
                console.error("Pagina a întâmpinat o eroare:", error);
        }, [error]);

        return (
                <div className="error-page">
                        <img src="/error.jpg" alt="" height={150} />
                        <h1 className="error-title">Oops! Something went wrong...</h1>

                        <p className="error-message">
                                {error.message ||
                                        "An error appeared in the system. And we don't know what went wrong."}
                        </p>

                        <div className="button-wrapper">
                                <ButtonPrimary text="Try again" onClick={() => reset()} />
                        </div>
                </div>
        );
}

function extractStatusCode(message: string): string {
        const match = message.match(/\d{3}/);
        return match ? match[0] : "UNK";
}
