"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { User } from "@/types";
import Input from "@/components/InputComponents/Input";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import LocalitySearch from "./LocalitySearch";
import { updateUserProfile } from "@/lib/api/users";
import { formatDate } from "@/lib/client/formatDate";
import InfoList from "@/components/CardComponents/InfoList/InfoList";
import InfoItem from "@/components/CardComponents/InfoItem/InfoItem";
import toast from "react-hot-toast";
import "./ProfileSection.css";

interface Props {
        user: User;
}

export default function ProfileSection({ user }: Props) {
        const router = useRouter();
        const [phone, setPhone] = useState(user.phone ?? "");
        const [street, setStreet] = useState(user.street ?? "");
        const [locality, setLocality] = useState(user.locality ?? "");
        const [isSubmitting, setIsSubmitting] = useState(false);

        const handleSave = async () => {
                setIsSubmitting(true);
                const response = await updateUserProfile({
                        phone: phone.trim() || null,
                        street: street.trim() || null,
                        locality: locality.trim() || null,
                });
                setIsSubmitting(false);

                if (response.success) {
                        toast.success("Profil actualizat cu succes.");
                        router.refresh();
                } else {
                        toast.error(response.message ?? "Eroare la actualizarea profilului.");
                }
        };

        return (
                <div className="profile-section">
                        <div className="profile-card">
                                <h3 className="profile-card-title">Informații cont</h3>
                                <InfoList>
                                        <InfoItem
                                                label="Nume"
                                                value={`${user.first_name} ${user.last_name}`}
                                        />
                                        <InfoItem label="Email" value={user.email} />
                                        <InfoItem label="Rol" value={user.role} />
                                        <InfoItem
                                                label="Status"
                                                value={user.is_active ? "Activ" : "Dezactivat"}
                                        />
                                        <InfoItem
                                                label="Membru din"
                                                value={formatDate(user.created_at)}
                                        />
                                </InfoList>
                        </div>

                        <div className="profile-card">
                                <h3 className="profile-card-title">Informații contact</h3>
                                <div className="profile-form">
                                        <Input
                                                label="Telefon"
                                                placeholder="ex: 0712 345 678"
                                                value={phone}
                                                onChange={(e) => setPhone(e.target.value)}
                                                fullWidth
                                        />
                                        <Input
                                                label="Stradă"
                                                placeholder="ex: Strada Florilor nr. 12"
                                                value={street}
                                                onChange={(e) => setStreet(e.target.value)}
                                                fullWidth
                                        />
                                        <LocalitySearch
                                                value={locality}
                                                onChange={setLocality}
                                        />
                                        <div className="profile-save">
                                                <ButtonPrimary
                                                        text={
                                                                isSubmitting
                                                                        ? "Se salvează..."
                                                                        : "Salvează"
                                                        }
                                                        onClick={handleSave}
                                                        disabled={isSubmitting}
                                                />
                                        </div>
                                </div>
                        </div>
                </div>
        );
}
