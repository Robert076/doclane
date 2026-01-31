import "./style.css";
import AddRequestForClientHeader from "./_components/AddRequestForClientHeader";
import AddRequestForClientForm from "./_components/AddRequestForClientForm";

interface PageProps {
  params: Promise<{ id: string }>;
}

const AddRequestForClient = async ({ params }: PageProps) => {
  const { id } = await params;

  return (
    <div className="add-request-for-client">
      <AddRequestForClientHeader id={id} />
      <AddRequestForClientForm id={id} />
    </div>
  );
};

export default AddRequestForClient;
