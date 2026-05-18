import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import { Toaster } from "react-hot-toast";
import "./globals.css";
import AmplifyProvider from "@/components/AmplifyProvider";

const roboto = Roboto({
  variable: "--font-roboto",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Doclane",
  description: "Your supercharged document sharing platform.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${roboto.variable}`}>
        <AmplifyProvider>
          {children}
          <Toaster position="bottom-center" reverseOrder={false} />
        </AmplifyProvider>
      </body>
    </html>
  );
}