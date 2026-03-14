import { useState, useEffect, useCallback } from 'react';
import type { Message } from '../../../domains/chat/entities/Message';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useChat(token: string | null, transactionId: number | undefined) {
  const [messages, setMessages] = useState<Message[]>([]);
  const { chatRepository } = getServices();

  const connect = useCallback(() => {
    if (!token || !transactionId) return;

    chatRepository.connect(token, (msg) => {
      if (msg.transaction_id === transactionId) {
        setMessages((prev) => [...prev, msg]);
      }
    });
  }, [token, transactionId]);

  const sendMessage = (content: string, receiverId: number) => {
    if (!transactionId || !content) return;
    chatRepository.send({
      transaction_id: transactionId,
      content,
      receiver_id: receiverId,
    });
  };

  useEffect(() => {
    connect();
    return () => chatRepository.disconnect();
  }, [connect]);

  return { messages, sendMessage };
}
