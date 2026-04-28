"use client";
import "./error.css";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";

export default function NotFoundPage() {
        return (
                <div className="error-page">
                        <div className="error-card">
                                <h1 className="error-code">404</h1>
                                <p className="error-title">Pagina nu a fost găsită.</p>
                                <p className="error-message">
                                        Ne pare rău, pagina pe care o cauți nu există.
                                </p>
                                <div className="button-back">
                                        <ButtonPrimary
                                                text="Înapoi"
                                                fullWidth
                                                onClick={() => {
                                                        window.location.href =
                                                                "/dashboard/requests";
                                                }}
                                        />
                                </div>
                        </div>
                </div>
        );
}
