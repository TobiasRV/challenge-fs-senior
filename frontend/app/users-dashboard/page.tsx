"use client";

import CreateUserModal from "@/components/modals/user/createUserModal";
import DeleteUserModal from "@/components/modals/user/deleteUserModal";
import EditUserModal from "@/components/modals/user/editUserModal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import * as Table from "@/components/ui/table";
import { useAuthStore } from "@/src/stores/auth/auth";
import { useUserStore } from "@/src/stores/users/users";
import { IGetUsersParams, IUser } from "@/src/stores/users/users.interface";
import { localStorageKeys } from "@/src/utils/consts";
import { UserRolesEnum } from "@/src/utils/enums";
import useDebounced from "@/src/utils/hooks/debounce";
import { getLsItem } from "@/src/utils/localStorage";
import { PencilLine, TrashIcon } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useShallow } from "zustand/shallow";

export default function UserDashboard() {
  const { user } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
    })),
  );

  const { users, getUsers, usersLoading } = useUserStore(
    useShallow((state) => ({
      users: state.users,
      getUsers: state.getUsers,
      usersLoading: state.loading,
    })),
  );

  const teamId = getLsItem(localStorageKeys.TEAM_ID);

  const [filter, setFilter] = useState<IGetUsersParams>({
    limit: 10,
    cursor: "",
    teamId,
  });

  const [showCreateUserModal, setShowCreateUserModal] = useState<boolean>(false)
  const [editUserData, setEditUserData] = useState<IUser | undefined>(undefined)
  const [deleteUserData, setDeleteUserData] = useState<IUser | undefined>(undefined)

  useEffect(() => {
    getUsers(filter);
  }, [filter]);

  const getNext = async () => {
    setFilter((prev) => ({
      ...prev,
      cursor: users.next,
    }));
  };

  const getPrev = async () => {
    setFilter((prev) => ({
      ...prev,
      cursor: users.prev,
    }));
  };


  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilter((prev) => ({
      ...prev,
      email: e.target.value,
    }));
  };

  const debounce = useDebounced(handleSearch, 1000);
  const onSearchChange = useMemo(() => debounce, []);

  const handleUserCreate = (success: boolean) => {
    if (success) {
      getUsers(filter);
    }
    
    setShowCreateUserModal(false);
  }

  const handleUserEdit = (success: boolean) => {

    if (success) {
      getUsers(filter);
    }
    setEditUserData(undefined);
  }

   const handleUserDelete = (success: boolean) => {

    if (success) {
      getUsers(filter);
    }
    setDeleteUserData(undefined);
  }

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="p-5">
        <h1 className="py-5 text-xl font-bold">Usuarios</h1>

       <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 px-4 sm:px-10">
          <Input className="w-full sm:w-1/2" onChange={onSearchChange} placeholder="Busqueda" />
          <Button className="w-full sm:w-auto" onClick={() => setShowCreateUserModal(true)}>
            Crear Usuario
          </Button>
        </div>

        <div className="py-5 px-4 sm:px-10">
          {/* Card view for mobile */}
          <div className="sm:hidden grid gap-3">
            {users.data.map((user) => (
              <div key={user.id} className="p-3 bg-white rounded-md shadow flex items-start justify-between">
                <div className="min-w-0">
                  <div className="text-sm font-semibold truncate">{user.username}</div>
                  <div className="text-xs text-muted-foreground truncate">{user.email}</div>
                  <div className="text-xs text-gray-500">{user.role}</div>
                </div>
                <div className="ml-3 flex-shrink-0">
                  <Button size="sm" variant="ghost" onClick={() => setEditUserData(user)}>
                    <PencilLine />
                  </Button>
                  <Button size="sm" variant="ghost" onClick={() => setDeleteUserData(user)}>
                    <TrashIcon />
                  </Button>
                </div>
              </div>
            ))}
          </div>

          <div className="hidden sm:block overflow-x-auto rounded-md bg-white shadow-sm">
          <Table.Table >
            <Table.TableHeader>
              <Table.TableRow>
                <Table.TableHead>Username</Table.TableHead>
                <Table.TableHead>Email</Table.TableHead>
                <Table.TableHead>Role</Table.TableHead>
                <Table.TableHead>Acciones</Table.TableHead>
              </Table.TableRow>
            </Table.TableHeader>
            <Table.TableBody>
              {users.data.map((user, key) => (
                <Table.TableRow key={key}>
                  <Table.TableCell>{user.username}</Table.TableCell>
                  <Table.TableCell>{user.email}</Table.TableCell>
                  <Table.TableCell>{user.role}</Table.TableCell>
                  <Table.TableCell>
                    <Button
                      variant={"ghost"}
                      onClick={() => setEditUserData(user)}
                    >
                      <PencilLine />
                    </Button>
                    <Button
                      variant={"ghost"}
                      onClick={() => setDeleteUserData(user)}
                    >
                      <TrashIcon />
                    </Button>
                  </Table.TableCell>
                </Table.TableRow>
              ))}
            </Table.TableBody>
          </Table.Table>
          </div>
        </div>
        <div className="flex justify-center sm:justify-end pr-4 sm:pr-10 gap-4">
          <Button disabled={!users.prev || usersLoading} onClick={getPrev}>
            Anterior
          </Button>
          <Button disabled={!users.next || usersLoading} onClick={getNext}>
            Siguiente
          </Button>
        </div>
      </div>
      <CreateUserModal isOpen={showCreateUserModal} handleClose={handleUserCreate} />
      {editUserData && <EditUserModal isOpen={!!editUserData} handleClose={handleUserEdit} user={editUserData} />}
      {deleteUserData && <DeleteUserModal isOpen={!!deleteUserData} handleClose={handleUserDelete} user={deleteUserData}/>}
    </main>
  );
}
