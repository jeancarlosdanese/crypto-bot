// File: components/Spinner.tsx

interface SpinnerProps {
  className?: string;
}

export default function Spinner({ className }: SpinnerProps) {
  return (
    <div className={`flex justify-center items-center min-h-screen ${className}`}>
      <div className="animate-spin rounded-full h-12 w-12 border-t-4 border-stone-800 border-solid"></div>
    </div>
  );
}
