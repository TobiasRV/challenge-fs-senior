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
import { useTaskStore } from "@/src/stores/tasks/tasks";
import MemberSelector from "@/components/selectors/memberSelectors";
import { Textarea } from "@/components/ui/textarea";
import { ITask } from "@/src/stores/tasks/tasks.interface";
import { TaskStatusEnum } from "@/src/utils/enums";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { TasksStatusTranslation } from "@/src/utils/translations";

type Inputs = {
  title: string;
  description?: string;
  userId: string;
  status: TaskStatusEnum;
};

type EditTaskModalProps = {
  isOpen: boolean;
  handleClose: (success: boolean) => void;
  teamId: string;
  task: ITask;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function EditTaskModal({
  isOpen,
  handleClose,
  teamId,
  task,
}: EditTaskModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<Inputs>({
    defaultValues: {
      title: task.title,
      description: task.description,
      userId: task.userId,
      status: task.status,
    },
  });
  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { updateTask, error, statusCode, loading, clearRequestState } =
    useTaskStore(
      useShallow((state) => ({
        updateTask: state.updateTask,
        error: state.error,
        statusCode: state.statusCode,
        loading: state.loading,
        clearRequestState: state.clearRequestState,
      }))
    );

  useEffect(() => {
    reset({
      title: task.title,
      description: task.description,
      userId: task.userId,
      status: task.status,
    });
  }, [task, reset]);

  useEffect(() => {
    if (error) {
      clearTimeout(alertTimeout);
      const timeoutAux = setTimeout(() => clearRequestState(), 7000);
      setAlertTimeout(timeoutAux);
    }
  }, [error, statusCode]);

  const onSubmit = async (data: Inputs) => {
    console.log({ task })
    const response = await updateTask({
      ...data,
      id: task.id,
    });
    if (response === HttpStatusCode.Ok) {
      close(true);
    }
  };

  const close = (success: boolean) => {
    reset();
    handleClose(success);
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al actualizar tarea. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    [HttpStatusCode.NotFound]: "Tarea no encontrada.",
    [HttpStatusCode.Forbidden]: "Error al actualizar tarea.",
    [HttpStatusCode.Unauthorized]: "Error al actualizar tarea.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Editar Tarea</Dialog.Title>
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100%-2rem)] max-w-md rounded-md bg-white p-4 sm:p-6 md:p-8 text-gray-900 shadow max-h-[90vh] overflow-y-auto">
          <h2 className="text-lg sm:text-xl">Editar tarea</h2>

          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-5">
              <Label>Titulo</Label>
              <Input
                className="mt-2"
                {...register("title", {
                  required: { value: true, message: "El titulo es requerido" },
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
                {errors.userId ? (
                  <p className="text-red-500 text-xs mt-2">
                    {errors.userId.message}
                  </p>
                ) : null}
              </div>
            </div>

            <div className="mt-5">
              <Label>Estado</Label>
              <div className="mt-2">
                <Controller
                  name="status"
                  control={control}
                  rules={{ required: { value: true, message: "El estado es requerido" } }}
                  render={({ field }) => (
                    <Select value={field.value} onValueChange={field.onChange}>
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Seleccionar estado" />
                      </SelectTrigger>
                      <SelectContent>
                        {[TaskStatusEnum.TODO, TaskStatusEnum.IN_PROGRESS, TaskStatusEnum.DONE].map((status, key) => (
                          <SelectItem key={key} value={status}>
                            {TasksStatusTranslation[status]}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.status ? (
                  <p className="text-red-500 text-xs mt-2">
                    {errors.status.message}
                  </p>
                ) : null}
              </div>
            </div>

            {error && statusCode && (
              <Alert className="mt-5" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

            <div className="flex end justify-end gap-3 mt-5">
              <Button type="button" variant="ghost" onClick={() => close(false)}>
                Cancelar
              </Button>
              <Button type="submit" disabled={loading}>
                Guardar cambios
              </Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
