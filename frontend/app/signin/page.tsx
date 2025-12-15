"use client";

import { useAuthStore } from "@/src/stores/auth/auth";
import { Label } from "@radix-ui/react-label";
import { HttpStatusCode } from "axios";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { useShallow } from "zustand/react/shallow";
import { Input } from "@/components/ui/input";
import { Alert, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { useUserStore } from "@/src/stores/users/users";
import { validateEmail } from "@/src/utils/validations";
import { ArrowLeftIcon } from "lucide-react";

type Inputs = {
  username: string;
  password: string;
  email: string;
};

type MapErrorMessage = {
  [key: number]: string;
};

export default function SignInPage() {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    setError,
  } = useForm<Inputs>();
  const router = useRouter();

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const { registerAdmin, error, authStatusCode, clearRequestState } =
    useAuthStore(
      useShallow((state) => ({
        registerAdmin: state.registerAdmin,
        error: state.error,
        authStatusCode: state.statusCode,
        clearRequestState: state.clearRequestState,
      })),
    );

  const { emailAlreadyExists, userLoading } = useUserStore(
    useShallow((state) => ({
      emailAlreadyExists: state.emailAlreadyExists,
      userLoading: state.loading,
      userError: state.error,
    })),
  );

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const statusCode = await registerAdmin(data);
    if (statusCode === HttpStatusCode.Created) {
      router.push("/");
    }
  };

  const mapErrorMessage: MapErrorMessage = {
    [HttpStatusCode.InternalServerError]:
      "Error al registrar usuario. Por favor intente nuevamente.",
    [HttpStatusCode.BadRequest]: "Datos incorrectos.",
  };

  useEffect(() => {
    if (error) {
      clearTimeout(alertTimeout);
      const timeoutAux = setTimeout(() => clearRequestState(), 7000);
      setAlertTimeout(timeoutAux);
    }
  }, [error, authStatusCode]);

  const goBack = () => {
    router.push("/login")
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4 py-8 sm:py-12">
      <main className="w-full max-w-sm sm:max-w-md lg:max-w-lg bg-white rounded-lg shadow-md p-6 sm:p-8 lg:p-10">
        <div className="w-full">
          <div className="flex align-center gap-25">
            <Button onClick={goBack}>
            <ArrowLeftIcon />
          </Button>
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 mb-4 sm:mb-6 text-center">
            Registro
          </h1>
          </div>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 sm:space-y-6">
            <div>
              <Label className="text-sm sm:text-base font-medium">Usuario</Label>
              <Input
                className="mt-1.5 sm:mt-2 w-full"
                {...register("username", {
                  required: { value: true, message: "El usuario es requerido" },
                })}
              />
              {errors.username ? (
                <p className="text-red-500 text-xs sm:text-sm mt-1.5">
                  {errors.username.message}
                </p>
              ) : null}
            </div>
            <div>
              <Label className="text-sm sm:text-base font-medium">Contraseña</Label>
              <Input
                type="password"
                className="mt-1.5 sm:mt-2 w-full"
                {...register("password", {
                  required: {
                    value: true,
                    message: "La contraseña es requerida",
                  },
                })}
              />
              {errors.password ? (
                <p className="text-red-500 text-xs sm:text-sm mt-1.5">
                  {errors.password.message}
                </p>
              ) : null}
            </div>
            <div>
              <Label className="text-sm sm:text-base font-medium">Email</Label>
              <Input
                className="mt-1.5 sm:mt-2 w-full"
                {...register("email", {
                  required: true,
                  validate: {
                    isValidEmail: (email) => {
                      const isValidEmail = validateEmail(email);
                      if (!isValidEmail) {
                        setError("email", {
                          message: "Email invalido",
                          type: "invalid-email",
                        });
                        return "Email invalido";
                      }
                      return isValidEmail;
                    },
                    emailAlreadyExists: async (email) => {
                      const emailExists = await emailAlreadyExists(email);
                      if (emailExists) {
                        setError("email", {
                          message: "Ya existe una cuenta con ese email",
                          type: "email-already-exists",
                        });

                        return "Ya existe una cuenta con ese email";
                      }

                      return true;
                    },
                  },
                })}
              />
              {errors.email ? (
                <p className="text-red-500 text-xs sm:text-sm mt-1.5">
                  {errors.email.message}
                </p>
              ) : null}
            </div>

            {error && authStatusCode && (
              <Alert className="w-full" variant="error">
                <AlertTitle className="font-normal text-sm sm:text-base">
                  {mapErrorMessage[authStatusCode]}
                </AlertTitle>
              </Alert>
            )}

            <div className="pt-2">
              <Button
                className="w-full sm:w-auto"
                disabled={userLoading || Object.keys(errors).length !== 0}
              >
                Registrarse
              </Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
