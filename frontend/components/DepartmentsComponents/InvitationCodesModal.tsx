"use client";
import { useState, useEffect } from "react";
import { InvitationCode } from "@/types";
import { getAllInvitationCodes } from "@/lib/api/invitation_codes";
import { formatDate } from "@/lib/client/formatDate";
import Modal from "@/components/Modals/Modal";
import NotFound from "@/components/OtherComponents/NotFound/NotFound";
import toast from "react-hot-toast";
import { MdContentCopy } from "react-icons/md";
import "./InvitationCodesModal.css";

interface Props {
  isOpen: boolean;
  onClose: () => void;
}

type Tab = "active" | "used";

export default function InvitationCodesModal({ isOpen, onClose }: Props) {
  const [codes, setCodes] = useState<InvitationCode[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [tab, setTab] = useState<Tab>("active");

  useEffect(() => {
    if (!isOpen) return;
    setIsLoading(true);
    getAllInvitationCodes().then((res) => {
      setCodes(res.data ?? []);
      setIsLoading(false);
    });
  }, [isOpen]);

  const now = new Date();

  const activeCodes = codes.filter(
    (c) => !c.used_at && (!c.expires_at || new Date(c.expires_at) > now),
  );

  const usedCodes = codes.filter((c) => !!c.used_at);

  const expiredCodes = codes.filter(
    (c) => !c.used_at && c.expires_at && new Date(c.expires_at) <= now,
  );

  const currentCodes = tab === "active" ? activeCodes : usedCodes;

  const handleCopy = (code: string) => {
    const link = `${window.location.origin}/register/invite?code=${code}`;
    navigator.clipboard.writeText(link);
    toast.success("Link copiat!");
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Coduri de invitație" hideFooter>
      <div className="inv-tabs">
        <button
          className={`inv-tab ${tab === "active" ? "inv-tab--active" : ""}`}
          onClick={() => setTab("active")}
        >
          Active ({activeCodes.length})
        </button>
        <button
          className={`inv-tab ${tab === "used" ? "inv-tab--active" : ""}`}
          onClick={() => setTab("used")}
        >
          Folosite ({usedCodes.length})
        </button>
      </div>

      {isLoading ? (
        <p className="inv-loading">Se încarcă...</p>
      ) : currentCodes.length === 0 ? (
        <NotFound
          text={tab === "active" ? "Niciun cod activ" : "Niciun cod folosit"}
          subtext={
            tab === "active"
              ? "Generează un cod nou dintr-un departament."
              : "Codurile folosite vor apărea aici."
          }
          background="#fff"
        />
      ) : (
        <div className="inv-list">
          {currentCodes.map((code) => (
            <div key={code.id} className="inv-item">
              <div className="inv-item-left">
                <span className="inv-code">{code.code}</span>
                <span className="inv-dept">{code.department_name}</span>
                {tab === "active" && code.expires_at && (
                  <span className="inv-meta">
                    Expiră: {formatDate(code.expires_at)}
                  </span>
                )}
                {tab === "used" && (
                  <>
                    {code.used_by_first_name && (
                      <span className="inv-used-by">
                        Folosit de: {code.used_by_first_name} {code.used_by_last_name}
                      </span>
                    )}
                    {code.used_by_email && (
                      <span className="inv-meta">{code.used_by_email}</span>
                    )}
                    {code.used_at && (
                      <span className="inv-meta">
                        Folosit la: {formatDate(code.used_at)}
                      </span>
                    )}
                  </>
                )}
              </div>
              {tab === "active" && (
                <button className="inv-copy-btn" onClick={() => handleCopy(code.code)}>
                  <MdContentCopy size={18} />
                </button>
              )}
              {tab === "used" && (
                <span className="inv-used-badge">Folosit</span>
              )}
            </div>
          ))}
        </div>
      )}

      {tab === "active" && expiredCodes.length > 0 && (
        <p className="inv-expired-note">
          {expiredCodes.length} {expiredCodes.length === 1 ? "cod expirat" : "coduri expirate"} ascunse.
        </p>
      )}
    </Modal>
  );
}