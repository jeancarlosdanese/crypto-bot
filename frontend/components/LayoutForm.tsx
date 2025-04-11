// File: components/LayoutForm.tsx

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useRouter } from "next/router";

type LayoutFormProps = {
  title?: string;
  children: React.ReactNode;
  onSave: () => void;
  onCancel?: () => void;
};

export default function LayoutForm({ title, children, onSave, onCancel }: LayoutFormProps) {
  const router = useRouter();

  return (
    <div className="p-0 sm:p-6 w-full max-w-4xl mx-auto">
      <Card className="w-full">
        {title ? (
          <CardHeader>
            <CardTitle className="text-lg sm:text-xl">{title}</CardTitle>
          </CardHeader>
        ) : (
          <div className="p-4"></div>
        )}
        <CardContent>
          {children}
          <div className="flex justify-end gap-4 mt-4 sm:mt-6">
            <Button variant="outline" onClick={onCancel || (() => router.back())}>
              Cancelar
            </Button>
            <Button onClick={onSave}>Salvar</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
