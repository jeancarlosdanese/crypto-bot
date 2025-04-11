// File: types/account.ts

export type Account = {
  id: string;
  name: string;
  email: string;
};

export type AccountContextType = {
  user: Account | null;
  setUser: (user: Account | null) => void;
  logout: () => void;
};
