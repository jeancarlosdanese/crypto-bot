// File: context/UserContext.tsx

import { createContext, useContext, useEffect, useState } from "react";
import { UserService } from "@/services/user";
import { useRouter } from "next/router";
import { toast } from "sonner";

// Tipo correto baseado no retorno do backend
type Account = {
  id: string; // UUID
  name: string;
  email: string;
  whatsapp: string;
};

type UserContextType = {
  user: Account | null;
  loading: boolean;
  refreshUser: () => Promise<void>;
  logout: () => void;
};

const UserContext = createContext<UserContextType>({
  user: null,
  loading: true,
  refreshUser: async () => {},
  logout: () => {},
});

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<Account | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const fetchUser = async () => {
    if (typeof window === "undefined") return;

    setLoading(true);

    const token = localStorage.getItem("token");
    if (!token) {
      setUser(null);
      router.push("/auth/login");
      return;
    }

    try {
      const userData: Account = await UserService.getAuthenticatedUser();
      setUser(userData);
    } catch (error: any) {
      console.error("Erro ao buscar usuário autenticado:", error);

      if (error.response?.status === 401) {
        localStorage.removeItem("token");
        router.push("/auth/login");
      } else {
        toast.error("Erro ao carregar dados do usuário. Tente novamente mais tarde.");
      }

      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    router.push("/auth/login");
  };

  useEffect(() => {
    fetchUser();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <UserContext.Provider value={{ user, loading, refreshUser: fetchUser, logout }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => useContext(UserContext);
