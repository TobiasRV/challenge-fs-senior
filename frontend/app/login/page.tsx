"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { SubmitHandler, useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { validateEmail } from "@/src/utils/validations";

import { useShallow } from "zustand/react/shallow";
import { HttpStatusCode } from "axios";
import { useEffect, useState } from "react";
import { Alert, AlertTitle } from "@/components/ui/alert";
import { clearLs } from "@/src/utils/localStorage";
import { useAuthStore } from "@/src/stores/auth/auth";

type Inputs = {
  email: string;
  password: string;
};

export default function LoginPage() {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<Inputs>();
  const router = useRouter();

  const [alert, setAlert] = useState<{
    type: "error"
    message: string
  } | undefined>(undefined)

  const [alertTimeout, setAlertTimeout] = useState<NodeJS.Timeout>();

  const {logIn, clearState} = useAuthStore(useShallow((state) => ({
    logIn: state.logIn,
    clearState: state.clearState
  })))

  useEffect(() => {
    clearLs()
    clearState()
  }, [])

  useEffect(() => {
      if (alert) {
        clearTimeout(alertTimeout);
        const timeoutAux = setTimeout(
          () => setAlert(undefined),
          7000,
        );
        setAlertTimeout(timeoutAux);
      }
    }, [alert]);

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    const response = await logIn({
      email: data.email,
      password: data.password
    });

    switch(response) {
      case HttpStatusCode.Ok:
        router.push("/")
        break;
      case HttpStatusCode.Conflict:
        setAlert({
          type: "error",
          message: "Usuario o contraseña invalidos"
        });
        break;
      default:
        setAlert({
          type: "error",
          message: "Error iniciando sesion"
        })
    }
  };

  const redirectToSignIn = () => {
    router.push("/signin");
  };

  return (
    <div className="min-h-screen p-4 sm:p-6 md:p-10 bg-gray-100 flex items-center justify-center">
      <main className="w-full max-w-md sm:max-w-lg md:max-w-xl lg:max-w-2xl mx-auto px-4 sm:px-6 md:px-8 py-6 sm:py-8 md:py-10 bg-white rounded-lg shadow-md">
        <div className="w-full">
          <h1 className="text-xl sm:text-2xl md:text-3xl font-bold text-gray-900 mb-4 sm:mb-6 text-center">
            Log In
          </h1>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 sm:space-y-6">
            <div>
              <Label className="text-sm sm:text-base">Email</Label>
              <Input
                className="mt-1 sm:mt-2 w-full"
                {...register("email", {
                  required: { value: true, message: "El email es requerido"},
                  validate: {
                    isValidEmail: (email) => {
                      const isValidEmail = validateEmail(email);
                      if (!isValidEmail) {
                        return "Email invalido";
                      }
                      return isValidEmail;
                    },
                  },
                })}
              />
              {errors.email ? (
                <p className="text-red-500 text-xs sm:text-sm mt-1 sm:mt-2">
                  {errors.email.message}
                </p>
              ) : null}
            </div>
            <div>
              <Label className="text-sm sm:text-base">Contraseña</Label>
              <Input
                type="password"
                className="mt-1 sm:mt-2 w-full"
                {...register("password", {
                  required: { value: true, message: "La contraseña es requerida"},
                })}
              />
              {errors.password ? (
                <p className="text-red-500 text-xs sm:text-sm mt-1 sm:mt-2">
                  {errors.password.message}
                </p>
              ) : null}
            </div>
            <div className="pt-1">
              <Button className="p-0" variant="link" onClick={redirectToSignIn} type="button">
                <p className="underline text-sm sm:text-base">Registrarse</p>
              </Button>
            </div>
            {alert && (
              <Alert className="mt-3 sm:mt-5" variant={alert.type}>
                <AlertTitle className="font-normal text-sm sm:text-base">
                  {alert.message}
                </AlertTitle>
              </Alert>
            )}
            <div className="pt-2 sm:pt-4">
              <Button 
                className="w-full sm:w-auto" 
                disabled={Object.keys(errors).length !== 0}
              >
                Iniciar Sesion
              </Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
