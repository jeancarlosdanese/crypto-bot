// File: components/Campaign/CampaignAudience.tsx

import { useEffect, useState, useCallback, useMemo } from "react";
import { CampaignService } from "@/services/campaign";
import { Contact } from "@/types/contact";
import { Audience } from "@/types/audience";

// import { debounce } from "lodash"; // para debouce se preferir
import Spinner from "@/components/Spinner";
import Paginator from "@/components/Paginator";

import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CheckSquare, PlusCircle, Square, Trash2, UserCheck, UserPlus } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { log } from "console";

type Filters = {
  name: string;
  email: string;
  whatsapp: string;
  gender: string;
  birth_date_start: string;
  birth_date_end: string;
  bairro: string;
  cidade: string;
  estado: string;
  tags: string;
};

type PaginationState = {
  total_records: number;
  total_pages: number;
  current_page: number;
  per_page: number;
};

/**
 * Componente que implementa toda a lógica de audiência
 * com filtros, paginação e manipulação (adicionar / remover contatos).
 */
const CampaignAudience = ({ campaignId, status }: { campaignId: string; status: string }) => {
  // -------------------------------------
  // Estados de lista / audiência
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [audiences, setAudiences] = useState<Audience[]>([]);
  const [selectedContacts, setSelectedContacts] = useState<string[]>([]);
  const [selectAll, setSelectAll] = useState(false);
  // Se estiver pendente ou cancelada, permitir edição
  const isEditing = useMemo(() => status === "pendente" || status === "cancelada", [status]);

  // -------------------------------------
  // Estados de filtros e paginação unificados
  const [filters, setFilters] = useState<Filters>({
    name: "",
    email: "",
    whatsapp: "",
    gender: "",
    birth_date_start: "",
    birth_date_end: "",
    bairro: "",
    cidade: "",
    estado: "",
    tags: "",
  });

  const [pagination, setPagination] = useState<{
    available: PaginationState;
    audience: PaginationState;
  }>({
    available: {
      total_records: 0,
      total_pages: 1,
      current_page: 1,
      per_page: 12,
    },
    audience: {
      total_records: 0,
      total_pages: 1,
      current_page: 1,
      per_page: 12,
    },
  });

  // 1) Mantenha 'isInitialLoad' apenas para o carregamento inicial
  const [isInitialLoad, setIsInitialLoad] = useState(true);

  // 2) Crie 'isFetching' para buscas nos filtros
  const [isFetching, setIsFetching] = useState(false);

  // -------------------------------------
  // Debounce para filtros (500ms)
  const [debouncedFilters, setDebouncedFilters] = useState<Filters>(filters);

  // Debounce
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedFilters(filters);
    }, 500);

    return () => clearTimeout(handler);
  }, [filters]);

  // useEffect único para buscar dados
  useEffect(() => {
    if (!campaignId) return;

    if (isInitialLoad) {
      // Carga inicial: exibe spinner global
      Promise.all([fetchContacts(), fetchAudience()]).finally(() => {
        setIsInitialLoad(false);
      });
    } else {
      // Buscas subsequentes: não remove toda a UI; usa spinner parcial (isFetching)
      setIsFetching(true);
      Promise.all([fetchContacts(), fetchAudience()]).finally(() => setIsFetching(false));
    }
  }, [
    campaignId,
    debouncedFilters,
    pagination.available.current_page,
    pagination.audience.current_page,
  ]);

  // -------------------------------------
  // Funções de busca de dados (useCallback p/ evitar recriações)
  const fetchContacts = useCallback(async () => {
    try {
      const resp = await CampaignService.getAvailableContacts(
        campaignId,
        debouncedFilters,
        pagination.available.current_page,
        pagination.available.per_page
      );
      setContacts(resp.data);
      setPagination((prev) => ({
        ...prev,
        available: {
          total_records: resp.total_records,
          total_pages: resp.total_pages,
          current_page: resp.current_page,
          per_page: resp.per_page,
        },
      }));
    } catch (error) {
      toast.error("Erro ao carregar contatos disponíveis.");
    }
  }, [
    campaignId,
    debouncedFilters,
    pagination.available.current_page,
    pagination.available.per_page,
  ]);

  const fetchAudience = useCallback(async () => {
    try {
      const resp = await CampaignService.getAudienceContacts(
        campaignId,
        pagination.audience.current_page,
        pagination.audience.per_page
      );
      setAudiences(resp.data);
      setPagination((prev) => ({
        ...prev,
        audience: {
          total_records: resp.total_records,
          total_pages: resp.total_pages,
          current_page: resp.current_page,
          per_page: resp.per_page,
        },
      }));
    } catch (error) {
      toast.error("Erro ao carregar contatos da audiência.");
    }
  }, [campaignId, pagination.audience.current_page, pagination.audience.per_page]);

  // -------------------------------------
  // Funções genéricas de manipulação de página (disponível e audiência)
  const handleAvailablePageChange = (newPage: number) => {
    setPagination((prev) => ({
      ...prev,
      available: { ...prev.available, current_page: newPage },
    }));
  };

  const handleAudiencePageChange = (newPage: number) => {
    setPagination((prev) => ({
      ...prev,
      audience: { ...prev.audience, current_page: newPage },
    }));
  };

  // -------------------------------------
  // Funções de manipulação de Filtros
  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFilters((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  // -------------------------------------
  // Ações de adicionar / remover
  const addAllContactsToAudience = async () => {
    if (contacts.length === 0) return;

    // Limpamos filtros "vazios"
    const cleanFilters = Object.fromEntries(
      Object.entries(debouncedFilters).filter(
        ([, value]) => value !== undefined && value !== "" && value !== "none"
      )
    );

    const filtersToAdd = {
      filters: cleanFilters,
      current_page: pagination.available.current_page,
      per_page: pagination.available.per_page,
    };

    try {
      const resp = await CampaignService.addAllFilteredContactsToAudience(campaignId, filtersToAdd);
      fetchContacts();
      setAudiences(resp.data);
      setPagination((prev) => ({
        ...prev,
        audience: {
          total_records: resp.total_records,
          total_pages: resp.total_pages,
          current_page: resp.current_page,
          per_page: resp.per_page,
        },
      }));
      toast.success("Todos os contatos foram adicionados com sucesso!");
    } catch (error) {
      toast.error("Erro ao adicionar contatos à audiência.");
    }
  };

  const addContactsToAudience = async () => {
    if (selectedContacts.length === 0) return;
    try {
      await CampaignService.addContactsToAudience(campaignId, selectedContacts);
      fetchContacts();
      fetchAudience();
      setSelectedContacts([]);
      setSelectAll(false);
      toast.success("Contatos adicionados com sucesso!");
    } catch (error) {
      toast.error("Erro ao adicionar contatos à audiência.");
    }
  };

  const removeContactFromAudience = async (contactId: string) => {
    try {
      await CampaignService.removeContactFromAudience(campaignId, contactId);
      fetchContacts();
      fetchAudience();
      toast.success("Contato removido com sucesso!");
    } catch (error) {
      toast.error("Erro ao remover contato da audiência.");
    }
  };

  const removeAllContactsFromAudience = async () => {
    if (audiences.length === 0) return;
    try {
      await CampaignService.removeAllAudience(campaignId);
      fetchContacts();
      fetchAudience();
      toast.success("Todos os contatos foram removidos com sucesso!");
    } catch (error) {
      toast.error("Erro ao remover contatos da audiência.");
    }
  };

  // -------------------------------------
  // Handlers de seleção de contatos
  const handleSelectContact = (contactId: string) => {
    setSelectedContacts((prev) =>
      prev.includes(contactId) ? prev.filter((id) => id !== contactId) : [...prev, contactId]
    );
  };

  const handleSelectAll = () => {
    if (selectAll) {
      setSelectedContacts([]);
    } else {
      setSelectedContacts(contacts.map((c) => c.id).filter(Boolean) as string[]);
    }
    setSelectAll(!selectAll);
  };

  // -------------------------------------
  // Render principal
  //   if (loading) return <Spinner />;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Audiência da Campanha</h1>

      {/* 🔹 Filtros */}
      {isEditing && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Filtros {isFetching && <Spinner className="sm" />}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {/* Nome, Email e WhatsApp */}
              <Input
                disabled={!isEditing}
                name="name"
                value={filters.name}
                onChange={handleFilterChange}
                placeholder="Nome"
              />
              <Input
                disabled={!isEditing}
                name="email"
                value={filters.email}
                onChange={handleFilterChange}
                placeholder="Email"
              />
              <Input
                disabled={!isEditing}
                name="whatsapp"
                value={filters.whatsapp}
                onChange={handleFilterChange}
                placeholder="WhatsApp"
              />

              {/* Filtro de Gênero */}
              <Select
                disabled={!isEditing}
                value={filters.gender}
                onValueChange={(value) =>
                  setFilters((prev) => ({
                    ...prev,
                    gender: value === "none" ? "" : value,
                  }))
                }
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Selecionar Gênero" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="masculino">Masculino</SelectItem>
                    <SelectItem value="feminino">Feminino</SelectItem>
                    <SelectItem value="outro">Outro</SelectItem>
                  </SelectGroup>
                  <SelectSeparator />
                  <SelectItem className="bg-secondary flex justify-center" value="none">
                    <span className="text-xs">Limpar</span>
                  </SelectItem>
                </SelectContent>
              </Select>

              {/* Data de nascimento */}
              <Input
                disabled={!isEditing}
                type="date"
                name="birth_date_start"
                value={filters.birth_date_start}
                onChange={handleFilterChange}
              />
              <Input
                disabled={!isEditing}
                type="date"
                name="birth_date_end"
                value={filters.birth_date_end}
                onChange={handleFilterChange}
              />

              {/* Endereço */}
              <Input
                disabled={!isEditing}
                name="bairro"
                value={filters.bairro}
                onChange={handleFilterChange}
                placeholder="Bairro"
              />
              <Input
                disabled={!isEditing}
                name="cidade"
                value={filters.cidade}
                onChange={handleFilterChange}
                placeholder="Cidade"
              />
              <Input
                disabled={!isEditing}
                name="estado"
                value={filters.estado}
                onChange={handleFilterChange}
                placeholder="Estado"
              />

              {/* Tags */}
              <Input
                disabled={!isEditing}
                name="tags"
                value={filters.tags}
                onChange={(e) => setFilters((prev) => ({ ...prev, tags: e.target.value }))}
                onBlur={() => {
                  setFilters((prev) => ({
                    ...prev,
                    tags: prev.tags
                      .split(",")
                      .map((tag) => tag.trim())
                      .filter((tag) => tag.length > 0)
                      .join(", "),
                  }));
                }}
                placeholder="Tags (separadas por vírgula)"
              />
            </div>

            {/* Botão Adicionar Todos dentro do painel de Filtros */}
            <div className="flex justify-end mt-4">
              <Button
                onClick={addAllContactsToAudience}
                disabled={!isEditing || !contacts || contacts.length === 0}
                variant="default"
              >
                <PlusCircle className="w-5 h-5 mr-2" />
                Adicionar Todos os Filtrados
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* 🔹 Contatos Disponíveis */}
      {isEditing && (
        <Card className="mb-6 bg-contacts-available">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <UserCheck size={20} className="text-primary" /> Contatos Disponíveis{" "}
              {isFetching && <Spinner className="sm" />}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {contacts?.length > 0 ? (
              <>
                {/* Botões dentro do painel */}
                <div className="flex justify-between mb-4">
                  <Button onClick={handleSelectAll} variant="default" disabled={!isEditing}>
                    {selectAll ? (
                      <>
                        <CheckSquare className="w-5 h-5 mr-2" />
                        Desmarcar Todos
                      </>
                    ) : (
                      <>
                        <Square className="w-5 h-5 mr-2" />
                        Selecionar Todos
                      </>
                    )}
                  </Button>

                  <Button
                    onClick={addContactsToAudience}
                    disabled={!isEditing || selectedContacts.length === 0}
                    variant="default"
                  >
                    <UserPlus className="w-5 h-5 mr-2" />
                    Adicionar Selecionados
                  </Button>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                  {contacts.map((contact) => (
                    <div
                      key={contact.id}
                      className="flex items-center bg-muted p-4 rounded-lg shadow-md"
                    >
                      <Checkbox
                        disabled={!isEditing}
                        checked={contact.id ? selectedContacts.includes(contact.id) : false}
                        onCheckedChange={() => contact.id && handleSelectContact(contact.id)}
                        className="mr-2"
                      />
                      <span className="text-sm font-medium">{contact.name}</span>
                    </div>
                  ))}
                </div>

                {/* Paginação dos Contatos Disponíveis */}
                {contacts.length > 0 &&
                  pagination.available.total_records > pagination.available.per_page && (
                    <div className="p-4">
                      <Paginator
                        totalRecords={pagination.available.total_records}
                        totalPages={pagination.available.total_pages}
                        currentPage={pagination.available.current_page}
                        perPage={pagination.available.per_page}
                        onPageChange={handleAvailablePageChange}
                      />
                    </div>
                  )}
              </>
            ) : (
              <p className="text-gray-500 text-center py-4">
                Nenhum contato disponível para adicionar.
              </p>
            )}
          </CardContent>
        </Card>
      )}

      {/* 🔹 Contatos Adicionados */}
      <Card className="mb-6 bg-contacts-added">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <UserCheck size={20} className="text-primary" /> Contatos Adicionados
          </CardTitle>
        </CardHeader>

        <CardContent>
          {audiences?.length > 0 ? (
            <>
              <div className="flex justify-end mb-4">
                <Button
                  onClick={removeAllContactsFromAudience}
                  variant="destructive"
                  disabled={!isEditing}
                >
                  <Trash2 className="w-5 h-5 mr-2" />
                  Remover Todos
                </Button>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {audiences?.map((aud) => (
                  <div
                    key={aud.id}
                    className="flex justify-between items-center bg-muted p-4 rounded-lg shadow-md"
                  >
                    <span className="text-sm font-medium">{aud.name}</span>
                    <Button
                      disabled={!isEditing}
                      variant="ghost"
                      size="icon"
                      onClick={() => removeContactFromAudience(aud.id)}
                      className="text-red-600 hover:text-red-800 transition"
                    >
                      <Trash2 size={18} />
                    </Button>
                  </div>
                ))}
              </div>

              {/* Paginação dos Contatos Adicionados */}
              {audiences.length > 0 &&
                pagination.audience.total_records > pagination.audience.per_page && (
                  <div className="p-4">
                    <Paginator
                      totalRecords={pagination.audience.total_records}
                      totalPages={pagination.audience.total_pages}
                      currentPage={pagination.audience.current_page}
                      perPage={pagination.audience.per_page}
                      onPageChange={handleAudiencePageChange}
                    />
                  </div>
                )}
            </>
          ) : (
            <p className="text-gray-500 text-center py-4">
              Nenhum contato adicionado à campanha ainda.
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default CampaignAudience;
