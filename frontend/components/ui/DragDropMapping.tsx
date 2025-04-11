// File: components/ui/DragDropMapping.tsx

import { useDrag, useDrop } from "react-dnd";
import { useRef } from "react";
import { X } from "lucide-react";

const ItemType = {
  CSV_FIELD: "csv_field",
};

// Componente de um item arrastável (cabeçalho do CSV)
function DraggableCSVField({ name }: { name: string }) {
  const ref = useRef<HTMLDivElement>(null);
  const [{ isDragging }, drag] = useDrag(() => ({
    type: ItemType.CSV_FIELD,
    item: { name },
    collect: (monitor) => ({
      isDragging: !!monitor.isDragging(),
    }),
  }));

  drag(ref);

  return (
    <div
      ref={ref}
      className={`px-3 py-1 text-sm border rounded-full cursor-grab bg-blue-500 text-white ${
        isDragging ? "opacity-50" : ""
      }`}
    >
      {name}
    </div>
  );
}

// Componente onde os campos do CSV podem ser soltos (cada campo do banco)
function DropTargetField({
  field,
  assignedCSVFields,
  onDrop,
  onRemove,
  onRuleChange,
  ruleText,
}: {
  field: string;
  assignedCSVFields: string[];
  onDrop: (csvField: string) => void;
  onRemove: (csvField: string) => void;
  onRuleChange: (field: string, value: string) => void;
  ruleText: string;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const [{ isOver }, drop] = useDrop(() => ({
    accept: ItemType.CSV_FIELD,
    drop: (item: { name: string }) => onDrop(item.name),
    collect: (monitor) => ({
      isOver: !!monitor.isOver(),
    }),
  }));

  drop(ref);

  return (
    <div
      ref={ref}
      className={`p-4 border rounded bg-white shadow-sm min-h-[80px] ${
        isOver ? "border-green-500 bg-green-100" : ""
      }`}
    >
      <h4 className="font-semibold text-gray-700">{field}</h4>

      {/* Lista de tags associadas */}
      <div className="flex flex-wrap gap-2 mt-2">
        {assignedCSVFields.length > 0 ? (
          assignedCSVFields.map((csvField) => (
            <div
              key={csvField}
              className="flex items-center px-3 py-1 text-sm border rounded-full bg-green-500 text-white"
            >
              {csvField}
              <button onClick={() => onRemove(csvField)} className="ml-2">
                <X size={14} />
              </button>
            </div>
          ))
        ) : (
          <p className="text-gray-400 text-sm">Arraste um campo CSV aqui</p>
        )}
      </div>

      {/* Campo de regras */}
      <label className="block mt-2 text-sm font-medium text-gray-600">Regras</label>
      <textarea
        className="w-full p-2 border rounded bg-white text-sm"
        rows={2}
        placeholder="Descreva como usar esse campo..."
        value={ruleText || ""}
        onChange={(e) => onRuleChange(field, e.target.value)}
      />
    </div>
  );
}

// Componente principal do mapeamento
export function DragDropMapping({
  csvHeaders,
  config,
  onUpdateConfig,
}: {
  csvHeaders: string[];
  config: Record<string, { sources: string[]; rules: string }>;
  onUpdateConfig: (newConfig: Record<string, { sources: string[]; rules: string }>) => void;
}) {
  const handleDrop = (field: string, csvField: string) => {
    const sources = config[field]?.sources ?? [];
    const updatedSources = [...new Set([...sources, csvField])];

    onUpdateConfig({
      ...config,
      [field]: { ...config[field], sources: updatedSources },
    });
  };

  const handleRemove = (field: string, csvField: string) => {
    const sources = config[field]?.sources ?? [];
    const updatedSources = sources.filter((s) => s !== csvField);

    onUpdateConfig({
      ...config,
      [field]: { ...config[field], sources: updatedSources },
    });
  };

  const handleRuleChange = (field: string, value: string) => {
    onUpdateConfig({
      ...config,
      [field]: { ...config[field], rules: value },
    });
  };

  return (
    <div className="flex gap-6">
      {/* Lista de campos do CSV para arrastar */}
      <div className="w-1/3 border p-4 bg-gray-50 rounded">
        <h3 className="text-lg font-semibold">Campos do CSV</h3>
        <div className="space-y-2">
          {csvHeaders.map((header) => (
            <DraggableCSVField key={header} name={header} />
          ))}
        </div>
      </div>

      {/* Área onde os campos do banco estão */}
      <div className="w-2/3 border p-4 bg-gray-50 rounded">
        <h3 className="text-lg font-semibold">Mapeamento de Campos</h3>
        <div className="grid grid-cols-2 gap-4">
          {Object.keys(config).map((field) => (
            <DropTargetField
              key={field}
              field={field}
              assignedCSVFields={config[field]?.sources ?? []}
              onDrop={(csvField) => handleDrop(field, csvField)}
              onRemove={(csvField) => handleRemove(field, csvField)}
              onRuleChange={handleRuleChange}
              ruleText={config[field]?.rules ?? ""}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
