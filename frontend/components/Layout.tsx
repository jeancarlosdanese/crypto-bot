// File: components/Layout.tsx

import { Sidebar } from "@/components/Sidebar";
import { Header } from "@/components/Header";
import Footer from "@/components/Footer";

type LayoutProps = {
  children: React.ReactNode;
};

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground">
      <div className="flex flex-1">
        <Sidebar />
        <div className="flex flex-col flex-1">
          <Header />
          <main className="flex-grow p-6">{children}</main>
        </div>
      </div>
    </div>
  );
}
