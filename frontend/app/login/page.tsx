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
    <div className="min-h-screen p-10 bg-gray-100">
      <main className="container mx-auto mt-10 px-4 py-4 md:py-6 flex justify-center border size-100 md:size-1/2 bg-[#FFFFFF]">
        <div className="mb-6 sm:mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
            Log In
          </h1>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mt-6">
              <Label>Email</Label>
              <Input
                className="mt-2"
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
                <p className="text-red-500 text-xs mt-2">
                  {errors.email.message}
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
            <div className="mt-3">
              <Button className="p-0" variant="link" onClick={redirectToSignIn}>
                <p className="underline">Registrarse</p>
              </Button>
            </div>
            {alert && (
              <Alert className="mt-5" variant={alert.type}>
                <AlertTitle className="font-normal">
                  {alert.message}
                </AlertTitle>
              </Alert>
            )}
            <div className="mt-6">
              <Button disabled={Object.keys(errors).length !== 0}>Iniciar Sesion</Button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
