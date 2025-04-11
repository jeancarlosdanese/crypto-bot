// File: /components/Sidebar.tsx

import Link from "next/link";
import { useRouter } from "next/router";
import { useUser } from "@/context/UserContext";
import { Home, FileText, Megaphone, Users, Settings, LogOut, Menu } from "lucide-react";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

// 🔹 Definição do menu de navegação
const menuItems = [
  { name: "Dashboard", href: "/dashboard", icon: <Home className="w-5 h-5" /> },
  { name: "Configurações", href: "/settings", icon: <Settings className="w-5 h-5" /> },
];

export function Sidebar() {
  const router = useRouter();
  const { user, logout } = useUser();

  return (
    <>
      {/* 🔹 Botão para abrir a sidebar no mobile */}
      <div className="sm:hidden fixed top-4 left-4 z-50">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="outline" size="icon">
              <Menu className="w-5 h-5" />
            </Button>
          </SheetTrigger>

          <SheetContent side="left" className="w-64 bg-background border-r border-border">
            {/* 🔹 Acessibilidade - Adicionando título e descrição */}
            <SheetHeader>
              <SheetTitle>hy-marketing</SheetTitle>
              <SheetDescription>Acesse as seções do sistema</SheetDescription>
            </SheetHeader>

            <SidebarContent user={user} logout={logout} router={router} />
          </SheetContent>
        </Sheet>
      </div>

      {/* 🔹 Sidebar fixa no desktop */}
      <aside className="hidden sm:flex sm:flex-col w-64 min-h-screen bg-background border-r border-border shadow-md">
        <SidebarContent user={user} logout={logout} router={router} />
      </aside>
    </>
  );
}

// 🔹 Conteúdo da Sidebar (Evita código duplicado)
function SidebarContent({ user, logout, router }: { user: any; logout: () => void; router: any }) {
  return (
    <div className="h-full flex flex-col p-4 overflow-y-auto">
      <h1 className="text-xl font-bold">Next-Marketing</h1>

      {/* 🔹 Menus */}
      <nav className="flex flex-col gap-2 mt-4 flex-grow">
        {menuItems.map((item) => (
          <Link key={item.href} href={item.href}>
            <Button
              variant={router.pathname === item.href ? "default" : "ghost"}
              className={cn("w-full flex items-center gap-3 justify-start")}
            >
              {item.icon}
              {item.name}
            </Button>
          </Link>
        ))}
      </nav>

      {/* 🔹 Botão de Logout fixo no final */}
      {user && (
        <div className="mt-auto">
          <Button
            variant="destructive"
            className="w-full flex items-center gap-3 justify-start"
            onClick={() => {
              logout();
              router.push("/auth/login");
            }}
          >
            <LogOut className="w-5 h-5" />
            Sair
          </Button>
        </div>
      )}
    </div>
  );
}
