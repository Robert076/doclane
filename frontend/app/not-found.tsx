"use client";
import "./error.css";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";

export default function NotFoundPage() {
        return (
                <div className="error-page">
                        <h1 className="error-code">404</h1>
                        <p className="error-message">Pagina nu a fost găsită.</p>
                        <div className="button-back">
                                <ButtonPrimary
                                        text="Înapoi"
                                        fullWidth
                                        onClick={() => {
                                                window.location.href = "/dashboard/requests";
                                        }}
                                />
                        </div>
                </div>
        );
}
