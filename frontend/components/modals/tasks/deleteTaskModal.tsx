"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { useTaskStore } from "@/src/stores/tasks/tasks";
import { ITask } from "@/src/stores/tasks/tasks.interface";

type DeleteTaskModalProps = {
  isOpen: boolean;
  task: ITask;
  handleClose: (success: boolean) => void;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function DeleteTaskModal({
  isOpen,
  handleClose,
  task,
}: DeleteTaskModalProps): ReactNode {
  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { loading, error, statusCode, clearRequestState, deleteTask } =
    useTaskStore(
      useShallow((state) => ({
        loading: state.loading,
        error: state.error,
        statusCode: state.statusCode,
        clearRequestState: state.clearRequestState,
        deleteTask: state.deleteTask,
      }))
    );

  useEffect(() => {
    if (error) {
      clearTimeout(alertTimeout);
      const timeoutAux = setTimeout(() => clearRequestState(), 7000);
      setAlertTimeout(timeoutAux);
    }
  }, [error, statusCode]);

  const onSubmit = async () => {
    const response = await deleteTask(task.id);

    if (response === HttpStatusCode.Ok) {
      close(true);
    }
  };

  const close = (success: boolean) => {
    handleClose(success);
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al eliminar la tarea. Por favor intente nuevamente.",
    [HttpStatusCode.Unauthorized]: "Error al eliminar la tarea.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Eliminar tarea</Dialog.Title>
        <Dialog.Content className="fixed inset-0 flex items-center justify-center p-4">
          <div className="w-full max-w-md rounded-md bg-white p-6 sm:p-8 text-gray-900 shadow max-h-[calc(100vh-16px)] sm:max-h-[calc(100vh-32px)] overflow-auto">
            <h2 className="text-xl">Eliminar tarea</h2>

            <Alert className="mt-5 w-full" variant="warning">
              <AlertTitle className="font-normal">
                ¿Está seguro que desea eliminar la tarea "{task.title}"?
              </AlertTitle>
            </Alert>

            {error && statusCode && (
              <Alert className="mt-5 w-full" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[statusCode]}
                </AlertTitle>
              </Alert>
            )}

            <div className="flex end justify-end gap-3 mt-10">
              <Button
                type="button"
                variant="ghost"
                onClick={() => close(false)}
              >
                Cancelar
              </Button>
              <Button disabled={loading} onClick={() => onSubmit()}>
                Eliminar Tarea
              </Button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
