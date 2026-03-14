import { useState, useEffect, useCallback, createContext, useContext } from 'react';
import { AlertTriangle, CheckCircle2, Info, X, XCircle } from 'lucide-react';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface Toast {
  id: number;
  type: ToastType;
  title: string;
  message?: string;
  duration?: number;
}

interface ToastContextValue {
  toast: (type: ToastType, title: string, message?: string, duration?: number) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

let nextId = 0;

const icons: Record<ToastType, React.ReactNode> = {
  success: <CheckCircle2 className="w-5 h-5" />,
  error: <XCircle className="w-5 h-5" />,
  warning: <AlertTriangle className="w-5 h-5" />,
  info: <Info className="w-5 h-5" />,
};

const styles: Record<ToastType, { ring: string; icon: string; bar: string }> = {
  success: {
    ring: 'ring-emerald-500/30',
    icon: 'text-emerald-400 bg-emerald-500/10',
    bar: 'bg-emerald-500',
  },
  error: {
    ring: 'ring-red-500/30',
    icon: 'text-red-400 bg-red-500/10',
    bar: 'bg-red-500',
  },
  warning: {
    ring: 'ring-amber-500/30',
    icon: 'text-amber-400 bg-amber-500/10',
    bar: 'bg-amber-500',
  },
  info: {
    ring: 'ring-blue-500/30',
    icon: 'text-blue-400 bg-blue-500/10',
    bar: 'bg-blue-500',
  },
};

function ToastItem({ toast: t, onRemove }: { toast: Toast; onRemove: (id: number) => void }) {
  const [progress, setProgress] = useState(100);
  const [exiting, setExiting] = useState(false);
  const duration = t.duration ?? 4000;
  const s = styles[t.type];

  useEffect(() => {
    const start = Date.now();
    const tick = () => {
      const elapsed = Date.now() - start;
      const pct = Math.max(0, 100 - (elapsed / duration) * 100);
      setProgress(pct);
      if (pct > 0) requestAnimationFrame(tick);
    };
    const raf = requestAnimationFrame(tick);

    const timer = setTimeout(() => {
      setExiting(true);
      setTimeout(() => onRemove(t.id), 300);
    }, duration);

    return () => {
      cancelAnimationFrame(raf);
      clearTimeout(timer);
    };
  }, [duration, t.id, onRemove]);

  const handleClose = () => {
    setExiting(true);
    setTimeout(() => onRemove(t.id), 300);
  };

  return (
    <div
      className={`
        relative overflow-hidden rounded-xl
        bg-slate-900/95 backdrop-blur-xl border border-slate-700/60
        ring-1 ${s.ring}
        shadow-2xl shadow-black/40
        transition-all duration-300 ease-out
        ${exiting ? 'opacity-0 translate-x-8 scale-95' : 'opacity-100 translate-x-0 scale-100'}
      `}
      style={{ animation: exiting ? undefined : 'toast-in 0.4s cubic-bezier(0.16, 1, 0.3, 1)' }}
    >
      <div className="flex items-start gap-3 px-4 py-3.5">
        <div className={`shrink-0 p-1.5 rounded-lg ${s.icon}`}>
          {icons[t.type]}
        </div>
        <div className="flex-1 min-w-0 pt-0.5">
          <p className="text-sm font-semibold text-white leading-tight">{t.title}</p>
          {t.message && (
            <p className="text-xs text-slate-400 mt-1 leading-relaxed">{t.message}</p>
          )}
        </div>
        <button
          onClick={handleClose}
          className="shrink-0 p-1 rounded-md text-slate-500 hover:text-slate-300 hover:bg-slate-800 transition-colors"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
      {/* Progress bar */}
      <div className="h-[2px] w-full bg-slate-800">
        <div
          className={`h-full ${s.bar} transition-none`}
          style={{ width: `${progress}%` }}
        />
      </div>
    </div>
  );
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const remove = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback((type: ToastType, title: string, message?: string, duration?: number) => {
    const id = ++nextId;
    setToasts((prev) => [...prev, { id, type, title, message, duration }]);
  }, []);

  return (
    <ToastContext.Provider value={{ toast: addToast }}>
      {children}
      {/* Toast container */}
      <div className="fixed top-20 right-4 z-[100] flex flex-col gap-3 w-[380px] max-w-[calc(100vw-2rem)] pointer-events-none">
        <style>{`
          @keyframes toast-in {
            0% { opacity: 0; transform: translateX(100%) scale(0.9); }
            100% { opacity: 1; transform: translateX(0) scale(1); }
          }
        `}</style>
        {toasts.map((t) => (
          <div key={t.id} className="pointer-events-auto">
            <ToastItem toast={t} onRemove={remove} />
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
}
