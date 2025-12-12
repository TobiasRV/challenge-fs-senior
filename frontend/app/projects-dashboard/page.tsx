import { useAuthStore } from "@/src/stores/auth/auth";
import { useShallow } from "zustand/shallow";

export default function ProjectsDashboard() {

    const { user } = useAuthStore(
        useShallow((state) => ({
        user: state.user,
        })),
    );    

    return (
        <main className="min-h-screen bg-gray-50">
            <div>
                <h1 className="p-5 text-xl font-bold">Bienvenido {user?.username}</h1>
            </div>
        </main>
    )
}