import { useEffect, useState } from 'react';
import { ShieldAlert, ArrowLeft, Sparkles } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface ErrorModalProps {
  title?: string;
  message: string;
  onClose?: () => void;
}

export default function ErrorModal({ title = 'Transaction Blocked', message, onClose }: ErrorModalProps) {
  const navigate = useNavigate();
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    requestAnimationFrame(() => setVisible(true));
  }, []);

  const handleBack = () => {
    setVisible(false);
    setTimeout(() => {
      onClose?.();
      navigate('/');
    }, 300);
  };

  return (
    <div
      className={`fixed inset-0 z-[90] flex items-center justify-center p-4 transition-all duration-300 ${
        visible ? 'bg-slate-950/80 backdrop-blur-md' : 'bg-transparent backdrop-blur-none'
      }`}
    >
      <div
        className={`
          relative w-full max-w-sm overflow-hidden rounded-2xl
          bg-slate-900 border border-slate-700/60
          shadow-2xl shadow-red-900/10
          transition-all duration-500 ease-out
          ${visible ? 'opacity-100 scale-100 translate-y-0' : 'opacity-0 scale-90 translate-y-8'}
        `}
      >
        {/* Glow effect top */}
        <div className="absolute -top-24 left-1/2 -translate-x-1/2 w-64 h-40 bg-red-500/20 rounded-full blur-3xl pointer-events-none" />

        <div className="relative px-6 pt-8 pb-6 flex flex-col items-center text-center">
          {/* Animated icon */}
          <div className="relative mb-5">
            <div className="absolute inset-0 animate-ping rounded-full bg-red-500/20" style={{ animationDuration: '2s' }} />
            <div className="relative p-4 rounded-full bg-gradient-to-br from-red-500/20 to-red-600/10 border border-red-500/30 shadow-lg shadow-red-500/10">
              <ShieldAlert className="w-8 h-8 text-red-400" />
            </div>
          </div>

          {/* Title */}
          <h2 className="text-xl font-bold text-white mb-2">{title}</h2>

          {/* Message */}
          <p className="text-sm text-slate-400 leading-relaxed max-w-[280px]">{message}</p>

          {/* Divider with sparkle */}
          <div className="flex items-center gap-3 my-5 w-full">
            <div className="flex-1 h-px bg-gradient-to-r from-transparent to-slate-700" />
            <Sparkles className="w-3.5 h-3.5 text-slate-600" />
            <div className="flex-1 h-px bg-gradient-to-l from-transparent to-slate-700" />
          </div>

          {/* Action */}
          <button
            onClick={handleBack}
            className="
              group w-full flex items-center justify-center gap-2.5
              bg-gradient-to-r from-slate-800 to-slate-800/80
              hover:from-blue-600 hover:to-blue-500
              border border-slate-700 hover:border-blue-500
              text-slate-300 hover:text-white
              px-5 py-3 rounded-xl font-medium text-sm
              transition-all duration-300
              shadow-lg hover:shadow-blue-500/20
            "
          >
            <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-0.5" />
            Back to Marketplace
          </button>
        </div>
      </div>
    </div>
  );
}
