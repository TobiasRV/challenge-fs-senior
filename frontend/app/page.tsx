"use client"
import { useAuthStore } from "@/src/stores/auth/auth";
import { UserRolesEnum } from "@/src/utils/enums";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useShallow } from "zustand/shallow";

export default function Home() {
  const { isLoggedIn, user } = useAuthStore(useShallow((state) => ({
    isLoggedIn: state.isLoggedIn,
    user: state.user
  })));

  const router = useRouter();

  useEffect(() => {
    if (isLoggedIn && user) {
      if ([UserRolesEnum.ADMIN, UserRolesEnum.MANAGER].includes(user.role) ) {
        router.push("/projects-dashboard")
      } else {
        router.push("/tasks-dashboard")
      }
    } else {
      router.push("/login")
    }
  }, [isLoggedIn, router])


  return (
      <div></div>
  );
}
