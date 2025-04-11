// File: components/Paginator.tsx

import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationEllipsis,
} from "@/components/ui/pagination";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/button";

type PaginatorProps = {
  totalRecords: number;
  totalPages: number;
  currentPage: number;
  perPage: number;
  onPageChange: (newPage: number) => void;
};

// 🔹 Criando os botões customizados de "Anterior" e "Próximo"
const PaginationPreviousCustom = ({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) => (
  <PaginationLink
    aria-label="Página Anterior"
    size="default"
    className={cn("gap-1 pl-2.5", className)}
    {...props}
  >
    <ChevronLeft className="h-4 w-4" />
    <span>Anterior</span>
  </PaginationLink>
);

const PaginationNextCustom = ({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) => (
  <PaginationLink
    aria-label="Próxima Página"
    size="default"
    className={cn("gap-1 pr-2.5", className)}
    {...props}
  >
    <span>Próximo</span>
    <ChevronRight className="h-4 w-4" />
  </PaginationLink>
);

export default function Paginator({
  totalRecords,
  totalPages,
  currentPage,
  perPage,
  onPageChange,
}: PaginatorProps) {
  return (
    <div className="flex flex-col items-center mt-6">
      <p className="text-sm text-muted-foreground">
        Exibindo <span className="font-semibold">{(currentPage - 1) * perPage + 1}</span> -{" "}
        <span className="font-semibold">{Math.min(currentPage * perPage, totalRecords)}</span> de{" "}
        <span className="font-semibold">{totalRecords}</span> registros
      </p>

      <Pagination>
        <PaginationContent>
          {/* Botão "Anterior" */}
          <PaginationItem>
            <PaginationPreviousCustom
              className={cn({ "pointer-events-none opacity-50": currentPage === 1 })}
              onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
            />
          </PaginationItem>

          {currentPage > 2 && (
            <PaginationItem>
              <PaginationLink onClick={() => onPageChange(1)}>1</PaginationLink>
            </PaginationItem>
          )}

          {currentPage > 3 && (
            <PaginationItem>
              <PaginationEllipsis />
            </PaginationItem>
          )}

          {currentPage > 1 && (
            <PaginationItem>
              <PaginationLink onClick={() => onPageChange(currentPage - 1)}>
                {currentPage - 1}
              </PaginationLink>
            </PaginationItem>
          )}

          {/* Página atual destacada */}
          <PaginationItem>
            <PaginationLink isActive>{currentPage}</PaginationLink>
          </PaginationItem>

          {currentPage < totalPages && (
            <PaginationItem>
              <PaginationLink onClick={() => onPageChange(currentPage + 1)}>
                {currentPage + 1}
              </PaginationLink>
            </PaginationItem>
          )}

          {currentPage < totalPages - 2 && (
            <PaginationItem>
              <PaginationEllipsis />
            </PaginationItem>
          )}

          {currentPage < totalPages - 1 && (
            <PaginationItem>
              <PaginationLink onClick={() => onPageChange(totalPages)}>{totalPages}</PaginationLink>
            </PaginationItem>
          )}

          {/* Botão "Próximo" */}
          <PaginationItem>
            <PaginationNextCustom
              className={cn({ "pointer-events-none opacity-50": currentPage >= totalPages })}
              onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  );
}
