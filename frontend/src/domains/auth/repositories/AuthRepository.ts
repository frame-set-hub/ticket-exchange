import { User } from '../entities/User';

export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  email: string;
  password: string;
  role: string;
}

export interface LoginResult {
  user: User;
  token: string;
}

export interface IAuthRepository {
  login(params: LoginParams): Promise<LoginResult>;
  register(params: RegisterParams): Promise<void>;
}
