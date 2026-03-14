import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom';
import { useAuthStore } from './store/useStore';
import { ToastProvider } from './components/Toast';
import Login from './pages/Login';
import Register from './pages/Register';
import Home from './pages/Home';
import AdminDashboard from './pages/AdminDashboard';
import Transaction from './pages/Transaction';
import MyTickets from './pages/MyTickets';

// Layout shell component with Navbar
const Layout = ({ children }: { children: React.ReactNode }) => {
  const { user, logout } = useAuthStore();

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col pt-16">
      <nav className="fixed top-0 left-0 right-0 h-16 glass-panel z-50 flex items-center justify-between px-6 border-b border-slate-800">
        <div className="flex items-center gap-2 text-xl font-bold text-white tracking-tight">
          <span className="text-blue-500">Ticket</span>X
          <span className="text-xs font-normal px-2 py-0.5 rounded-full bg-blue-500/20 text-blue-300 border border-blue-500/30">Escrow POC</span>
        </div>

        <div className="flex items-center gap-6">
          <div className="hidden md:flex items-center gap-6 text-sm font-medium">
            <Link to="/" className="text-slate-300 hover:text-white transition-colors">Marketplace</Link>
            {user && <Link to="/my-tickets" className="text-slate-300 hover:text-white transition-colors">My Tickets</Link>}
            {user?.role === 'Admin' && <Link to="/admin" className="text-slate-300 hover:text-white transition-colors">Dashboard</Link>}
          </div>

          {user ? (
            <div className="flex items-center gap-4">
              <div className="text-sm">
                <span className="text-slate-400">Signed in as </span>
                <span className="text-white font-medium">{user.username} <span className="opacity-50 text-xs">({user.role})</span></span>
              </div>
              <button onClick={logout} className="text-sm bg-slate-800 hover:bg-slate-700 text-white px-4 py-2 rounded-lg transition-colors border border-slate-700">
                Log Out
              </button>
            </div>
          ) : (
            <div className="space-x-3">
              <Link to="/login" className="text-sm font-medium text-slate-300 hover:text-white transition-colors">Log In</Link>
              <Link to="/register" className="text-sm font-medium bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg transition-colors">Sign Up</Link>
            </div>
          )}
        </div>
      </nav>
      <main className="flex-1 w-full max-w-7xl mx-auto p-6">
        {children}
      </main>
    </div>
  );
};

// Protect Routes
const ProtectedRoute = ({ children, adminOnly = false }: { children: React.ReactNode, adminOnly?: boolean }) => {
  const { user, token } = useAuthStore();
  if (!token || !user) return <Navigate to="/login" replace />;
  if (adminOnly && user.role !== 'Admin') return <Navigate to="/" replace />;
  return children;
};

export default function App() {
  return (
    <ToastProvider>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route path="/" element={<Layout><ProtectedRoute><Home /></ProtectedRoute></Layout>} />
        <Route path="/my-tickets" element={<Layout><ProtectedRoute><MyTickets /></ProtectedRoute></Layout>} />
        <Route path="/checkout/:id" element={<Layout><ProtectedRoute><Transaction /></ProtectedRoute></Layout>} />
        <Route path="/admin" element={<Layout><ProtectedRoute adminOnly><AdminDashboard /></ProtectedRoute></Layout>} />
      </Routes>
    </BrowserRouter>
    </ToastProvider>
  );
}
