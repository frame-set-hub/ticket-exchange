import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useRegister() {
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { authRepository } = getServices();

  const handleRegister = async (username: string, email: string, password: string, role: string) => {
    setError('');
    setLoading(true);
    try {
      await authRepository.register({ username, email, password, role });
      navigate('/login');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return { handleRegister, error, loading };
}
