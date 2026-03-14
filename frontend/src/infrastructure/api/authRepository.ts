import type { IAuthRepository, LoginParams, LoginResult, RegisterParams } from '../../domains/auth/repositories/AuthRepository';
import apiClient from './apiClient';

export class AuthRepository implements IAuthRepository {
  async login(params: LoginParams): Promise<LoginResult> {
    const { data } = await apiClient.post('/auth/login', params);
    return data;
  }

  async register(params: RegisterParams): Promise<void> {
    await apiClient.post('/auth/register', params);
  }
}
