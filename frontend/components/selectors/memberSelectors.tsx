"use client";

import React, { useEffect } from "react";
import { IUser } from "@/src/stores/users/users.interface";
import { getUsers as getUsersService } from "@/src/services/apiServices/users/service";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
  SelectGroup,
} from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import { UserRolesEnum } from "@/src/utils/enums";

interface MemberSelectorProps {
  value?: string;
  onChange?: (id?: string) => void;
  placeholder?: string;
  teamId?: string;
  limit?: number;
}

export default function MemberSelector({
  value,
  onChange,
  placeholder = "Seleccionar usuario",
  teamId = "",
  limit = 10,
}: MemberSelectorProps) {
  const [localUsers, setLocalUsers] = React.useState<IUser[]>([]);
  const [localNext, setLocalNext] = React.useState<string>("");
  const [localLoading, setLocalLoading] = React.useState(false);

  const fetchUsers = React.useCallback(
    async (cursor = "") => {
      try {
        setLocalLoading(true);
        const resp = await getUsersService({ teamId, limit, cursor, role: UserRolesEnum.MEMBER });
        if (resp.error) {
          // keep existing users on error
          setLocalLoading(false);
          return;
        }

        const incoming: IUser[] = resp.data || [];
        setLocalUsers((prev) => (cursor ? [...prev, ...incoming] : incoming));
        setLocalNext(resp.pagination?.next_cursor || "");
        setLocalLoading(false);
      } catch (e) {
        setLocalLoading(false);
      }
    },
    [teamId, limit]
  );

  useEffect(() => {
    fetchUsers("");
  }, [teamId, fetchUsers]);

  // track previous value so we can restore selection after using the "load more" item
  const prevValueRef = React.useRef<string | undefined>(value);
  useEffect(() => {
    prevValueRef.current = value;
  }, [value]);

  const handleValueChange = (v: string) => {
    if (onChange) onChange(v);
  };

  return (
    <Select value={value} onValueChange={handleValueChange}>
      <SelectTrigger
        className="
                        inline-flex w-full h-16 items-center justify-start
                        rounded-lg
                        border border-gray-300
                        bg-white
                        px-3
                        py-7
                        text-sm text-gray-800
                        shadow-sm
                        outline-none
                        transition
                        hover:border-gray-400
                        focus:border-gray-400
                        focus:ring-2 focus:ring-gray-300/40
                        data-[placeholder]:text-gray-400
                        [&>svg]:ml-auto
                      "
      >
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {localUsers.length === 0 && (
            <SelectItem value="__no-users" disabled>
              {localLoading ? "Cargando..." : "Sin usuarios"}
            </SelectItem>
          )}
          {localUsers.map((u) => (
            <SelectItem key={u.id} value={u.id} className="
                relative flex cursor-pointer select-none items-center
                rounded-md
                px-3 py-2
                text-sm text-gray-700
                outline-none
                transition
                focus:bg-gray-100
                data-[highlighted]:bg-gray-200
                data-[state=checked]:bg-gray-100
            ">
              <div className="flex flex-col items-start py-2">
                <span className="font-medium text-base leading-5">{u.username}</span>
                <span className="text-sm text-muted-foreground truncate leading-4">{u.email}</span>
              </div>
            </SelectItem>
          ))}
          {localNext && (
            <SelectItem
              value="__load-more"
              onMouseDown={(e) => {
                e.preventDefault();
                if (!localLoading) fetchUsers(localNext);
              }}
            >
              <div className="flex w-full items-center py-2 text-sm">
                {localLoading ? (
                  <>
                    <Loader2 className="size-4 animate-spin mr-2" />
                    Cargando...
                  </>
                ) : (
                  "Cargar más"
                )}
              </div>
            </SelectItem>
          )}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
