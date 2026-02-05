"use client";
import { useState } from "react";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import "./InvitationCodeGenerator.css";
import { MdContentCopy, MdCheck, MdClose } from "react-icons/md";

const InvitationCodeGenerator = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [code, setCode] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const generateCode = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch("/api/backend/invitations/generate", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          expires_in_days: 7,
        }),
      });

      if (!response.ok) {
        throw new Error(
          "Failed to generate invitation code. Please make sure you have less than 3 active codes.",
        );
      }

      const data = await response.json();
      setCode(data.data.code);
    } catch (err) {
      setError(
        "Failed to generate invitation code. Please make sure you have less than 3 active codes.",
      );
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = async () => {
    if (!code) return;

    try {
      await navigator.clipboard.writeText(code);
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setCode(null);
    setIsCopied(false);
    setError(null);
  };

  return (
    <>
      <ButtonPrimary text="Add New Client" onClick={() => setIsModalOpen(true)} />

      {isModalOpen && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button className="modal-close" onClick={closeModal} aria-label="Close">
              <MdClose size={24} />
            </button>

            {!code ? (
              <>
                <h2 className="modal-title">Generate Invitation Code</h2>
                <p className="modal-description">
                  This will create a one-time invitation code that expires in 7 days. Share
                  this code with your client so they can register and access their document
                  portal.
                </p>
                <p className="modal-note">
                  Make sure you're ready to copy and send the code immediately after
                  generation.
                </p>
                <div className="modal-actions">
                  <ButtonPrimary
                    text={isLoading ? "Generating..." : "Continue"}
                    onClick={generateCode}
                    disabled={isLoading}
                  />
                  <ButtonPrimary text="Cancel" variant="ghost" onClick={closeModal} />
                </div>
              </>
            ) : (
              <>
                <h2 className="modal-title">Invitation Code Generated</h2>
                <p className="modal-description">
                  Share this code with your client. They'll use it to register on the platform.
                </p>
                <div className="code-box">
                  <span className="code-text">{code}</span>
                  <button
                    className="copy-button"
                    onClick={copyToClipboard}
                    aria-label="Copy code"
                  >
                    {isCopied ? <MdCheck size={20} /> : <MdContentCopy size={20} />}
                  </button>
                </div>
                <p className="code-expiry">Expires in 7 days</p>
                <div className="modal-actions">
                  <ButtonPrimary text="Done" onClick={closeModal} />
                </div>
              </>
            )}

            {error && <p className="error-message">{error}</p>}
          </div>
        </div>
      )}
    </>
  );
};

export default InvitationCodeGenerator;
