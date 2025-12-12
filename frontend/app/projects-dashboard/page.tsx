"use client"

import CreateTeamModal from "@/components/modals/createTeamModal";
import { useAuthStore } from "@/src/stores/auth/auth";
import { useTeamStore } from "@/src/stores/teams/teams";
import { UserRolesEnum } from "@/src/utils/enums";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useShallow } from "zustand/shallow";

export default function ProjectsDashboard() {
  const { user, isLoggedIn } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
      isLoggedIn: state.isLoggedIn,
    })),
  );

  const { team, getTeamByOwner, teamLoading } = useTeamStore(
    useShallow((state) => ({
        team: state.team,
        getTeamByOwner: state.getTeamByOwner,
        teamLoading: state.loading
    }))
  )

  const router = useRouter();

  const [showCreateTeamModal, setShowCreateTeamModal] = useState<boolean>(false)

  useEffect(() => {
    if (!isLoggedIn || !user || user.role === UserRolesEnum.MEMBER) {
        router.replace("/")
    }
  }, [user, isLoggedIn]);


  useEffect(() => {
    getTeamByOwner();
  }, [])

  useEffect(() => {
    if (!team.exists) {
        setShowCreateTeamModal(true)
    } else {
        setShowCreateTeamModal(false)
    }
  }, [team])


  return (
    <main className="min-h-screen bg-gray-50">
      <div>
        <h1 className="p-5 text-xl font-bold">Bienvenido {user?.username}</h1>
      </div>
        

      <CreateTeamModal isOpen={showCreateTeamModal} handleClose={() => setShowCreateTeamModal(false)}/>
    </main>
  );
}
