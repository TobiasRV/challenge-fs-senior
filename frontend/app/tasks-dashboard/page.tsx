"use client";

import CreateTaskModal from "@/components/modals/tasks/createTaskModal";
import DeleteTaskModal from "@/components/modals/tasks/deleteTaskModal";
import EditTaskModal from "@/components/modals/tasks/editTaskModal";
import ShowTaskDetails from "@/components/modals/tasks/showTaskDetails";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useAuthStore } from "@/src/stores/auth/auth";
import { useTaskStore } from "@/src/stores/tasks/tasks";
import { IGetTasksParams, ITask } from "@/src/stores/tasks/tasks.interface";
import { UserRolesEnum } from "@/src/utils/enums";
import useDebounced from "@/src/utils/hooks/debounce";
import { TasksStatusTranslation } from "@/src/utils/translations";
import { ArrowLeft, Pencil, Trash } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useMemo, useState } from "react";
import { useShallow } from "zustand/shallow";

export default function TasksDashboard() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-gray-50 p-5">Cargando...</div>}>
      <TasksDashboardContent />
    </Suspense>
  );
}

function TasksDashboardContent() {
  const { user, isLoggedIn } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
      isLoggedIn: state.isLoggedIn,
    }))
  );

  const { tasks, getTasks, tasksLoading } = useTaskStore(
    useShallow((state) => ({
      tasks: state.tasks,
      getTasks: state.getTasks,
      tasksLoading: state.loading,
    }))
  );

  const searchParams = useSearchParams();
  const projectId = searchParams?.get("projectId");

  const router = useRouter();

  const [showCreateTaskModal, setShowCreateTaskModal] =
    useState<boolean>(false);

  const [updateTaskData, setUpdateTaskData] = useState<ITask | undefined>(
    undefined
  );

  const [deleteTaskData, setDeleteTaskData] = useState<ITask | undefined>(
    undefined
  );

  const [selectedTask, setSelectedTask] = useState<ITask | undefined>(
    undefined
  );

  const [tasksFilters, setTasksFilters] = useState<IGetTasksParams>({
    limit: 10,
    cursor: "",
  });

  useEffect(() => {
    if (!isLoggedIn || !user) {
      router.replace("/");
      return;
    }

    // If the user is an admin or manager it should have the project id as a query param
    if (user.role !== "Member" && !projectId) {
      router.replace("/");
    }
  }, [user, isLoggedIn]);

  const handleGetTasks = () => {
    const filters: IGetTasksParams = {
      ...tasksFilters,
    };

    if (user?.role != "Member") {
      filters.projectId = projectId!;
    }

    getTasks(filters);
  };

  useEffect(() => {
    handleGetTasks();
  }, [tasksFilters]);

  const getNext = async () => {
    setTasksFilters((prev) => ({
      ...prev,
      cursor: tasks.next,
    }));
  };

  const getPrev = async () => {
    setTasksFilters((prev) => ({
      ...prev,
      cursor: tasks.prev,
    }));
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setTasksFilters((prev) => ({
      ...prev,
      title: e.target.value,
    }));
  };

  const debounce = useDebounced(handleSearch, 1000);
  const onSearchChange = useMemo(() => debounce, []);

  const handleCreateTask = (success: boolean) => {
    if (success) {
      handleGetTasks();
    }

    setShowCreateTaskModal(false);
  };

  const handleEditTask = (success: boolean) => {
    if (success) {
      handleGetTasks();
    }

    setUpdateTaskData(undefined);
  };

  const handleDeleteTask = (success: boolean) => {
    if (success) {
      handleGetTasks();
    }

    setDeleteTaskData(undefined);
  };

  return (
    <main className="min-h-screen bg-gray-50">
      <div>
        {user?.role !== "Member" && (
          <Button variant="ghost" onClick={() => router.back()} className="m-5">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Volver
          </Button>
        )}
        <h1 className="p-5 text-4xl font-bold">Bienvenido {user?.username}!</h1>
        <h2 className="p-5 text-2xl font-bold">Estas son tus tareas:</h2>
      </div>
      {
        <div className="p-5 grid">
          <div className="flex justify-space-between gap-5">
            <Input onChange={onSearchChange} placeholder="Busqueda"></Input>
            {user?.role === "Manager" && (
              <Button onClick={() => setShowCreateTaskModal(true)}>
                Crear tarea
              </Button>
            )}
          </div>
          <div className="mt-10">
            {tasks.data.length ? (
              <div>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6 p-5">
                  {tasks.data?.map((task, key) => (
                    <Card
                      className="hover:shadow-lg transition-shadow cursor-pointer group"
                      key={key}
                      onClick={() => setSelectedTask(task)}
                    >
                      <CardContent>
                        <CardTitle className="flex justify-between items-center gap-2 min-w-0">
                          <span title={task.title}>{task.title}</span>
                          {user?.role === UserRolesEnum.MANAGER && (
                            <div className="flex-shrink-0 flex">
                              <Button
                                variant={"ghost"}
                                onClick={(
                                  e: React.MouseEvent<HTMLButtonElement>
                                ) => {
                                  e.stopPropagation();
                                  e.preventDefault();
                                  setUpdateTaskData(task);
                                }}
                              >
                                <Pencil />
                              </Button>
                              <Button
                                variant={"ghost"}
                                onClick={(
                                  e: React.MouseEvent<HTMLButtonElement>
                                ) => {
                                  e.stopPropagation();
                                  e.preventDefault();
                                  setDeleteTaskData(task);
                                }}
                              >
                                <Trash />
                              </Button>
                            </div>
                          )}
                        </CardTitle>
                        <p className="mt-5">
                          {`Estado: ${TasksStatusTranslation[task.status]}`}
                        </p>
                        {task.userName && (
                          <p className="mt-5">
                            {`Responsable: ${task.userName}`}
                          </p>
                        )}
                        <p className="mt-5">
                          {`Proyecto: ${task.projectName}`}
                        </p>
                      </CardContent>
                    </Card>
                  ))}
                </div>
                <div className="flex justify-end pr-10 gap-5">
                  <Button
                    disabled={!tasks.prev || tasksLoading}
                    onClick={getPrev}
                  >
                    Anterior
                  </Button>
                  <Button
                    disabled={!tasks.next || tasksLoading}
                    onClick={getNext}
                  >
                    Siguiente
                  </Button>
                </div>
              </div>
            ) : (
              <p className="text-center text-2xl w-full">
                No se encontraron datos
              </p>
            )}
          </div>
        </div>
      }

      <CreateTaskModal
        isOpen={showCreateTaskModal}
        handleClose={handleCreateTask}
        teamId={user?.teamId!}
        projectId={projectId!}
      />
      {updateTaskData && (
        <EditTaskModal
          isOpen={!!updateTaskData}
          teamId={user?.teamId!}
          task={updateTaskData!}
          handleClose={handleEditTask}
        />
      )}
      {deleteTaskData && (
        <DeleteTaskModal
          isOpen={!!deleteTaskData}
          task={deleteTaskData}
          handleClose={handleDeleteTask}
        />
      )}
      {selectedTask && (
        <ShowTaskDetails
          isOpen={!!selectedTask}
          task={selectedTask}
          handleClose={() => {
            setSelectedTask(undefined);
          }}
        />
      )}
    </main>
  );
}
