import type { IChatRepository } from '../../domains/chat/repositories/ChatRepository';
import type { Message, SendMessageParams } from '../../domains/chat/entities/Message';
import apiClient from './apiClient';

export class ChatRepository implements IChatRepository {
  async getMessages(transactionId: number): Promise<Message[]> {
    const { data } = await apiClient.get(`/chat/transactions/${transactionId}/messages`);
    return data;
  }

  async sendMessage(transactionId: number, params: SendMessageParams): Promise<Message> {
    const { data } = await apiClient.post(`/chat/transactions/${transactionId}/messages`, params);
    return data;
  }
}
