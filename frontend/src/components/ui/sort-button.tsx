import { Button } from "@/components/ui/button";
import { ArrowUpDown, ArrowUp, ArrowDown } from "lucide-react";

type SortDirection = "asc" | "desc" | null;

interface SortButtonProps {
  label: string;
  active: boolean;
  direction: SortDirection;
  onClick: () => void;
}

export function SortButton({
  label,
  active,
  direction,
  onClick,
}: SortButtonProps) {
  return (
    <Button
      variant="ghost"
      onClick={onClick}
      className={`h-8 ${active ? "text-primary" : "text-muted-foreground"}`}
    >
      {label}
      {direction === null && <ArrowUpDown className="ml-2 h-4 w-4" />}
      {direction === "asc" && <ArrowUp className="ml-2 h-4 w-4" />}
      {direction === "desc" && <ArrowDown className="ml-2 h-4 w-4" />}
    </Button>
  );
}
