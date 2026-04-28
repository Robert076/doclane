"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { updatePassword } from "@/lib/api/users";
import toast from "react-hot-toast";
import "@/components/Pages/SettingsComponents/ProfileSection.css";

export default function PasswordSection() {
        const router = useRouter();
        const [currentPassword, setCurrentPassword] = useState("");
        const [newPassword, setNewPassword] = useState("");
        const [confirmPassword, setConfirmPassword] = useState("");
        const [isSubmitting, setIsSubmitting] = useState(false);

        const handleSave = async () => {
                if (!currentPassword || !newPassword || !confirmPassword) {
                        toast.error("Completează toate câmpurile.");
                        return;
                }
                if (newPassword.length < 8) {
                        toast.error("Parola nouă trebuie să aibă cel puțin 8 caractere.");
                        return;
                }
                if (newPassword !== confirmPassword) {
                        toast.error("Parolele nu coincid.");
                        return;
                }
                if (currentPassword === newPassword) {
                        toast.error("Parola nouă trebuie să fie diferită de cea actuală.");
                        return;
                }

                setIsSubmitting(true);
                const response = await updatePassword({
                        current_password: currentPassword,
                        new_password: newPassword,
                });
                setIsSubmitting(false);

                if (response.success) {
                        toast.success(
                                "Parolă schimbată cu succes. Te rugăm să te autentifici din nou.",
                        );
                        setCurrentPassword("");
                        setNewPassword("");
                        setConfirmPassword("");
                        router.push("/login");
                } else {
                        toast.error(response.message ?? "Eroare la schimbarea parolei.");
                }
        };

        return (
                <div className="profile-card">
                        <h3 className="profile-card-title">Schimbă parola</h3>
                        <div className="profile-form">
                                <Input
                                        label="Parola actuală"
                                        placeholder="Parola ta actuală"
                                        value={currentPassword}
                                        onChange={(e) => setCurrentPassword(e.target.value)}
                                        isPassword
                                        fullWidth
                                />
                                <Input
                                        label="Parola nouă"
                                        placeholder="Minimum 8 caractere"
                                        value={newPassword}
                                        onChange={(e) => setNewPassword(e.target.value)}
                                        isPassword
                                        fullWidth
                                />
                                <Input
                                        label="Confirmă parola nouă"
                                        placeholder="Repetă parola nouă"
                                        value={confirmPassword}
                                        onChange={(e) => setConfirmPassword(e.target.value)}
                                        isPassword
                                        fullWidth
                                />
                                <div className="profile-save">
                                        <ButtonPrimary
                                                text={
                                                        isSubmitting
                                                                ? "Se salvează..."
                                                                : "Schimbă parola"
                                                }
                                                onClick={handleSave}
                                                disabled={isSubmitting}
                                        />
                                </div>
                        </div>
                </div>
        );
}
