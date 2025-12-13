"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { useProjectStore } from "@/src/stores/projects/projects";
import { IProject } from "@/src/stores/projects/projects.interface";

type DeleteProjectModalProps = {
  isOpen: boolean;
  project: IProject;
  handleClose: (success: boolean) => void;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function DeleteProjectModal({
  isOpen,
  handleClose,
  project,
}: DeleteProjectModalProps): ReactNode {
  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { loading, error, statusCode, clearRequestState, deleteProject } =
    useProjectStore(
      useShallow((state) => ({
        loading: state.loading,
        error: state.error,
        statusCode: state.statusCode,
        clearRequestState: state.clearRequestState,
        deleteProject: state.deleteProject,
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
    const response = await deleteProject(project.id);

    if (response === HttpStatusCode.Ok) {
      close(true);
    }
  };

  const close = (success: boolean) => {
    handleClose(success);
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al eliminar el usuario. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "El usuario no existe o ya esta eliminado",
    [HttpStatusCode.NotFound]: "El usuario no existe o ya esta eliminado",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.Unauthorized]: "Error al eliminar el usuario.",
  };

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Eliminar projecto</Dialog.Title>
        <Dialog.Content className="fixed inset-0 flex items-center justify-center p-4">
          <div className="w-full max-w-md rounded-md bg-white p-6 sm:p-8 text-gray-900 shadow max-h-[calc(100vh-16px)] sm:max-h-[calc(100vh-32px)] overflow-auto">
            <h2 className="text-xl">Eliminar projecto</h2>

            <Alert className="mt-5 w-full" variant="warning">
              <AlertTitle className="font-normal">
                Al eliminar el projecto todas sus tareas seran eliminadas
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
                Eliminar Projecto
              </Button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
