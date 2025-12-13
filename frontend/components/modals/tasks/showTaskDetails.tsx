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
      <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 max-w-lg min-w-md rounded-md bg-white p-8 text-gray-900 shadow">
        <h2 className="text-xl font-semibold mb-6">Detalles de la Tarea</h2>

        <div className="space-y-4">
          <div>
            <p className="text-sm font-medium text-gray-500">Título</p>
            <p className="text-base mt-1">{task.title}</p>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">Descripción</p>
            <p className="text-base mt-1">
              {task.description || "Sin descripción"}
            </p>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">Estado</p>
            <div className="mt-1">
                {TasksStatusTranslation[task.status]}
            </div>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">Proyecto</p>
            <p className="text-base mt-1">{task.projectName}</p>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">Usuario Asignado</p>
            <p className="text-base mt-1">{task.userName || "Sin asignar"}</p>
          </div>
        </div>

        <div className="mt-8 flex justify-end">
          <Button type="button" variant="outline" onClick={handleClose}>
            Cerrar
          </Button>
        </div>
      </Dialog.Content>
    </Dialog.Root>
  );
}
