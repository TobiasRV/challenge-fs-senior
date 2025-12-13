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
import { localStorageKeys } from "@/src/utils/consts";

type Inputs = {
  username: string;
  password: string;
  email: string;
  role: UserRolesEnum;
};

type CreateUserModalProps = {
  isOpen: boolean;
  handleClose: (success: boolean) => void;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function CreateUserModal({
  isOpen,
  handleClose,
}: CreateUserModalProps): ReactNode {
  const {
    control,
    register,
    handleSubmit,
    formState: { errors },
    setError,
    reset,
  } = useForm<Inputs>({
    reValidateMode: "onChange"
  });

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const {
    emailAlreadyExists,
    userLoading,
    error,
    statusCode,
    clearRequestState,
    createUser
  } = useUserStore(
    useShallow((state) => ({
      emailAlreadyExists: state.emailAlreadyExists,
      userLoading: state.loading,
      error: state.error,
      statusCode: state.statusCode,
      clearRequestState: state.clearRequestState,
      createUser: state.createUser,
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
    const teamId = localStorage.getItem(localStorageKeys.TEAM_ID) || "";
    const response = await createUser({
      ...data,
      teamId
    });

    if (response === HttpStatusCode.Created) {
        close(true)
    }
  };

  const close = (success: boolean) => {
    reset();
    handleClose(success);
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al crear usuario. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    [HttpStatusCode.Conflict]:
      "Ya existe un usuario con ese email",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.Unauthorized]: "Error al crear usuario.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Crear nuevo usuario</Dialog.Title>
        <Dialog.Content className="fixed inset-0 flex items-center justify-center p-4">
          <div className="w-full max-w-md rounded-md bg-white p-6 sm:p-8 text-gray-900 shadow max-h-[calc(100vh-16px)] sm:max-h-[calc(100vh-32px)] overflow-auto">
          <h2 className="text-xl">Nuevo Usuario</h2>

          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-6">
              <Label>Usuario</Label>
              <Input
                className="mt-2"
                {...register("username", {
                  required: { value: true, message: "El usuario es requerido" },
                })}
              />
              {errors.username ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.username.message}
                </p>
              ) : null}
            </div>
            <div className="mt-6">
              <Label>Contraseña</Label>
              <Input
                type="password"
                className="mt-2"
                {...register("password", {
                  required: {
                    value: true,
                    message: "La contraseña es requerida",
                  },
                })}
              />
              {errors.password ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.password.message}
                </p>
              ) : null}
            </div>
            <div className="mt-6">
              <Label>Email</Label>
              <Input
                className="mt-2"
                {...register("email", {
                  required: "El email es requerido",
                  validate: {
                    isValidEmail: (email) => {
                      return validateEmail(email) || "Email invalido";
                    },
                    emailAlreadyExists: async (email) => {
                      const exists = await emailAlreadyExists(email);
                      return (!exists) || "Ya existe una cuenta con ese email";
                    },
                  },
                })}
              />
              {errors.email ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.email.message}
                </p>
              ) : null}
            </div>

            <div className="mt-6">
              <div>
                <Label>Rol</Label>
              </div>

              <Controller
                name="role"
                control={control}
                rules={{
                  required: { value: true, message: "El rol es requerido" },
                }}
                render={({ field }) => (
                  <Select.Root value={field.value} onValueChange={field.onChange}>
                    <Select.Trigger
                      className="
                        inline-flex w-full h-10 items-center justify-between
                        rounded-lg
                        border border-gray-300
                        bg-white
                        px-3
                        text-sm text-gray-800
                        shadow-sm
                        outline-none
                        transition
                        hover:border-gray-400
                        focus:border-gray-400
                        focus:ring-2 focus:ring-gray-300/40
                        data-[placeholder]:text-gray-400
                      "
                    >
                      <Select.Value placeholder="Seleccione un rol" />
                      <Select.Icon className="ml-2 text-gray-400">
                        <ChevronDown className="h-4 w-4" />
                      </Select.Icon>
                    </Select.Trigger>
                    <Select.Portal>
                      <Select.Content
                        position="popper"
                        className="
                          z-50
                          min-w-[var(--radix-select-trigger-width)]
                          overflow-hidden
                          rounded-lg
                          border border-gray-200
                          bg-white
                          shadow-md
                        "
                      >
                        <Select.Viewport className="p-1">
                          <Select.Group>
                            {[UserRolesEnum.MANAGER, UserRolesEnum.MEMBER].map((role, idx) => (
                              <Select.Item
                                value={role}
                                key={idx}
                                className="
                                  relative flex cursor-pointer select-none items-center
                                  rounded-md
                                  px-3 py-2
                                  text-sm text-gray-700
                                  outline-none
                                  transition
                                  focus:bg-gray-100
                                  data-[highlighted]:bg-gray-200
                                  data-[state=checked]:bg-gray-100
                                "
                              >
                                <Select.ItemText>
                                  {UserRolesTranslation[role as UserRolesEnum]}
                                </Select.ItemText>
                              </Select.Item>
                            ))}
                          </Select.Group>
                        </Select.Viewport>
                      </Select.Content>
                    </Select.Portal>
                  </Select.Root>
                )}
              />
            </div>

            {error && statusCode && (
              <Alert className="mt-5 w-100" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

              <div className="flex end justify-end gap-3 mt-10">
              <Button type="button" variant="ghost" onClick={() => close(false)}>Cancelar</Button>
              <Button
                disabled={userLoading || Object.keys(errors).length !== 0}
                type="submit"
              >
                Crear Usuario
              </Button>
            </div>
          </form>
          </div>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
