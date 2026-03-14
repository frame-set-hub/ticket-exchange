import type { IChatRepository } from '../../domains/chat/repositories/ChatRepository';
import type { Message, SendMessageParams } from '../../domains/chat/entities/Message';

export class ChatRepository implements IChatRepository {
  private ws: WebSocket | null = null;

  connect(token: string, onMessage: (msg: Message) => void): void {
    this.ws = new WebSocket(`ws://localhost:8080/api/chat?token=${token}`);
    this.ws.onmessage = (event) => {
      const msg: Message = JSON.parse(event.data);
      onMessage(msg);
    };
  }

  send(params: SendMessageParams): void {
    if (!this.ws) return;
    this.ws.send(JSON.stringify(params));
  }

  disconnect(): void {
    this.ws?.close();
    this.ws = null;
  }
}
