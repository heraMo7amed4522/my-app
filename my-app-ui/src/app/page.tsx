'use client';


export default function SplashPage() {

  return (
    <main
      className="relative h-screen bg-cover bg-center flex items-center justify-center"
      style={{
        backgroundImage: "url('/images/splash-background.jpg')",
      }}
    >
      <div className="absolute inset-0 bg-black opacity-60"></div>

      <div className="relative z-10 text-center text-white px-4">
        <h1 className="text-5xl font-bold">Welcome</h1>
        <p className="text-xl mt-2">Loading your journey...</p>
      </div>
    </main>
  );
}