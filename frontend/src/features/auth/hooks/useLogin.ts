import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../store/useStore';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useLogin() {
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const login = useAuthStore((state) => state.login);
  const navigate = useNavigate();
  const { authRepository } = getServices();

  const handleLogin = async (username: string, password: string) => {
    setError('');
    setLoading(true);
    try {
      const result = await authRepository.login({ username, password });
      login(result.user, result.token);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return { handleLogin, error, loading };
}
