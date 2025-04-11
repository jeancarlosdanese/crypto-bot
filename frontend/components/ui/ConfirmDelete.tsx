// File: components/ui/ConfirmDelete.tsx

import { useState } from "react";
import {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogTitle,
  AlertDialogDescription,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import Spinner from "../Spinner";

interface ConfirmDeleteProps {
  onConfirm: () => void;
  entityName: string;
  label?: string;
  disabled?: boolean;
}

export function ConfirmDelete({ onConfirm, entityName, label, disabled }: ConfirmDeleteProps) {
  const [open, setOpen] = useState(false);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="destructive" disabled={disabled} onClick={() => setOpen(true)}>
          {disabled ? "Excluindo..." : label ? label : "Excluir"}
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent className="bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 border border-border shadow-lg rounded-lg">
        <AlertDialogHeader>
          <AlertDialogTitle>Excluir {entityName}</AlertDialogTitle>
          <AlertDialogDescription>
            Você tem certeza que deseja excluir {entityName}? Essa ação não pode ser desfeita.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={() => {
              onConfirm();
              setOpen(false);
            }}
            disabled={disabled}
          >
            {disabled ? <Spinner className="w-4 h-4" /> : "Excluir"}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
