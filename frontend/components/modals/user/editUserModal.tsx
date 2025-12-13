"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Input } from "../../ui/input";
import { UserRolesEnum } from "@/src/utils/enums";
import { validateEmail } from "@/src/utils/validations";
import { useUserStore } from "@/src/stores/users/users";
import * as Select from "@radix-ui/react-select";
import { ChevronDown } from "lucide-react";
import { UserRolesTranslation } from "@/src/utils/translations";
import { IUser } from "@/src/stores/users/users.interface";

type Inputs = {
  username: string;
  email: string;
};

type EditUserModalProps = {
  isOpen: boolean;
  user: IUser;
  handleClose: (success: boolean) => void;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function EditUserModal({
  isOpen,
  handleClose,
  user,
}: EditUserModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
  } = useForm<Inputs>({
    reValidateMode: "onChange",
    defaultValues: {
      username: user.username,
      email: user.email,
    },
  });

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const {
    emailAlreadyExists,
    userLoading,
    error,
    statusCode,
    clearRequestState,
    updateUser,
  } = useUserStore(
    useShallow((state) => ({
      emailAlreadyExists: state.emailAlreadyExists,
      userLoading: state.loading,
      error: state.error,
      statusCode: state.statusCode,
      clearRequestState: state.clearRequestState,
      updateUser: state.updateUser,
    }))
  );

  useEffect(() => {
    if (error) {
      clearTimeout(alertTimeout);
      const timeoutAux = setTimeout(() => clearRequestState(), 7000);
      setAlertTimeout(timeoutAux);
    }
  }, [error, statusCode]);

  const onSubmit = async (data: Inputs) => {
    const response = await updateUser({
      username: data.username,
      email: data.email,
      id: user.id,
    });

    if (response === HttpStatusCode.Created) {
      close(true);
    }
  };

  const close = (success: boolean) => {
    reset();
    handleClose(success);
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al actualizar el usuario. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    [HttpStatusCode.Conflict]: "Ya existe un usuario con ese email",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.NotFound]: "Usuario no encontrado",
    [HttpStatusCode.Unauthorized]: "Error al actualizar el usuario.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Editar usuario</Dialog.Title>
        <Dialog.Content className="fixed inset-0 flex items-center justify-center p-4">
          <div className="w-full max-w-md rounded-md bg-white p-6 sm:p-8 text-gray-900 shadow max-h-[calc(100vh-16px)] sm:max-h-[calc(100vh-32px)] overflow-auto">
            <h2 className="text-xl">Editar usuario</h2>

            <form onSubmit={handleSubmit(onSubmit)}>
              <div className="mt-6">
                <Label>Usuario</Label>
                <Input
                  className="mt-2"
                  {...register("username", {
                    validate: (username) => {
                      if (username === user.username) return true;
                      if (!username || username.trim() === "")
                        return "El usuario es requerido";
                      return true;
                    },
                  })}
                />
                {errors.username ? (
                  <p className="text-red-500 text-xs mt-2">
                    {errors.username.message}
                  </p>
                ) : null}
              </div>
              <div className="mt-6">
                <Label>Email</Label>
                <Input
                  className="mt-2"
                  {...register("email", {
                    validate: async (email) => {
                      if (email === user.email) return true;
                      if (!email) return "El email es requerido";
                      if (!validateEmail(email)) return "Email invalido";
                      const exists = await emailAlreadyExists(email);
                      return !exists || "Ya existe una cuenta con ese email";
                    },
                  })}
                />
                {errors.email ? (
                  <p className="text-red-500 text-xs mt-2">
                    {errors.email.message}
                  </p>
                ) : null}
              </div>

              {error && statusCode && (
                <Alert className="mt-5 w-100" variant="error">
                  <AlertTitle className="font-normal">
                    {mapErrorMessage[statusCode]}
                  </AlertTitle>
                </Alert>
              )}

              <div className="flex end justify-end gap-3 mt-10">
                <Button type="button" variant="ghost" onClick={() => close(false)}>
                  Cancelar
                </Button>
                <Button
                  disabled={userLoading || Object.keys(errors).length !== 0 || !isDirty}
                  type="submit"
                >
                  Editar Usuario
                </Button>
              </div>
            </form>
          </div>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
