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
                        <div className="error-card">
                                <img src="/error.jpg" alt="" />
                                <h1 className="error-title">Oops! Something went wrong...</h1>
                                <p className="error-message">
                                        {error.message || "An error appeared in the system."}
                                </p>
                                <div className="button-wrapper">
                                        <ButtonPrimary
                                                text="Try again"
                                                fullWidth
                                                onClick={() => reset()}
                                        />
                                </div>
                        </div>
                </div>
        );
}
