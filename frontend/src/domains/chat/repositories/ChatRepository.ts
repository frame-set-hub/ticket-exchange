import { Message, SendMessageParams } from '../entities/Message';

export interface IChatRepository {
  connect(token: string, onMessage: (msg: Message) => void): void;
  send(params: SendMessageParams): void;
  disconnect(): void;
}
