"use client";

import * as Dialog from "@radix-ui/react-dialog";
import * as Select from "@radix-ui/react-select";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { ReactNode, useState, useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
import { useShallow } from "zustand/shallow";
import { Alert, AlertTitle } from "../../ui/alert";
import { Button } from "../../ui/button";
import { Input } from "../../ui/input";
import { useProjectStore } from "@/src/stores/projects/projects";
import { ChevronDown } from "lucide-react";
import { ProjectStatusEnum } from "@/src/utils/enums";
import { ProjectStatusTranslation } from "@/src/utils/translations";
import { IProject } from "@/src/stores/projects/projects.interface";


type Inputs = {
  name: string;
  status: ProjectStatusEnum;
};

type EditProjectModalProps = {
    isOpen: boolean;
    project: IProject;
    handleClose: (success: boolean) => void;
}


type MapErrorMessage = {
  [key: number]: string;
};

export default function EditProjectModal({isOpen, project, handleClose}: EditProjectModalProps): ReactNode {
  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
    control
  } = useForm<Inputs>({
    reValidateMode: "onChange",
    defaultValues: {
      name: project.name,
      status: project.status
    }
  });

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { updateProject, error, statusCode, loading, clearRequestState } = useProjectStore(
    useShallow((state) => ({
        updateProject: state.updateProject,
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
    const response = await updateProject({
      name: data.name,
      status: data.status,
      id: project.id
    });
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
        "Error al actualizando projecto. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
    // Generic responses for edge case errors that should not happend but could happend
    [HttpStatusCode.NotFound]: "Error al actualizar projecto.",
    [HttpStatusCode.Forbidden]: "Error al actualizar projecto.",
    [HttpStatusCode.Unauthorized]: "Error al actualizar projecto.",
};

  return (
    <div>
      <Dialog.Root open={isOpen} modal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Title className="hidden">Editar projecto</Dialog.Title>
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 max-w-md min-w-md rounded-md bg-white p-8 text-gray-900 shadow">
          <h2 className="text-xl">Editar projecto</h2>

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
            </div>
            <div className="mt-6">
              <div>
                <Label>Estado</Label>
              </div>

              <Controller
                name="status"
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
                            {[ProjectStatusEnum.ON_HOLD, ProjectStatusEnum.IN_PROGRESS, ProjectStatusEnum.COMPLETED].map((status, idx) => (
                              <Select.Item
                                value={status}
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
                                  {ProjectStatusTranslation[status as ProjectStatusEnum]}
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
            <div className="flex end justify-end gap-3 mt-10">
                <Button type="button" variant="ghost" onClick={() => close(false)}>
                  Cancelar
                </Button>
                <Button
                  disabled={loading || Object.keys(errors).length !== 0 || !isDirty}
                  type="submit"
                >
                  Editar projecto
                </Button>
              </div>
          </form>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
}
