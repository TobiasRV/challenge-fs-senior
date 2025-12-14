"use client";

import { useAuthStore } from "@/src/stores/auth/auth";
import { UserRolesEnum } from "@/src/utils/enums";
import { ArrowRight, LogOut, Menu, Trello, X } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";
import { useState } from "react";
import { useShallow } from "zustand/shallow";
import { Button } from "./ui/button";

export default function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const { isLoggedIn, user, logOut } = useAuthStore(
    useShallow((state) => ({
      user: state.user,
      isLoggedIn: state.isLoggedIn,
      logOut: state.logOut
    }))
  );

  const router = useRouter();
  const pathname = usePathname();

  const handleLogOut = () => {
    logOut();
    setIsMenuOpen(false);
    router.push("/login");
  };

  const handleNavigateToUsers = () => {
    setIsMenuOpen(false);
    router.push("/users-dashboard");
  };

  const showUserDashboardButton = isLoggedIn && user?.role === UserRolesEnum.ADMIN && pathname !== "/users-dashboard";

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

        <div className="hidden sm:flex items-center">
          {showUserDashboardButton && (
            <Button
              type="button"
              onClick={() => router.push("/users-dashboard")}
              className="ml-4 px-3 py-1 rounded text-white"
            >
              <div className="flex items-center space-x-2 gap-4">
                <p className="m-0">Users dashboard</p>
                <ArrowRight />
              </div>
            </Button>
          )}
          {isLoggedIn && (
            <Button
              type="button"
              onClick={handleLogOut}
              className="ml-4 px-3 py-1 rounded text-white"
            >
              <div className="flex items-center space-x-2 gap-4">
                <p className="m-0">Log Out</p>
                <LogOut />
              </div>
            </Button>
          )}
        </div>

        {isLoggedIn && (
          <button
            type="button"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="sm:hidden p-2 rounded-md hover:bg-gray-100 transition-colors"
            aria-label="Toggle menu"
          >
            {isMenuOpen ? (
              <X className="h-6 w-6 text-gray-600" />
            ) : (
              <Menu className="h-6 w-6 text-gray-600" />
            )}
          </button>
        )}
      </div>

      {isMenuOpen && isLoggedIn && (
        <div className="sm:hidden border-t bg-white px-4 py-2 space-y-2">
          {showUserDashboardButton && (
            <Button
              type="button"
              onClick={handleNavigateToUsers}
              className="w-full px-3 py-2 rounded text-white"
            >
              <div className="flex items-center justify-center space-x-2 gap-4">
                <p className="m-0">Users dashboard</p>
                <ArrowRight />
              </div>
            </Button>
          )}
          <Button
            type="button"
            onClick={handleLogOut}
            className="w-full px-3 py-2 rounded text-white"
          >
            <div className="flex items-center justify-center space-x-2 gap-4">
              <p className="m-0">Log Out</p>
              <LogOut />
            </div>
          </Button>
        </div>
      )}
    </header>
  );
}
