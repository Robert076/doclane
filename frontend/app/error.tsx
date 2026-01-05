"use client";

import { useEffect } from "react";
import "./error.css";
import Button from "@/components/Button/Button";

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
    <div className="error-page-wrapper">
      <div className="error-page">
        <h1 className="error-title">Oops! Something went wrong...</h1>

        <p className="error-message">
          {error.message ||
            "An error appeared in the system. And we don't know what went wrong."}
        </p>

        <div style={{ marginTop: "1rem" }}>
          <span className="error-status-code">Status: {extractStatusCode(error.message)}</span>
        </div>

        {error.digest && <span className="error-digest">Error ID: {error.digest}</span>}

        <Button text="Try again" onClick={() => reset()} />
      </div>
    </div>
  );
}

function extractStatusCode(message: string): string {
  const match = message.match(/\d{3}/);
  return match ? match[0] : "UNK";
}
