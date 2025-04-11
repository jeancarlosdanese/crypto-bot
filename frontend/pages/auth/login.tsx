// File: pages/auth/login.tsx

import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import { useUser } from "@/context/UserContext";
import { useRouter } from "next/router";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp";
import { toast } from "sonner";
import { Mail, Loader2 } from "lucide-react";
import { useRecaptcha } from "@/hooks/useRecaptcha"; // Hook de reCAPTCHA

const loginSchema = z.object({
  identifier: z.string().min(5, "Informe seu e-mail ou WhatsApp"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const [otpSent, setOtpSent] = useState(false);
  const [identifier, setIdentifier] = useState("");
  const [otp, setOtp] = useState("");
  const { refreshUser } = useUser();
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const recaptchaToken = useRecaptcha("login"); // Captura o token do reCAPTCHA

  const otpSchema = z
    .string()
    .length(8, "O OTP deve ter exatamente 8 dígitos.")
    .regex(/^\d+$/, "O OTP deve conter apenas números.");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) });

  const sendOTP = async (data: LoginForm) => {
    if (!recaptchaToken) {
      toast.error("❌ Erro ao validar reCAPTCHA. O token está indefinido.");
      return;
    }

    setIsLoading(true);
    try {
      await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/request-otp`, {
        ...data,
        recaptchaToken, // Enviando o token do reCAPTCHA para o backend
      });

      setIdentifier(data.identifier);
      setOtpSent(true);
      toast.success("Código OTP enviado!");
    } catch (error: any) {
      const message =
        error.response?.data?.message || "Erro ao solicitar OTP. Verifique seus dados.";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const verifyOTP = async (otpValue?: string) => {
    const otpToValidate = otpValue || otp;

    const validation = otpSchema.safeParse(otpToValidate);

    if (!validation.success) {
      toast.error(validation.error.errors[0].message);
      return;
    }

    setIsLoading(true);
    try {
      const response = await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/auth/verify-otp`, {
        identifier,
        otp: otpToValidate,
      });

      localStorage.setItem("token", response.data.token);
      await refreshUser();
      toast.success("Login realizado com sucesso!");
      router.push("/dashboard");
    } catch (error: any) {
      const message = error.response?.data?.message || "Código inválido. Tente novamente.";
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <Card className="w-full max-w-md shadow-lg">
        <CardHeader>
          <CardTitle className="text-center text-2xl">Login</CardTitle>
        </CardHeader>
        <CardContent>
          {!otpSent ? (
            <form onSubmit={handleSubmit(sendOTP)} className="space-y-4">
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground" />
                <Input
                  {...register("identifier")}
                  placeholder="E-mail ou WhatsApp"
                  className="pl-10"
                  disabled={isLoading}
                />
              </div>
              {errors.identifier && (
                <p className="text-destructive text-sm">{errors.identifier.message}</p>
              )}
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? <Loader2 className="animate-spin mr-2" /> : "Enviar OTP"}
              </Button>
            </form>
          ) : (
            <div className="space-y-4">
              <p className="text-center text-sm text-muted-foreground">
                Digite o código enviado para <strong>{identifier}</strong>
              </p>

              <InputOTP
                maxLength={8}
                value={otp}
                onChange={(value) => {
                  if (/^\d*$/.test(value)) {
                    setOtp(value);
                    if (value.length === 8) {
                      verifyOTP(value);
                    }
                  }
                }}
              >
                <InputOTPGroup className="flex gap-1 sm:gap-2 justify-center w-full max-w-md">
                  {[...Array(8)].map((_, index) => (
                    <InputOTPSlot
                      key={index}
                      index={index}
                      className="w-10 h-12 sm:w-12 sm:h-14 text-xl text-center border border-gray-300 rounded-md focus:ring-2 focus:ring-primary"
                    />
                  ))}
                </InputOTPGroup>
              </InputOTP>

              <Button
                onClick={() => verifyOTP()}
                className="w-full"
                disabled={isLoading || otp.length < 8}
              >
                {isLoading ? <Loader2 className="animate-spin mr-2" /> : "Confirmar Código"}
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
