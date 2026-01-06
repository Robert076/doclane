import { notFound } from "next/navigation";
import getRequestById from "@/lib/getRequestById";
import getFilesByRequestId from "@/lib/getFilesByRequestId";
import FileSection from "@/components/FileSection/FileSection";
import DetailsHeader from "./_components/DetailsHeader";
import DetailsCard from "./_components/DetailsCard";
import DetailCardsActionSidebar from "./_components/DetailCardsActionSidebar";
import "./style.css";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default async function RequestDetailsPage({ params }: PageProps) {
  const { id } = await params;

  const [request, filesResponse] = await Promise.all([
    getRequestById(id),
    getFilesByRequestId(id),
  ]);

  if (!request || !request.data) {
    notFound();
  }

  const data = request.data;
  const files = filesResponse?.data || [];

  return (
    <div className="details-container">
      <DetailsHeader data={data} />
      <div className="details-grid">
        <div className="main-content">
          <DetailsCard data={data} />
          <FileSection files={files} />
        </div>
        <DetailCardsActionSidebar id={id} />
      </div>
    </div>
  );
}
