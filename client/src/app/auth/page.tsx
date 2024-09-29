import React from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import LoginForm from '@/components/auth/login-form'
import RegisterForm from '@/components/auth/register-form'

const AuthForm = () => {
  return (
    <div className="flex items-center justify-center min-h-screen bg-green-50">
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle className="text-2xl font-bold text-center text-green-700">Welcome to Dotwork</CardTitle>
        <CardDescription className="text-center text-green-600">Log in or create an account</CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="login">Login</TabsTrigger>
            <TabsTrigger value="register">Register</TabsTrigger>
          </TabsList>
          <TabsContent value="login">
            <LoginForm />

          </TabsContent>
          <TabsContent value="register">
            <RegisterForm />
          </TabsContent>
        </Tabs>
      </CardContent>
      <CardFooter className="flex justify-center">
        <Button variant="link" className="text-green-600 hover:text-green-700">
          Forgot your password?
        </Button>
      </CardFooter>
    </Card>
  </div>
  )
}

export default AuthForm