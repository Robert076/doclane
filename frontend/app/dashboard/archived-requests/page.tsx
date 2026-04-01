import { notFound } from "next/navigation";
import { Request } from "@/types";
import ArchivedRequestsSection from "@/components/Pages/ArchivedRequestsComponents/ArchivedRequestsSection";
import PageHeader from "@/components/PageHeader/PageHeader";

const ArchivedRequests = async () => {
        // if (!requestsResponse?.data) {
        //         notFound();
        // }
        // const requests = requestsResponse.data as Request[];
        // return (
        //         <div>
        //                 <PageHeader
        //                         title="Dosare arhivate"
        //                         subtitle="Restaurează şi gestionează dosarele arhivate."
        //                 />
        //                 <ArchivedRequestsSection requests={requests} />
        //         </div>
        // );
};

export default ArchivedRequests;
