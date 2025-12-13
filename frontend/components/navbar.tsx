"use client";

import { useAuthStore } from "@/src/stores/auth/auth";
import { UserRolesEnum } from "@/src/utils/enums";
import { ArrowBigLeft, ArrowRight, Trello } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";
import { useShallow } from "zustand/shallow";
import { Button } from "./ui/button";

export default function Navbar() {
  const { isLoggedIn, user } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
      isLoggedIn: state.isLoggedIn
    }))
  );

  const router = useRouter();
  const pathname = usePathname();

  return (
    <header className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
      <div className="conatiner mx-auto px-4 py-3 sm:py-4 flex items-center justify-between">
        <div className="flex items-center">
          <button
            type="button"
            onClick={() => router.push("/")}
            className="flex items-center space-x-2 hover:cursor-pointer"
          >
            <Trello className="h-6 w-6 sm:h-8 sm:w-8 text-blue-600" />
            <span className="text-xl sm:text-2xl font-bold text-gray-900">
              Task Manager
            </span>
          </button>
        </div>

        {isLoggedIn && user?.role === UserRolesEnum.ADMIN && pathname !== "/users-dashboard" && (
          <Button
            type="button"
            onClick={() => router.push("/users-dashboard")}
            className="ml-4 px-3 py-1 rounded text-white"
          >
            <div className="flex items-center space-x-2">
                <p className="m-0">Users dashboard</p>
                <ArrowRight />
            </div>
          </Button>
        )}
      </div>
    </header>
  );
}
