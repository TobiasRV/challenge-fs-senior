"use client";

import * as Dialog from "@radix-ui/react-dialog";
import { ReactNode } from "react";
import { Button } from "../../ui/button";
import { Badge } from "../../ui/badge";
import { ITask } from "@/src/stores/tasks/tasks.interface";
import { TaskStatusEnum } from "@/src/utils/enums";
import { TasksStatusTranslation } from "@/src/utils/translations";

type ShowTaskDetailsProps = {
  isOpen: boolean;
  handleClose: () => void;
  task: ITask;
};

export default function ShowTaskDetails({
  isOpen,
  handleClose,
  task,
}: ShowTaskDetailsProps): ReactNode {
  return (
    <Dialog.Root open={isOpen} modal>
      <Dialog.Overlay className="fixed inset-0 bg-black/50" />
      <Dialog.Title className="hidden">Detalles de Tarea</Dialog.Title>
      <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100%-2rem)] sm:w-full max-w-lg rounded-md bg-white p-4 sm:p-6 md:p-8 text-gray-900 shadow max-h-[90vh] overflow-y-auto">
        <h2 className="text-lg sm:text-xl font-semibold mb-4 sm:mb-6">Detalles de la Tarea</h2>

        <div className="space-y-3 sm:space-y-4">
          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-500">Título</p>
            <p className="text-sm sm:text-base mt-1 break-words">{task.title}</p>
          </div>

          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-500">Descripción</p>
            <p className="text-sm sm:text-base mt-1 break-words">
              {task.description || "Sin descripción"}
            </p>
          </div>

          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-500">Estado</p>
            <div className="mt-1">
                {TasksStatusTranslation[task.status]}
            </div>
          </div>

          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-500">Proyecto</p>
            <p className="text-sm sm:text-base mt-1 break-words">{task.projectName}</p>
          </div>

          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-500">Usuario Asignado</p>
            <p className="text-sm sm:text-base mt-1">{task.userName || "Sin asignar"}</p>
          </div>
        </div>

        <div className="mt-6 sm:mt-8 flex justify-end">
          <Button type="button" variant="outline" onClick={handleClose}>
            Cerrar
          </Button>
        </div>
      </Dialog.Content>
    </Dialog.Root>
  );
}
