"use client";

import CreateProjectModal from "@/components/modals/project/createProjectModal";
import DeleteProjectModal from "@/components/modals/project/deleteProjectModal";
import EditProjectModal from "@/components/modals/project/editProjectModal";
import CreateTeamModal from "@/components/modals/team/createTeamModal";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useAuthStore } from "@/src/stores/auth/auth";
import { useProjectStore } from "@/src/stores/projects/projects";
import {
  IGetProjectParams,
  IProject,
} from "@/src/stores/projects/projects.interface";
import { useTeamStore } from "@/src/stores/teams/teams";
import { UserRolesEnum } from "@/src/utils/enums";
import useDebounced from "@/src/utils/hooks/debounce";
import { ProjectStatusTranslation } from "@/src/utils/translations";
import { Pencil, Trash } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { useShallow } from "zustand/shallow";

export default function ProjectsDashboard() {
  const { user, isLoggedIn } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
      isLoggedIn: state.isLoggedIn,
    }))
  );

  const { team, getTeamByOwner, teamLoading } = useTeamStore(
    useShallow((state) => ({
      team: state.team,
      getTeamByOwner: state.getTeamByOwner,
      teamLoading: state.loading,
    }))
  );

  const { projects, getProjects, projectLoading } = useProjectStore(
    useShallow((state) => ({
      projects: state.projects,
      getProjects: state.getProjects,
      projectLoading: state.loading,
    }))
  );

  const router = useRouter();

  const [showCreateTeamModal, setShowCreateTeamModal] =
    useState<boolean>(false);

  const [showCreateProjectModal, setShowCreateProjectModal] =
    useState<boolean>(false);

  // prettier-ignore
  const [updateProjectData, setUpdateProjectData] = useState<IProject | undefined>(undefined);
  // prettier-ignore
  const [deleteProjectData, setDeleteProjectData] = useState<IProject | undefined>(undefined);

  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const [projectFilters, setProjectFilter] = useState<IGetProjectParams>({
    cursor: "",
    limit: 20,
    withStats: true
  });

  useEffect(() => {
    if (!isLoggedIn || !user || user.role === UserRolesEnum.MEMBER) {
      router.replace("/");
    }
  }, [user, isLoggedIn]);

  useEffect(() => {
    if (user?.role === UserRolesEnum.ADMIN) {
      getTeamByOwner();
    }
  }, []);

  useEffect(() => {
    if (
      user?.role === UserRolesEnum.ADMIN &&
      !teamLoading &&
      team &&
      !team.exists
    ) {
      setShowCreateTeamModal(true);
    } else {
      setShowCreateTeamModal(false);
    }
  }, [team, teamLoading, user]);

  useEffect(() => {
    handleGetProjects();
  }, [projectFilters, team]);

  const handleGetProjects = () => {
    if (!user) {
      router.replace("/");
    }

    let extraFilters = {};

    if (user?.role === UserRolesEnum.ADMIN) {
      if (!team.team?.id) {
        return;
      }

      extraFilters = {
        teamId: team.team?.id,
      };
    }

    if (user?.role === UserRolesEnum.MANAGER) {
      extraFilters = {
        managerId: user.id,
      };
    }

    getProjects({
      ...projectFilters,
      ...extraFilters,
    });
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setProjectFilter((prev) => ({
      ...prev,
      name: e.target.value,
    }));
  };

  const debounce = useDebounced(handleSearch, 1000);
  const onSearchChange = useMemo(() => debounce, []);

  const getNext = async () => {
    setProjectFilter((prev) => ({
      ...prev,
      cursor: projects.next,
    }));
  };

  const getPrev = async () => {
    setProjectFilter((prev) => ({
      ...prev,
      cursor: projects.prev,
    }));
  };

  const handleCreateProject = (success: boolean) => {
    if (success) {
      handleGetProjects();
    }
    setShowCreateProjectModal(false);
  };

  const handleEditProject = (success: boolean) => {
    if (success) {
      handleGetProjects();
    }

    setUpdateProjectData(undefined);
  };

  const handleDeleteProject = (success: boolean) => {
    if (success) {
      handleGetProjects();
    }

    setDeleteProjectData(undefined);
  };

  if (!mounted) {
    return null;
  }

  return (
    <main className="min-h-screen bg-gray-50">
      <div>
        <h1 className="p-5 text-4xl font-bold">Bienvenido {user?.username}!</h1>
        <h2 className="p-5 text-2xl font-bold">Estos son tus proyectos:</h2>
      </div>
      {
        <div className="p-5 grid">
          <div className="flex justify-space-between gap-5">
            <Input onChange={onSearchChange} placeholder="Busqueda"></Input>
            {user?.role === "Manager" && (
              <Button onClick={() => setShowCreateProjectModal(true)}>
                Crear proyecto
              </Button>
            )}
          </div>
          <div className="mt-10">
            {projects.data.length ? (
              <div>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6 p-5">
                  {projects.data?.map((project, key) => (
                    <Link href={`/tasks-dashboard?projectId=${project.id}`} key={key}>
                      <Card className="hover:shadow-lg transition-shadow cursor-pointer group">
                        <CardContent>
                          <CardTitle className="flex justify-between items-center gap-2">
                            <span title={project.name}>
                              {project.name}
                            </span>
                            {user?.role === UserRolesEnum.MANAGER && (
                              <div className="flex-shrink-0 flex">
                                <Button
                                  variant={"ghost"}
                                  size="icon"
                                  onClick={(
                                    e: React.MouseEvent<HTMLButtonElement>
                                  ) => {
                                    e.stopPropagation();
                                    e.preventDefault();
                                    setUpdateProjectData(project);
                                  }}
                                >
                                  <Pencil />
                                </Button>
                                <Button
                                  variant={"ghost"}
                                  size="icon"
                                  onClick={(
                                    e: React.MouseEvent<HTMLButtonElement>
                                  ) => {
                                    e.stopPropagation();
                                    e.preventDefault();
                                    setDeleteProjectData(project);
                                  }}
                                >
                                  <Trash />
                                </Button>
                              </div>
                            )}
                          </CardTitle>
                          <p className="mt-5">
                            {ProjectStatusTranslation[project.status]}
                          </p>
                          <ul className="list-disc ml-5 mt-5">
                            <li>Tareas por hacer: {project.toDoTasks}</li>
                            <li>Tareas en curso: {project.inProgressTasks}</li>
                            <li>Tareas completas: {project.doneTasks}</li>
                          </ul>
                        </CardContent>
                      </Card>
                    </Link>
                  ))}
                </div>
                <div className="flex justify-end pr-10 gap-5">
                  <Button
                    disabled={!projects.prev || projectLoading}
                    onClick={getPrev}
                  >
                    Anterior
                  </Button>
                  <Button
                    disabled={!projects.next || projectLoading}
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

      <CreateTeamModal
        isOpen={showCreateTeamModal}
        handleClose={() => setShowCreateTeamModal(false)}
      />
      <CreateProjectModal
        isOpen={showCreateProjectModal}
        handleClose={handleCreateProject}
      />
      {updateProjectData && (
        <EditProjectModal
          isOpen={!!updateProjectData}
          project={updateProjectData}
          handleClose={handleEditProject}
        />
      )}
      {deleteProjectData && (
        <DeleteProjectModal
          isOpen={!!deleteProjectData}
          project={deleteProjectData}
          handleClose={handleDeleteProject}
        />
      )}
    </main>
  );
}
