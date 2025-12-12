"use client";

import { useTeamStore } from "@/src/stores/teams/teams";
import * as Dialog from "@radix-ui/react-dialog";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../ui/alert";
import { Button } from "../ui/button";
import { Input } from "../ui/input";


type Inputs = {
  name: string;
};

type CreateTeamModalProps = {
    isOpen: boolean;
    handleClose: () => void
}


type MapErrorMessage = {
  [key: number]: string;
};

export default function CreateTeamModal({isOpen, handleClose}: CreateTeamModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset
  } = useForm<Inputs>();

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { createTeam, error, statusCode, loading, clearRequestState } = useTeamStore(useShallow((state) => ({
    createTeam: state.createTeam,
    error: state.error,
    statusCode: state.statusCode,
    loading: state.loading,
    clearRequestState: state.clearRequestState,
  })))

  useEffect(() => {
    if (error) {
      clearTimeout(alertTimeout);
      const timeoutAux = setTimeout(
        () => clearRequestState(),
        7000,
      );
      setAlertTimeout(timeoutAux);
    }
  }, [error, statusCode]);

  const onSubmit = async (data: Inputs) => {
    const response = await createTeam(data);
    if (response === HttpStatusCode.Created) {
        close()
    }
  };

  const close = () => {
    reset()
    handleClose()
  }

const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
        "Error al crear equipo. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    [HttpStatusCode.Conflict]: "El usuario ya tiene un equipo creado con ese nombre.",
    // If the endpoint responde unauthorized it should have been redirected to login so it uses a generic message
    [HttpStatusCode.Unauthorized]: "Error al crear equipo.",
    // If the endpoint responde not found the user was deleted but the token is valid so it uses a generic message
    [HttpStatusCode.NotFound]: "Error al crear equipo."
};

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Crear nuevo equipo</Dialog.Title>
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 max-w-md min-w-md rounded-md bg-white p-8 text-gray-900 shadow">
          <h2 className="text-xl">Ponle nombre a tu equipo!</h2>

          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-5">
              <Label>Nombre</Label>
              <Input
                className="mt-2"
                {...register("name", {
                  required: { value: true, message: "El nombre es requerido" },
                })}
              />
              {errors.name ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.name.message}
                </p>
              ) : null}

              {(error && statusCode) && (
              <Alert className="mt-5" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

              <div className="flex end justify-end gap-3 mt-5">
                <Button variant="ghost" onClick={close}>Cancelar</Button>
                <Button type="submit">Confirmar</Button>
              </div>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
