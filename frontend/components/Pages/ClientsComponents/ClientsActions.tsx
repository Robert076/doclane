import InvitationCodeGenerator from "@/components/ClientComponents/Invitation/InvitationCodeGenerator/InvitationCodeGenerator";
import InvitationCodesModal from "@/components/ClientComponents/Invitation/InvitationCodesModal/InvitationCodesModal";
import "./ClientsActions.css";

const ClientsActions = () => {
        return (
                <div className="clients-actions has-margin-bottom">
                        <InvitationCodeGenerator />
                        <InvitationCodesModal />
                </div>
        );
};

export default ClientsActions;
