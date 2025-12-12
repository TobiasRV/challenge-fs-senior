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

  const { registerAdmin, error, authStatusCode, clearRequestState } = useAuthStore(
    useShallow((state) => ({
      registerAdmin: state.registerAdmin,
      error: state.error,
      authStatusCode: state.statusCode,
      clearRequestState: state.clearRequestState
    })),
  );

    const { emailAlreadyExists, userLoading } = useUserStore(
    useShallow((state) => ({
      emailAlreadyExists: state.emailAlreadyExists,
      userLoading: state.loading,
      userError: state.error
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
      const timeoutAux = setTimeout(
        () => clearRequestState(),
        7000,
      );
      setAlertTimeout(timeoutAux);
    }
  }, [error, authStatusCode]);

  return (
    <div className="min-h-screen bg-gray-50">
      <main className="container mx-auto mt-10 px-4 py-6 sm:py-8 flex justify-center border size-100 md:size-1/2">
        <div className="mb-6 sm:mb-8 w-100">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
            Registro
          </h1>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-6">
              <Label>Usuario</Label>
              <Input
                className="mt-2"
                {...register("username", {
                  required: { value: true, message: "El usuario es requerido"},
                })}
              />
              {errors.username ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.username.message}
                </p>
              ) : null}
            </div>
            <div className="mt-6">
              <Label>Contraseña</Label>
              <Input
                type="password"
                className="mt-2"
                {...register("password", {
                  required: { value: true, message: "La contraseña es requerida"},
                })}
              />
              {errors.password ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.password.message}
                </p>
              ) : null}
            </div>
            <div className="mt-6">
              <Label>Email</Label>
              <Input
                className="mt-2"
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
                    }
                  },
                })}
              />
              {errors.email ? (
                <p className="text-red-500 text-xs mt-2">
                  {errors.email.message}
                </p>
              ) : null}
            </div>

            {(error && authStatusCode) && (
              <Alert className="mt-5 w-100" variant="error">
                <AlertTitle className="font-normal">
                  {mapErrorMessage[authStatusCode]}
                </AlertTitle>
              </Alert>
            )}

            <div className="mt-6">
              <Button disabled={userLoading || Object.keys(errors).length !== 0}>Registrarse</Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
