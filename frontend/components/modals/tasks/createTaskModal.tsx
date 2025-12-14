"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { useForm, Controller } from "react-hook-form";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Input } from "../../ui/input";
import { useProjectStore } from "@/src/stores/projects/projects";
import { useTaskStore } from "@/src/stores/tasks/tasks";
import MemberSelector from "@/components/selectors/memberSelectors";
import { Textarea } from "@/components/ui/textarea";

type Inputs = {
  title: string;
  description?: string;
  userId?: string;
};

type CreateTaskModalProps = {
  isOpen: boolean;
  handleClose: (success: boolean) => void;
  teamId: string;
  projectId: string;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function CreateTaskModal({
  isOpen,
  handleClose,
  teamId,
  projectId,
}: CreateTaskModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<Inputs>();
  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { createTask, error, statusCode, loading, clearRequestState } =
    useTaskStore(
      useShallow((state) => ({
        createTask: state.createTask,
        error: state.error,
        statusCode: state.statusCode,
        loading: state.loading,
        clearRequestState: state.clearRequestState,
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
    const response = await createTask({
      ...data,
      projectId,
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
      "Error al crear tarea. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.NotFound]: "Proyecto no valido",
    [HttpStatusCode.Forbidden]: "Error al crear tarea.",
    [HttpStatusCode.Unauthorized]: "Error al crear tarea.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Crear nuevaTarea</Dialog.Title>
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100%-2rem)] max-w-md rounded-md bg-white p-4 sm:p-6 md:p-8 text-gray-900 shadow max-h-[90vh] overflow-y-auto">
          <h2 className="text-lg sm:text-xl">Crear nueva tarea</h2>

          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-5">
              <Label>Titulo</Label>
              <Input
                className="mt-2"
                {...register("title", {
                  required: { value: true, message: "El nombre es requerido" },
                })}
              />
              {errors.title ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.title.message}
                </p>
              ) : null}
            </div>

            <div className="mt-5">
              <Label>Descripcion</Label>
              <Textarea
                className="mt-2"
                {...register("description")}
              />
              {errors.description ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.description.message}
                </p>
              ) : null}
            </div>

            <div className="mt-5">
              <Label>Responsable</Label>
              <div className="mt-2">
                <Controller
                  name="userId"
                  control={control}
                  render={({ field }) => (
                    <MemberSelector
                      limit={2}
                      teamId={teamId}
                      value={field.value}
                      onChange={field.onChange}
                    />
                  )}
                />
              </div>
            </div>

            {error && statusCode && (
              <Alert className="mt-5" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

            <div className="flex flex-col-reverse sm:flex-row sm:justify-end gap-3 mt-5">
              <Button type="button" variant="ghost" onClick={() => close(false)} className="w-full sm:w-auto">
                Cancelar
              </Button>
              <Button type="submit" className="w-full sm:w-auto">Crear tarea</Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
