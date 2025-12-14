"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Input } from "../../ui/input";
import { useProjectStore } from "@/src/stores/projects/projects";


type Inputs = {
  name: string;
};

type CreateProjectModalProps = {
    isOpen: boolean;
    handleClose: (success: boolean) => void
}


type MapErrorMessage = {
  [key: number]: string;
};

export default function CreateProjectModal({isOpen, handleClose}: CreateProjectModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset
  } = useForm<Inputs>();

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { createProject, error, statusCode, loading, clearRequestState } = useProjectStore(
    useShallow((state) => ({
        createProject: state.createProject,
        error: state.error,
        statusCode: state.statusCode,
        loading: state.loading,
        clearRequestState: state.clearRequestState
    }))
  )

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
    const response = await createProject(data);
    if (response === HttpStatusCode.Created) {
        close(true)
    }
  };

  const close = (success: boolean) => {
    reset()
    handleClose(success)
  }

const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
        "Error al crear proyecto. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.Forbidden]: "Error al crear proyecto.",
    [HttpStatusCode.Unauthorized]: "Error al crear proyecto.",
};

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Crear nuevo proyecto</Dialog.Title>
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100%-2rem)] max-w-md rounded-md bg-white p-4 sm:p-6 md:p-8 text-gray-900 shadow max-h-[90vh] overflow-y-auto">
          <h2 className="text-lg sm:text-xl">Crear nuevo proyecto</h2>

          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-4 sm:mt-5">
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
              <Alert className="mt-4 sm:mt-5" variant="error">
                <AlertTitle className="font-normal text-sm sm:text-base">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

              <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 sm:gap-3 mt-4 sm:mt-5">
                <Button variant="ghost" onClick={() => close(false)} className="w-full sm:w-auto">Cancelar</Button>
                <Button type="submit" className="w-full sm:w-auto">Confirmar</Button>
              </div>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
