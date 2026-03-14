import type { Message, SendMessageParams } from '../entities/Message';

export interface IChatRepository {
  getMessages(transactionId: number): Promise<Message[]>;
  sendMessage(transactionId: number, params: SendMessageParams): Promise<Message>;
}
